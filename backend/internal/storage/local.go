package storage

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"skylogix-storage/internal/crypto"
)

type LocalStorage struct {
	config    StorageConfig
	basePath  string
	raidMgr   *RAIDManager
	encryptor *crypto.Encryptor
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (ls *LocalStorage) Configure(config StorageConfig) error {
	ls.config = config
	
	if config.Endpoint != "" {
		ls.basePath = config.Endpoint
	} else {
		ls.basePath = "/data/minio"
	}
	
	if err := os.MkdirAll(ls.basePath, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}
	
	if config.RAID.Enabled {
		ls.raidMgr = NewRAIDManager()
		if err := ls.raidMgr.Configure(config.RAID); err != nil {
			return fmt.Errorf("failed to configure RAID: %w", err)
		}
	}
	
	if config.Encryption.Enabled {
		ls.encryptor = crypto.NewEncryptor()
		if err := ls.encryptor.Configure(config.Encryption); err != nil {
			return fmt.Errorf("failed to configure encryption: %w", err)
		}
	}
	
	return nil
}

func (ls *LocalStorage) WriteData(ctx context.Context, key string, data []byte) error {
	filePath := filepath.Join(ls.basePath, key)
	
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	if ls.encryptor != nil {
		encryptedData, err := ls.encryptor.Encrypt(data)
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

func (ls *LocalStorage) ReadData(ctx context.Context, key string) ([]byte, error) {
	filePath := filepath.Join(ls.basePath, key)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	if ls.encryptor != nil {
		decryptedData, err := ls.encryptor.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}
		data = decryptedData
	}
	
	return data, nil
}

func (ls *LocalStorage) DeleteData(ctx context.Context, key string) error {
	filePath := filepath.Join(ls.basePath, key)
	
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

func (ls *LocalStorage) ListData(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	
	err := filepath.WalkDir(ls.basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() {
			return nil
		}
		
		relPath, err := filepath.Rel(ls.basePath, path)
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

func (ls *LocalStorage) GetHealth() HealthStatus {
	health := HealthStatus{
		Status:    HealthStateHealthy,
		LastCheck: time.Now(),
		Errors:    []string{},
		Warnings:  []string{},
	}
	
	if _, err := os.Stat(ls.basePath); err != nil {
		health.Status = HealthStateCritical
		health.Errors = append(health.Errors, fmt.Sprintf("Base path not accessible: %v", err))
		return health
	}
	
	if ls.raidMgr != nil {
		raidHealth := ls.raidMgr.GetHealth()
		health.RAIDStatus = raidHealth
		
		if raidHealth.FailedDisks > 0 {
			health.Status = HealthStateCritical
			health.Errors = append(health.Errors, fmt.Sprintf("RAID has %d failed disks", raidHealth.FailedDisks))
		}
	}
	
	if usage, err := getDiskUsage(ls.basePath); err == nil {
		if usage.UsedPercent > 90 {
			health.Status = HealthStateCritical
			health.Errors = append(health.Errors, "Disk usage above 90%")
		} else if usage.UsedPercent > 80 {
			if health.Status == HealthStateHealthy {
				health.Status = HealthStateWarning
			}
			health.Warnings = append(health.Warnings, "Disk usage above 80%")
		}
	}
	
	return health
}

func (ls *LocalStorage) GetMetrics() StorageMetrics {
	metrics := StorageMetrics{
		LastUpdated: time.Now(),
	}
	
	if usage, err := getDiskUsage(ls.basePath); err == nil {
		metrics.TotalSpace = usage.Total
		metrics.UsedSpace = usage.Used
		metrics.FreeSpace = usage.Free
	}
	
	
	return metrics
}

func (ls *LocalStorage) Close() error {
	if ls.raidMgr != nil {
		return ls.raidMgr.Close()
	}
	return nil
}

type DiskUsage struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func getDiskUsage(path string) (*DiskUsage, error) {
	
	var usage DiskUsage
	
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if !d.IsDir() {
			if info, err := d.Info(); err == nil {
				usage.Used += uint64(info.Size())
			}
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	usage.Total = 1024 * 1024 * 1024 * 1024 // 1TB
	usage.Free = usage.Total - usage.Used
	usage.UsedPercent = float64(usage.Used) / float64(usage.Total) * 100
	
	return &usage, nil
}
