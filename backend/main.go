package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"skylogix-storage/internal/config"
	"skylogix-storage/internal/database"
	"skylogix-storage/internal/snmp"
	"skylogix-storage/internal/storage"
	"skylogix-storage/internal/web"
)

func main() {
	fmt.Println("Starting Skylogix Storage Management System...")

	configManager := config.NewManager("/etc/skylogix/config.json")
	if err := configManager.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cfg := configManager.GetConfig()
	if cfg == nil {
		log.Fatal("Configuration is nil")
	}

	dbService, err := database.NewService(database.Config{
		Type: cfg.Database.Type,
		Path: cfg.Database.Path,
		Host: cfg.Database.Host,
		Port: cfg.Database.Port,
	})
	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}

	if err := dbService.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbService.Close()

	storageManager := storage.NewManager()
	
	managerConfig := storage.ManagerConfig{
		DefaultBackend: cfg.Storage.DefaultBackend,
		Monitoring: storage.MonitoringConfig{
			Enabled:       true,
			CheckInterval: 30 * time.Second,
		},
	}
	if err := storageManager.Configure(managerConfig); err != nil {
		log.Fatalf("Failed to configure storage manager: %v", err)
	}

	for name, backendCfg := range cfg.Storage.Backends {
		if !backendCfg.Enabled {
			continue
		}

		var backend storage.StorageBackend
		switch backendCfg.Type {
		case "local":
			backend = storage.NewLocalStorage()
		case "nas":
			backend = storage.NewNASStorage()
		case "san":
			backend = storage.NewSANStorage()
		case "object":
			backend = storage.NewMinIOClient()
		default:
			log.Printf("Unknown storage backend type: %s", backendCfg.Type)
			continue
		}

		storageCfg := storage.StorageConfig{
			Type:     storage.StorageType(backendCfg.Type),
			Endpoint: backendCfg.Endpoint,
			Credentials: storage.StorageCredentials{
				AccessKey: backendCfg.Credentials.AccessKey,
				SecretKey: backendCfg.Credentials.SecretKey,
				Username:  backendCfg.Credentials.Username,
				Password:  backendCfg.Credentials.Password,
			},
			RAID: storage.RAIDConfig{
				Level:   storage.RAIDLevel(backendCfg.RAID.Level),
				Devices: backendCfg.RAID.Devices,
				Enabled: backendCfg.RAID.Enabled,
			},
			Encryption: storage.EncryptionConfig{
				Enabled:   cfg.Storage.Encryption.Enabled,
				Algorithm: cfg.Storage.Encryption.Algorithm,
				KeyID:     cfg.Storage.Encryption.KeyID,
			},
			Options: backendCfg.Options,
		}

		if err := backend.Configure(storageCfg); err != nil {
			log.Printf("Failed to configure backend %s: %v", name, err)
			continue
		}

		if err := storageManager.AddBackend(name, backend); err != nil {
			log.Printf("Failed to add backend %s: %v", name, err)
			continue
		}

		fmt.Printf("Added storage backend: %s (%s)\n", name, backendCfg.Type)
	}

	storageManager.StartMonitoring()

	snmpAgent := snmp.NewAgent()
	if cfg.SNMP.Enabled {
		snmpConfig := snmp.Config{
			ListenAddress: cfg.SNMP.ListenAddress,
			Port:          cfg.SNMP.Port,
			Community:     cfg.SNMP.Community,
			Version:       cfg.SNMP.Version,
		}

		if err := snmpAgent.Configure(snmpConfig); err != nil {
			log.Printf("Failed to configure SNMP agent: %v", err)
		} else {
			snmpAgent.RegisterStorageMetrics()

			if err := snmpAgent.Start(); err != nil {
				log.Printf("Failed to start SNMP agent: %v", err)
			} else {
				fmt.Printf("SNMP agent started on %s:%d\n", cfg.SNMP.ListenAddress, cfg.SNMP.Port)
			}
		}
	}

	webConfig := web.Config{
		Host:         cfg.Server.Host,
		Port:         cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		TLS: web.TLSConfig{
			Enabled:  cfg.Server.TLS.Enabled,
			CertFile: cfg.Server.TLS.CertFile,
			KeyFile:  cfg.Server.TLS.KeyFile,
		},
		CORS: web.CORSConfig{
			Enabled:        true,
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
	}

	webServer := web.NewServer(webConfig)

	handlers := web.NewHandlers(
		&StorageManagerAdapter{storageManager},
		&ConfigManagerAdapter{configManager},
		&SNMPAgentAdapter{snmpAgent},
		&DatabaseServiceAdapter{dbService},
	)

	webServer.SetupRoutes(handlers)

	go func() {
		if err := webServer.Start(); err != nil {
			log.Printf("Web server error: %v", err)
		}
	}()

	fmt.Printf("Web server started on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Println("Skylogix Storage Management System is running...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := webServer.Stop(ctx); err != nil {
		log.Printf("Error stopping web server: %v", err)
	}

	if err := snmpAgent.Stop(); err != nil {
		log.Printf("Error stopping SNMP agent: %v", err)
	}

	if err := storageManager.Close(); err != nil {
		log.Printf("Error closing storage manager: %v", err)
	}

	fmt.Println("Shutdown complete")
}


type StorageManagerAdapter struct {
	*storage.Manager
}

func (s *StorageManagerAdapter) GetBackends() map[string]interface{} {
	backends := make(map[string]interface{})
	for _, name := range s.Manager.ListBackends() {
		if backend, err := s.Manager.GetBackend(name); err == nil {
			backends[name] = map[string]interface{}{
				"type":   "unknown",
				"health": backend.GetHealth(),
			}
		}
	}
	return backends
}

func (s *StorageManagerAdapter) GetBackend(name string) (interface{}, error) {
	backend, err := s.Manager.GetBackend(name)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name":    name,
		"health":  backend.GetHealth(),
		"metrics": backend.GetMetrics(),
	}, nil
}

