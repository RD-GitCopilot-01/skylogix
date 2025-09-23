package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Manager struct {
	config     *Config
	configPath string
	mutex      sync.RWMutex
	watchers   []ConfigWatcher
}

type Config struct {
	Server    ServerConfig    `json:"server"`
	Storage   StorageConfig   `json:"storage"`
	SNMP      SNMPConfig      `json:"snmp"`
	Database  DatabaseConfig  `json:"database"`
	Security  SecurityConfig  `json:"security"`
	Logging   LoggingConfig   `json:"logging"`
	Retention RetentionConfig `json:"retention"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ServerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	TLS          TLSConfig     `json:"tls"`
}

type TLSConfig struct {
	Enabled  bool   `json:"enabled"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type StorageConfig struct {
	DefaultBackend string                    `json:"default_backend"`
	Backends       map[string]BackendConfig  `json:"backends"`
	TieringRules   []TieringRuleConfig       `json:"tiering_rules"`
	Replication    ReplicationConfig         `json:"replication"`
	Encryption     EncryptionConfig          `json:"encryption"`
}

type BackendConfig struct {
	Type        string                 `json:"type"`
	Endpoint    string                 `json:"endpoint"`
	Credentials CredentialsConfig      `json:"credentials"`
	RAID        RAIDConfig             `json:"raid"`
	Options     map[string]interface{} `json:"options"`
	Enabled     bool                   `json:"enabled"`
}

type CredentialsConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

type RAIDConfig struct {
	Level        string   `json:"level"`
	Devices      []string `json:"devices"`
	SpareDevices []string `json:"spare_devices"`
	ChunkSize    int      `json:"chunk_size"`
	Enabled      bool     `json:"enabled"`
}

type TieringRuleConfig struct {
	Name        string        `json:"name"`
	Condition   string        `json:"condition"`
	SourceTier  string        `json:"source_tier"`
	TargetTier  string        `json:"target_tier"`
	Age         time.Duration `json:"age"`
	Size        int64         `json:"size"`
	AccessCount int           `json:"access_count"`
	Enabled     bool          `json:"enabled"`
}

type ReplicationConfig struct {
	Enabled     bool     `json:"enabled"`
	Targets     []string `json:"targets"`
	SyncMode    string   `json:"sync_mode"`
	Compression bool     `json:"compression"`
}

type EncryptionConfig struct {
	Enabled   bool              `json:"enabled"`
	Algorithm string            `json:"algorithm"`
	KeyID     string            `json:"key_id"`
	Options   map[string]string `json:"options"`
}

type SNMPConfig struct {
	Enabled       bool               `json:"enabled"`
	ListenAddress string             `json:"listen_address"`
	Port          int                `json:"port"`
	Community     string             `json:"community"`
	Version       string             `json:"version"`
	Users         []SNMPUserConfig   `json:"users"`
	TrapTargets   []TrapTargetConfig `json:"trap_targets"`
}

type SNMPUserConfig struct {
	Username     string `json:"username"`
	AuthProtocol string `json:"auth_protocol"`
	AuthKey      string `json:"auth_key"`
	PrivProtocol string `json:"priv_protocol"`
	PrivKey      string `json:"priv_key"`
}

type TrapTargetConfig struct {
	Address   string `json:"address"`
	Port      int    `json:"port"`
	Community string `json:"community"`
	Version   string `json:"version"`
}

type DatabaseConfig struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Options  map[string]interface{} `json:"options"`
}

type SecurityConfig struct {
	JWTSecret       string        `json:"jwt_secret"`
	SessionTimeout  time.Duration `json:"session_timeout"`
	PasswordPolicy  PasswordPolicy `json:"password_policy"`
	TwoFactorAuth   bool          `json:"two_factor_auth"`
	AuditLogging    bool          `json:"audit_logging"`
}

type PasswordPolicy struct {
	MinLength    int  `json:"min_length"`
	RequireUpper bool `json:"require_upper"`
	RequireLower bool `json:"require_lower"`
	RequireDigit bool `json:"require_digit"`
	RequireSymbol bool `json:"require_symbol"`
}

type LoggingConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

type RetentionConfig struct {
	DefaultTTL     time.Duration            `json:"default_ttl"`
	ChannelRules   map[string]time.Duration `json:"channel_rules"`
	QuarantineTTL  time.Duration            `json:"quarantine_ttl"`
	ProtectedMode  bool                     `json:"protected_mode"`
	SafetyThreshold float64                 `json:"safety_threshold"`
}

