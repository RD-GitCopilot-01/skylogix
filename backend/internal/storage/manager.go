package storage

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Manager struct {
	backends map[string]StorageBackend
	mutex    sync.RWMutex
	config   ManagerConfig
}

type ManagerConfig struct {
	DefaultBackend string                    `json:"default_backend"`
	TieringRules   []TieringRule            `json:"tiering_rules"`
	Replication    ReplicationConfig        `json:"replication"`
	Retention      RetentionConfig          `json:"retention"`
	Monitoring     MonitoringConfig         `json:"monitoring"`
}

type TieringRule struct {
	Name        string        `json:"name"`
	Condition   string        `json:"condition"`
	SourceTier  string        `json:"source_tier"`
	TargetTier  string        `json:"target_tier"`
	Age         time.Duration `json:"age"`
	Size        int64         `json:"size"`
	AccessCount int           `json:"access_count"`
}

type ReplicationConfig struct {
	Enabled     bool     `json:"enabled"`
	Targets     []string `json:"targets"`
	SyncMode    string   `json:"sync_mode"` // async, sync
	Compression bool     `json:"compression"`
}

type RetentionConfig struct {
	DefaultTTL    time.Duration            `json:"default_ttl"`
	ChannelRules  map[string]time.Duration `json:"channel_rules"`
	QuarantineTTL time.Duration            `json:"quarantine_ttl"`
	ProtectedMode bool                     `json:"protected_mode"`
}

type MonitoringConfig struct {
	Enabled         bool          `json:"enabled"`
	CheckInterval   time.Duration `json:"check_interval"`
	AlertThresholds AlertThresholds `json:"alert_thresholds"`
}

type AlertThresholds struct {
	DiskUsageWarning  float64 `json:"disk_usage_warning"`
	DiskUsageCritical float64 `json:"disk_usage_critical"`
	IOPSThreshold     uint64  `json:"iops_threshold"`
	LatencyThreshold  float64 `json:"latency_threshold"`
}

func NewManager() *Manager {
	return &Manager{
		backends: make(map[string]StorageBackend),
		config: ManagerConfig{
			DefaultBackend: "local",
			Monitoring: MonitoringConfig{
				Enabled:       true,
				CheckInterval: 30 * time.Second,
				AlertThresholds: AlertThresholds{
					DiskUsageWarning:  80.0,
					DiskUsageCritical: 90.0,
					IOPSThreshold:     1000,
					LatencyThreshold:  100.0,
				},
			},
		},
	}
}

func (m *Manager) Configure(config ManagerConfig) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.config = config
	return nil
}

func (m *Manager) AddBackend(name string, backend StorageBackend) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, exists := m.backends[name]; exists {
		return fmt.Errorf("backend %s already exists", name)
	}
	
	m.backends[name] = backend
	return nil
}

func (m *Manager) RemoveBackend(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	backend, exists := m.backends[name]
	if !exists {
		return fmt.Errorf("backend %s not found", name)
	}
	
	if err := backend.Close(); err != nil {
		return fmt.Errorf("failed to close backend %s: %w", name, err)
	}
	
	delete(m.backends, name)
	return nil
}

func (m *Manager) GetBackend(name string) (StorageBackend, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	backend, exists := m.backends[name]
	if !exists {
		return nil, fmt.Errorf("backend %s not found", name)
	}
	
	return backend, nil
}

func (m *Manager) ListBackends() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	names := make([]string, 0, len(m.backends))
	for name := range m.backends {
		names = append(names, name)
	}
	
	return names
}

func (m *Manager) WriteData(ctx context.Context, key string, data []byte) error {
	backend, err := m.selectBackendForWrite(key, data)
	if err != nil {
		return err
	}
	
	return backend.WriteData(ctx, key, data)
}

func (m *Manager) ReadData(ctx context.Context, key string) ([]byte, error) {
	for name, backend := range m.backends {
		data, err := backend.ReadData(ctx, key)
		if err == nil {
			return data, nil
		}
		
		fmt.Printf("Backend %s: failed to read key %s: %v\n", name, key, err)
	}
	
	return nil, fmt.Errorf("key not found in any backend: %s", key)
}

func (m *Manager) DeleteData(ctx context.Context, key string) error {
	var lastErr error
	deleted := false
	
	for name, backend := range m.backends {
		err := backend.DeleteData(ctx, key)
		if err == nil {
			deleted = true
		} else {
			lastErr = err
			fmt.Printf("Backend %s: failed to delete key %s: %v\n", name, key, err)
		}
	}
	
	if !deleted && lastErr != nil {
		return lastErr
	}
	
	return nil
}