func (s *StorageManagerAdapter) AddBackend(name string, config interface{}) error {
	return fmt.Errorf("not implemented")
}

func (s *StorageManagerAdapter) GetHealth() interface{} {
	return s.Manager.GetAggregatedHealth()
}

func (s *StorageManagerAdapter) GetMetrics() interface{} {
	return s.Manager.GetAggregatedMetrics()
}

type ConfigManagerAdapter struct {
	*config.Manager
}

func (c *ConfigManagerAdapter) GetConfig() interface{} {
	return c.Manager.GetConfig()
}

func (c *ConfigManagerAdapter) UpdateConfig(cfg interface{}) error {
	return fmt.Errorf("not implemented")
}

func (c *ConfigManagerAdapter) BackupConfig() error {
	return c.Manager.BackupConfig()
}

func (c *ConfigManagerAdapter) RestoreConfig(path string) error {
	return c.Manager.RestoreConfig(path)
}

func (c *ConfigManagerAdapter) GetHistory() ([]interface{}, error) {
	return []interface{}{}, nil
}

type SNMPAgentAdapter struct {
	*snmp.Agent
}

func (s *SNMPAgentAdapter) IsRunning() bool {
	return true
}

func (s *SNMPAgentAdapter) GetMIBObjects() map[string]interface{} {
	return map[string]interface{}{}
}

func (s *SNMPAgentAdapter) SendTrap(oid string, vars []interface{}) error {
	return nil
}

type DatabaseServiceAdapter struct {
	*database.Service
}

func (d *DatabaseServiceAdapter) GetStats() interface{} {
	return d.Service.GetStats()
}

func (d *DatabaseServiceAdapter) Backup(path string) error {
	return d.Service.BackupDatabase(path)
}

func (d *DatabaseServiceAdapter) Restore(path string) error {
	return d.Service.RestoreDatabase(path)
}

func (d *DatabaseServiceAdapter) Cleanup() error {
	return d.Service.CleanupOldRecords()
}

func (d *DatabaseServiceAdapter) GetAlerts() ([]interface{}, error) {
	alerts, err := d.Service.GetAlerts(false, 100)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(alerts))
	for i, alert := range alerts {
		result[i] = alert
	}
	return result, nil
}

func (d *DatabaseServiceAdapter) GetAuditLogs(userID string, since time.Time) ([]interface{}, error) {
	logs, err := d.Service.GetAuditLogs(userID, since, 100)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(logs))
	for i, log := range logs {
		result[i] = log
	}
	return result, nil
}
