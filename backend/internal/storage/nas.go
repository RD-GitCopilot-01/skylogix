package storage

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	
	"skylogix-storage/internal/crypto"
)

type NASStorage struct {
	config     StorageConfig
	mountPoint string
	protocol   string // nfs, smb, etc.
	encryptor  *crypto.Encryptor
}

func NewNASStorage() *NASStorage {
	return &NASStorage{}
}

func (ns *NASStorage) Configure(config StorageConfig) error {
	ns.config = config
	
	if mountPoint, ok := config.Options["mount_point"].(string); ok {
		ns.mountPoint = mountPoint
	} else {
		ns.mountPoint = "/mnt/nas"
	}
	
	if protocol, ok := config.Options["protocol"].(string); ok {
		ns.protocol = protocol
	} else {
		ns.protocol = "nfs"
	}
	
	if config.Encryption.Enabled {
		ns.encryptor = crypto.NewEncryptor()
		if err := ns.encryptor.Configure(config.Encryption); err != nil {
			return fmt.Errorf("failed to configure encryption: %w", err)
		}
	}
	
	if err := ns.ensureMounted(); err != nil {
		return fmt.Errorf("failed to mount NAS: %w", err)
	}
	
	return nil
}

func (ns *NASStorage) ensureMounted() error {
	if ns.isMounted() {
		return nil
	}
	
	if err := os.MkdirAll(ns.mountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %w", err)
	}
	
	switch ns.protocol {
	case "nfs":
		return ns.mountNFS()
	case "smb", "cifs":
		return ns.mountSMB()
	default:
		return fmt.Errorf("unsupported NAS protocol: %s", ns.protocol)
	}
}

func (ns *NASStorage) isMounted() bool {
	cmd := exec.Command("mountpoint", "-q", ns.mountPoint)
	return cmd.Run() == nil
}

func (ns *NASStorage) mountNFS() error {
	args := []string{
		"-t", "nfs",
		ns.config.Endpoint,
		ns.mountPoint,
	}
	
	options := []string{
		"rw",
		"noatime",
		"rsize=1048576",
		"wsize=1048576",
		"tcp",
	}
	
	if len(options) > 0 {
		args = append(args, "-o", strings.Join(options, ","))
	}
	
	cmd := exec.Command("mount", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount NFS: %s, output: %s", err, string(output))
	}
	
	return nil
}

func (ns *NASStorage) mountSMB() error {
	args := []string{
		"-t", "cifs",
		ns.config.Endpoint,
		ns.mountPoint,
	}
	
	options := []string{
		"rw",
		"noatime",
	}
	
	if ns.config.Credentials.Username != "" {
		options = append(options, fmt.Sprintf("username=%s", ns.config.Credentials.Username))
	}
	if ns.config.Credentials.Password != "" {
		options = append(options, fmt.Sprintf("password=%s", ns.config.Credentials.Password))
	}
	
	if len(options) > 0 {
		args = append(args, "-o", strings.Join(options, ","))
	}
	
	cmd := exec.Command("mount", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount SMB: %s, output: %s", err, string(output))
	}
	
	return nil
}

func (ns *NASStorage) WriteData(ctx context.Context, key string, data []byte) error {
	if !ns.isMounted() {
		return fmt.Errorf("NAS not mounted")
	}
	
	filePath := filepath.Join(ns.mountPoint, key)
	
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	if ns.encryptor != nil {
		encryptedData, err := ns.encryptor.Encrypt(data)
		if err != nil {
			return fmt.Errorf("failed to encrypt data: %w", err)
		}
		data = encryptedData
	}
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

func (ns *NASStorage) ReadData(ctx context.Context, key string) ([]byte, error) {
	if !ns.isMounted() {
		return nil, fmt.Errorf("NAS not mounted")
	}
	
	filePath := filepath.Join(ns.mountPoint, key)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	if ns.encryptor != nil {
		decryptedData, err := ns.encryptor.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}
		data = decryptedData
	}
	
	return data, nil
}

func (ns *NASStorage) DeleteData(ctx context.Context, key string) error {
	if !ns.isMounted() {
		return fmt.Errorf("NAS not mounted")
	}
	
	filePath := filepath.Join(ns.mountPoint, key)
	
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

func (ns *NASStorage) ListData(ctx context.Context, prefix string) ([]string, error) {
	if !ns.isMounted() {
		return nil, fmt.Errorf("NAS not mounted")
	}
	
	var keys []string
	
	err := filepath.WalkDir(ns.mountPoint, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() {
			return nil
		}
		
		relPath, err := filepath.Rel(ns.mountPoint, path)
		if err != nil {
			return err
		}
		
		if strings.HasPrefix(relPath, prefix) {
			keys = append(keys, relPath)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	
	return keys, nil
}

func (ns *NASStorage) GetHealth() HealthStatus {
	health := HealthStatus{
		Status:    HealthStateHealthy,
		LastCheck: time.Now(),
		Errors:    []string{},
		Warnings:  []string{},
	}
	
	if !ns.isMounted() {
		health.Status = HealthStateCritical
		health.Errors = append(health.Errors, "NAS not mounted")
		return health
	}
	
	if _, err := os.Stat(ns.mountPoint); err != nil {
		health.Status = HealthStateCritical
		health.Errors = append(health.Errors, fmt.Sprintf("Mount point not accessible: %v", err))
		return health
	}
	
	testFile := filepath.Join(ns.mountPoint, ".health_check")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		health.Status = HealthStateWarning
		health.Warnings = append(health.Warnings, "Write access test failed")
	} else {
		os.Remove(testFile) // Clean up test file
	}
	
	return health
}

func (ns *NASStorage) GetMetrics() StorageMetrics {
	metrics := StorageMetrics{
		LastUpdated: time.Now(),
	}
	
	if ns.isMounted() {
		if usage, err := getDiskUsage(ns.mountPoint); err == nil {
			metrics.TotalSpace = usage.Total
			metrics.UsedSpace = usage.Used
			metrics.FreeSpace = usage.Free
		}
	}
	
	return metrics
}

func (ns *NASStorage) Close() error {
	return nil
}

func (ns *NASStorage) Unmount() error {
	if !ns.isMounted() {
		return nil
	}
	
	cmd := exec.Command("umount", ns.mountPoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unmount NAS: %s, output: %s", err, string(output))
	}
	
	return nil
}
