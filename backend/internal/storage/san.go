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

type SANStorage struct {
	config     StorageConfig
	devicePath string
	mountPoint string
	protocol   string // iscsi, fc
	encryptor  *crypto.Encryptor
}

func NewSANStorage() *SANStorage {
	return &SANStorage{}
}

func (ss *SANStorage) Configure(config StorageConfig) error {
	ss.config = config
	
	if devicePath, ok := config.Options["device_path"].(string); ok {
		ss.devicePath = devicePath
	} else {
		ss.devicePath = "/dev/sdb1"
	}
	
	if mountPoint, ok := config.Options["mount_point"].(string); ok {
		ss.mountPoint = mountPoint
	} else {
		ss.mountPoint = "/mnt/san"
	}
	
	if protocol, ok := config.Options["protocol"].(string); ok {
		ss.protocol = protocol
	} else {
		ss.protocol = "iscsi"
	}
	
	if config.Encryption.Enabled {
		ss.encryptor = crypto.NewEncryptor()
		if err := ss.encryptor.Configure(config.Encryption); err != nil {
			return fmt.Errorf("failed to configure encryption: %w", err)
		}
	}
	
	if err := ss.connectAndMount(); err != nil {
		return fmt.Errorf("failed to connect to SAN: %w", err)
	}
	
	return nil
}

func (ss *SANStorage) connectAndMount() error {
	switch ss.protocol {
	case "iscsi":
		if err := ss.connectiSCSI(); err != nil {
			return err
		}
	case "fc":
		if err := ss.connectFC(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported SAN protocol: %s", ss.protocol)
	}
	
	return ss.mountDevice()
}

func (ss *SANStorage) connectiSCSI() error {
	parts := strings.Split(ss.config.Endpoint, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid iSCSI endpoint format: %s", ss.config.Endpoint)
	}
	
	target := parts[0]
	portal := strings.Join(parts[1:], ":")
	
	cmd := exec.Command("iscsiadm", "-m", "discovery", "-t", "st", "-p", portal)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("iSCSI discovery failed: %s, output: %s", err, string(output))
	}
	
	cmd = exec.Command("iscsiadm", "-m", "node", "-T", target, "-p", portal, "--login")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("iSCSI login failed: %s, output: %s", err, string(output))
	}
	
	time.Sleep(2 * time.Second)
	
	return nil
}

func (ss *SANStorage) connectFC() error {
	if _, err := os.Stat(ss.devicePath); err != nil {
		return fmt.Errorf("FC device not found: %s", ss.devicePath)
	}
	
	return nil
}

func (ss *SANStorage) mountDevice() error {
	if ss.isMounted() {
		return nil
	}
	
	if err := os.MkdirAll(ss.mountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %w", err)
	}
	
	if _, err := os.Stat(ss.devicePath); err != nil {
		return fmt.Errorf("device not found: %s", err)
	}
	
	cmd := exec.Command("mount", ss.devicePath, ss.mountPoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount device: %s, output: %s", err, string(output))
	}
	
	return nil
}

func (ss *SANStorage) isMounted() bool {
	cmd := exec.Command("mountpoint", "-q", ss.mountPoint)
	return cmd.Run() == nil
}

func (ss *SANStorage) WriteData(ctx context.Context, key string, data []byte) error {
	if !ss.isMounted() {
		return fmt.Errorf("SAN not mounted")
	}
	
	filePath := filepath.Join(ss.mountPoint, key)
	
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	if ss.encryptor != nil {
		encryptedData, err := ss.encryptor.Encrypt(data)
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

func (ss *SANStorage) ReadData(ctx context.Context, key string) ([]byte, error) {
	if !ss.isMounted() {
		return nil, fmt.Errorf("SAN not mounted")
	}
	
	filePath := filepath.Join(ss.mountPoint, key)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	if ss.encryptor != nil {
		decryptedData, err := ss.encryptor.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}
		data = decryptedData
	}
	
	return data, nil
}

func (ss *SANStorage) DeleteData(ctx context.Context, key string) error {
	if !ss.isMounted() {
		return fmt.Errorf("SAN not mounted")
	}
	
	filePath := filepath.Join(ss.mountPoint, key)
	
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

func (ss *SANStorage) ListData(ctx context.Context, prefix string) ([]string, error) {
	if !ss.isMounted() {
		return nil, fmt.Errorf("SAN not mounted")
	}
	
	var keys []string
	
	err := filepath.WalkDir(ss.mountPoint, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() {
			return nil
		}
		
		relPath, err := filepath.Rel(ss.mountPoint, path)
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

func (ss *SANStorage) GetHealth() HealthStatus {
	health := HealthStatus{
		Status:    HealthStateHealthy,
		LastCheck: time.Now(),
		Errors:    []string{},
		Warnings:  []string{},
	}
	
	if !ss.isMounted() {
		health.Status = HealthStateCritical
		health.Errors = append(health.Errors, "SAN not mounted")
		return health
	}
	
	if _, err := os.Stat(ss.devicePath); err != nil {
		health.Status = HealthStateCritical
		health.Errors = append(health.Errors, fmt.Sprintf("Device not accessible: %v", err))
		return health
	}
	
	if _, err := os.Stat(ss.mountPoint); err != nil {
		health.Status = HealthStateCritical
		health.Errors = append(health.Errors, fmt.Sprintf("Mount point not accessible: %v", err))
		return health
	}
	
	testFile := filepath.Join(ss.mountPoint, ".health_check")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		health.Status = HealthStateWarning
		health.Warnings = append(health.Warnings, "Write access test failed")
	} else {
		os.Remove(testFile) // Clean up test file
	}
	
	return health
}

func (ss *SANStorage) GetMetrics() StorageMetrics {
	metrics := StorageMetrics{
		LastUpdated: time.Now(),
	}
	
	if ss.isMounted() {
		if usage, err := getDiskUsage(ss.mountPoint); err == nil {
			metrics.TotalSpace = usage.Total
			metrics.UsedSpace = usage.Used
			metrics.FreeSpace = usage.Free
		}
	}
	
	return metrics
}

func (ss *SANStorage) Close() error {
	return nil
}

func (ss *SANStorage) Disconnect() error {
	if ss.isMounted() {
		cmd := exec.Command("umount", ss.mountPoint)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to unmount SAN: %s, output: %s", err, string(output))
		}
	}
	
	switch ss.protocol {
	case "iscsi":
		return ss.disconnectiSCSI()
	case "fc":
		return nil
	}
	
	return nil
}

func (ss *SANStorage) disconnectiSCSI() error {
	parts := strings.Split(ss.config.Endpoint, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid iSCSI endpoint format: %s", ss.config.Endpoint)
	}
	
	target := parts[0]
	portal := strings.Join(parts[1:], ":")
	
	cmd := exec.Command("iscsiadm", "-m", "node", "-T", target, "-p", portal, "--logout")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("iSCSI logout failed: %s, output: %s", err, string(output))
	}
	
	return nil
}
