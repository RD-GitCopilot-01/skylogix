package storage

import (
	"context"
	"time"
	
	"skylogix-storage/internal/crypto"
)

type StorageBackend interface {
	Configure(config StorageConfig) error
	WriteData(ctx context.Context, key string, data []byte) error
	ReadData(ctx context.Context, key string) ([]byte, error)
	DeleteData(ctx context.Context, key string) error
	ListData(ctx context.Context, prefix string) ([]string, error)
	GetHealth() HealthStatus
	GetMetrics() StorageMetrics
	Close() error
}

type StorageConfig struct {
	Type        StorageType            `json:"type"`
	Endpoint    string                 `json:"endpoint"`
	Credentials StorageCredentials     `json:"credentials"`
	RAID        RAIDConfig             `json:"raid"`
	Encryption  EncryptionConfig       `json:"encryption"`
	Options     map[string]interface{} `json:"options"`
}

type StorageType string

const (
	StorageTypeLocal  StorageType = "local"
	StorageTypeNAS    StorageType = "nas"
	StorageTypeSAN    StorageType = "san"
	StorageTypeObject StorageType = "object"
)

type StorageCredentials struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

type RAIDConfig struct {
	Level       RAIDLevel `json:"level"`
	Devices     []string  `json:"devices"`
	SpareDevices []string `json:"spare_devices"`
	ChunkSize   int       `json:"chunk_size"`
	Enabled     bool      `json:"enabled"`
}

type RAIDLevel string

const (
	RAID0  RAIDLevel = "0"
	RAID1  RAIDLevel = "1"
	RAID5  RAIDLevel = "5"
	RAID6  RAIDLevel = "6"
	RAID10 RAIDLevel = "10"
)

type EncryptionConfig = crypto.EncryptionConfig

type HealthStatus struct {
	Status      HealthState       `json:"status"`
	LastCheck   time.Time         `json:"last_check"`
	Errors      []string          `json:"errors"`
	Warnings    []string          `json:"warnings"`
	RAIDStatus  RAIDHealthStatus  `json:"raid_status"`
	DiskStatus  []DiskHealthStatus `json:"disk_status"`
}

type HealthState string

const (
	HealthStateHealthy  HealthState = "healthy"
	HealthStateWarning  HealthState = "warning"
	HealthStateCritical HealthState = "critical"
	HealthStateUnknown  HealthState = "unknown"
)

type RAIDHealthStatus struct {
	State       string    `json:"state"`
	Level       RAIDLevel `json:"level"`
	ActiveDisks int       `json:"active_disks"`
	FailedDisks int       `json:"failed_disks"`
	SpareDisks  int       `json:"spare_disks"`
	SyncStatus  string    `json:"sync_status"`
	LastSync    time.Time `json:"last_sync"`
}

type DiskHealthStatus struct {
	Device      string    `json:"device"`
	Status      string    `json:"status"`
	Temperature int       `json:"temperature"`
	Errors      int       `json:"errors"`
	LastCheck   time.Time `json:"last_check"`
}

type StorageMetrics struct {
	TotalSpace     uint64        `json:"total_space"`
	UsedSpace      uint64        `json:"used_space"`
	FreeSpace      uint64        `json:"free_space"`
	IOPS           IOPSMetrics   `json:"iops"`
	Throughput     ThroughputMetrics `json:"throughput"`
	Latency        LatencyMetrics `json:"latency"`
	ErrorRate      float64       `json:"error_rate"`
	LastUpdated    time.Time     `json:"last_updated"`
}

type IOPSMetrics struct {
	Read  uint64 `json:"read"`
	Write uint64 `json:"write"`
	Total uint64 `json:"total"`
}

type ThroughputMetrics struct {
	ReadMBps  float64 `json:"read_mbps"`
	WriteMBps float64 `json:"write_mbps"`
	TotalMBps float64 `json:"total_mbps"`
}

type LatencyMetrics struct {
	ReadMs  float64 `json:"read_ms"`
	WriteMs float64 `json:"write_ms"`
	AvgMs   float64 `json:"avg_ms"`
}

type StorageManager interface {
	AddBackend(name string, backend StorageBackend) error
	RemoveBackend(name string) error
	GetBackend(name string) (StorageBackend, error)
	ListBackends() []string
	GetAggregatedMetrics() StorageMetrics
	GetAggregatedHealth() HealthStatus
	MigrateData(from, to string, keys []string) error
}