func (m *Manager) selectBackendForWrite(key string, data []byte) (StorageBackend, error) {
	for _, rule := range m.config.TieringRules {
		if m.matchesTieringRule(rule, key, data) {
			if backend, err := m.GetBackend(rule.TargetTier); err == nil {
				return backend, nil
			}
		}
	}
	
	return m.GetBackend(m.config.DefaultBackend)
}

func (m *Manager) matchesTieringRule(rule TieringRule, key string, data []byte) bool {
	if rule.Size > 0 && int64(len(data)) > rule.Size {
		return true
	}
	
	return false
}

func (m *Manager) GetAggregatedMetrics() StorageMetrics {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	aggregated := StorageMetrics{
		LastUpdated: time.Now(),
	}
	
	for _, backend := range m.backends {
		metrics := backend.GetMetrics()
		aggregated.TotalSpace += metrics.TotalSpace
		aggregated.UsedSpace += metrics.UsedSpace
		aggregated.FreeSpace += metrics.FreeSpace
		aggregated.IOPS.Read += metrics.IOPS.Read
		aggregated.IOPS.Write += metrics.IOPS.Write
		aggregated.IOPS.Total += metrics.IOPS.Total
		aggregated.Throughput.ReadMBps += metrics.Throughput.ReadMBps
		aggregated.Throughput.WriteMBps += metrics.Throughput.WriteMBps
		aggregated.Throughput.TotalMBps += metrics.Throughput.TotalMBps
	}
	
	return aggregated
}

func (m *Manager) GetAggregatedHealth() HealthStatus {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	aggregated := HealthStatus{
		Status:    HealthStateHealthy,
		LastCheck: time.Now(),
		Errors:    []string{},
		Warnings:  []string{},
	}
	
	for name, backend := range m.backends {
		health := backend.GetHealth()
		
		for _, err := range health.Errors {
			aggregated.Errors = append(aggregated.Errors, fmt.Sprintf("%s: %s", name, err))
		}
		for _, warning := range health.Warnings {
			aggregated.Warnings = append(aggregated.Warnings, fmt.Sprintf("%s: %s", name, warning))
		}
		
		if health.Status == HealthStateCritical {
			aggregated.Status = HealthStateCritical
		} else if health.Status == HealthStateWarning && aggregated.Status == HealthStateHealthy {
			aggregated.Status = HealthStateWarning
		}
	}
	
	return aggregated
}

func (m *Manager) MigrateData(from, to string, keys []string) error {
	sourceBackend, err := m.GetBackend(from)
	if err != nil {
		return fmt.Errorf("source backend not found: %w", err)
	}
	
	targetBackend, err := m.GetBackend(to)
	if err != nil {
		return fmt.Errorf("target backend not found: %w", err)
	}
	
	ctx := context.Background()
	
	for _, key := range keys {
		data, err := sourceBackend.ReadData(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to read key %s from source: %w", key, err)
		}
		
		if err := targetBackend.WriteData(ctx, key, data); err != nil {
			return fmt.Errorf("failed to write key %s to target: %w", key, err)
		}
		
	}
	
	return nil
}

func (m *Manager) StartMonitoring() {
	if !m.config.Monitoring.Enabled {
		return
	}
	
	go func() {
		ticker := time.NewTicker(m.config.Monitoring.CheckInterval)
		defer ticker.Stop()
		
		for range ticker.C {
			m.performHealthCheck()
		}
	}()
}

func (m *Manager) performHealthCheck() {
	health := m.GetAggregatedHealth()
	metrics := m.GetAggregatedMetrics()
	
	if metrics.TotalSpace > 0 {
		usagePercent := float64(metrics.UsedSpace) / float64(metrics.TotalSpace) * 100
		
		if usagePercent >= m.config.Monitoring.AlertThresholds.DiskUsageCritical {
			m.sendAlert("CRITICAL", fmt.Sprintf("Disk usage at %.1f%%", usagePercent))
		} else if usagePercent >= m.config.Monitoring.AlertThresholds.DiskUsageWarning {
			m.sendAlert("WARNING", fmt.Sprintf("Disk usage at %.1f%%", usagePercent))
		}
	}
	
	if health.Status == HealthStateCritical {
		m.sendAlert("CRITICAL", "Storage backend critical errors detected")
	} else if health.Status == HealthStateWarning {
		m.sendAlert("WARNING", "Storage backend warnings detected")
	}
}

func (m *Manager) sendAlert(level, message string) {
	fmt.Printf("ALERT [%s]: %s\n", level, message)
}

func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	var lastErr error
	for name, backend := range m.backends {
		if err := backend.Close(); err != nil {
			lastErr = err
			fmt.Printf("Failed to close backend %s: %v\n", name, err)
		}
	}
	
	return lastErr
}