type ConfigWatcher interface {
	OnConfigChanged(config *Config) error
}

func NewManager(configPath string) *Manager {
	return &Manager{
		configPath: configPath,
		watchers:   make([]ConfigWatcher, 0),
	}
}

func (m *Manager) LoadConfig() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		m.config = m.getDefaultConfig()
		return m.saveConfigLocked()
	}
	
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	m.config = &config
	return nil
}

func (m *Manager) SaveConfig() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	return m.saveConfigLocked()
}

func (m *Manager) saveConfigLocked() error {
	if m.config == nil {
		return fmt.Errorf("no config to save")
	}
	
	m.config.UpdatedAt = time.Now()
	
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	for _, watcher := range m.watchers {
		if err := watcher.OnConfigChanged(m.config); err != nil {
			fmt.Printf("Config watcher error: %v\n", err)
		}
	}
	
	return nil
}

func (m *Manager) GetConfig() *Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	if m.config == nil {
		return nil
	}
	
	data, _ := json.Marshal(m.config)
	var copy Config
	json.Unmarshal(data, &copy)
	return &copy
}

func (m *Manager) UpdateConfig(config *Config) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if err := m.validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	
	m.config = config
	return m.saveConfigLocked()
}

func (m *Manager) validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}
	
	for name, backend := range config.Storage.Backends {
		if backend.Type == "" {
			return fmt.Errorf("backend %s has no type specified", name)
		}
		if backend.Endpoint == "" && backend.Type != "local" {
			return fmt.Errorf("backend %s has no endpoint specified", name)
		}
	}
	
	if config.SNMP.Enabled {
		if config.SNMP.Port <= 0 || config.SNMP.Port > 65535 {
			return fmt.Errorf("invalid SNMP port: %d", config.SNMP.Port)
		}
	}
	
	return nil
}

func (m *Manager) getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			TLS: TLSConfig{
				Enabled: false,
			},
		},
		Storage: StorageConfig{
			DefaultBackend: "local",
			Backends: map[string]BackendConfig{
				"local": {
					Type:     "local",
					Endpoint: "/data/minio",
					Enabled:  true,
					RAID: RAIDConfig{
						Level:   "1",
						Enabled: false,
					},
				},
			},
			Encryption: EncryptionConfig{
				Enabled:   true,
				Algorithm: "AES-256-GCM",
				KeyID:     "default",
			},
		},
		SNMP: SNMPConfig{
			Enabled:       true,
			ListenAddress: "0.0.0.0",
			Port:          161,
			Community:     "public",
			Version:       "2c",
		},
		Database: DatabaseConfig{
			Type: "sqlite",
			Path: "/data/skylogix.db",
		},
		Security: SecurityConfig{
			JWTSecret:      "change-me-in-production",
			SessionTimeout: 24 * time.Hour,
			PasswordPolicy: PasswordPolicy{
				MinLength:     8,
				RequireUpper:  true,
				RequireLower:  true,
				RequireDigit:  true,
				RequireSymbol: false,
			},
			AuditLogging: true,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
		Retention: RetentionConfig{
			DefaultTTL:      30 * 24 * time.Hour, // 30 days
			QuarantineTTL:   7 * 24 * time.Hour,  // 7 days
			ProtectedMode:   false,
			SafetyThreshold: 0.1, // 10%
		},
		UpdatedAt: time.Now(),
	}
}

func (m *Manager) AddWatcher(watcher ConfigWatcher) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.watchers = append(m.watchers, watcher)
}

func (m *Manager) RemoveWatcher(watcher ConfigWatcher) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	for i, w := range m.watchers {
		if w == watcher {
			m.watchers = append(m.watchers[:i], m.watchers[i+1:]...)
			break
		}
	}
}

func (m *Manager) BackupConfig() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	if m.config == nil {
		return fmt.Errorf("no config to backup")
	}
	
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup.%s", m.configPath, timestamp)
	
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config for backup: %w", err)
	}
	
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}
	
	fmt.Printf("Configuration backed up to: %s\n", backupPath)
	return nil
}

func (m *Manager) RestoreConfig(backupPath string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse backup file: %w", err)
	}
	
	return m.UpdateConfig(&config)
}
