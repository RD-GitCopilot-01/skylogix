package storage

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type RAIDManager struct {
	config RAIDConfig
}

func NewRAIDManager() *RAIDManager {
	return &RAIDManager{}
}

func (rm *RAIDManager) Configure(config RAIDConfig) error {
	rm.config = config
	
	if err := rm.validateConfig(); err != nil {
		return fmt.Errorf("invalid RAID configuration: %w", err)
	}
	
	return nil
}

func (rm *RAIDManager) validateConfig() error {
	if len(rm.config.Devices) < 2 {
		return fmt.Errorf("RAID requires at least 2 devices")
	}
	
	switch rm.config.Level {
	case RAID0:
		if len(rm.config.Devices) < 2 {
			return fmt.Errorf("RAID 0 requires at least 2 devices")
		}
	case RAID1:
		if len(rm.config.Devices) != 2 {
			return fmt.Errorf("RAID 1 requires exactly 2 devices")
		}
	case RAID5:
		if len(rm.config.Devices) < 3 {
			return fmt.Errorf("RAID 5 requires at least 3 devices")
		}
	case RAID6:
		if len(rm.config.Devices) < 4 {
			return fmt.Errorf("RAID 6 requires at least 4 devices")
		}
	case RAID10:
		if len(rm.config.Devices) < 4 || len(rm.config.Devices)%2 != 0 {
			return fmt.Errorf("RAID 10 requires at least 4 devices and even number")
		}
	default:
		return fmt.Errorf("unsupported RAID level: %s", rm.config.Level)
	}
	
	return nil
}

func (rm *RAIDManager) CreateArray(arrayName string) error {
	args := []string{
		"--create",
		arrayName,
		"--level=" + string(rm.config.Level),
		"--raid-devices=" + strconv.Itoa(len(rm.config.Devices)),
	}
	
	if rm.config.ChunkSize > 0 {
		args = append(args, "--chunk="+strconv.Itoa(rm.config.ChunkSize))
	}
	
	if len(rm.config.SpareDevices) > 0 {
		args = append(args, "--spare-devices="+strconv.Itoa(len(rm.config.SpareDevices)))
	}
	
	args = append(args, rm.config.Devices...)
	args = append(args, rm.config.SpareDevices...)
	
	cmd := exec.Command("mdadm", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create RAID array: %s, output: %s", err, string(output))
	}
	
	return nil
}

func (rm *RAIDManager) GetHealth() RAIDHealthStatus {
	status := RAIDHealthStatus{
		Level:       rm.config.Level,
		ActiveDisks: len(rm.config.Devices),
		FailedDisks: 0,
		SpareDisks:  len(rm.config.SpareDevices),
		State:       "unknown",
		LastSync:    time.Now(),
	}
	
	if mdstat, err := rm.readMdstat(); err == nil {
		status = rm.parseMdstat(mdstat, status)
	}
	
	return status
}

func (rm *RAIDManager) readMdstat() (string, error) {
	cmd := exec.Command("cat", "/proc/mdstat")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (rm *RAIDManager) parseMdstat(mdstat string, status RAIDHealthStatus) RAIDHealthStatus {
	lines := strings.Split(mdstat, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.Contains(line, "raid") {
			if strings.Contains(line, "active") {
				status.State = "active"
			} else if strings.Contains(line, "inactive") {
				status.State = "inactive"
			} else if strings.Contains(line, "degraded") {
				status.State = "degraded"
			}
			
			if strings.Contains(line, "[") && strings.Contains(line, "]") {
				start := strings.Index(line, "[")
				end := strings.Index(line, "]")
				if start != -1 && end != -1 && end > start {
					deviceStatus := line[start+1 : end]
					status.ActiveDisks = strings.Count(deviceStatus, "U")
					status.FailedDisks = strings.Count(deviceStatus, "_")
				}
			}
		}
		
		if strings.Contains(line, "resync") || strings.Contains(line, "recovery") {
			status.SyncStatus = "syncing"
		} else if strings.Contains(line, "check") {
			status.SyncStatus = "checking"
		} else {
			status.SyncStatus = "idle"
		}
	}
	
	return status
}

func (rm *RAIDManager) RepairArray(arrayName string, newDevice string) error {
	cmd := exec.Command("mdadm", "--manage", arrayName, "--add", newDevice)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add device to RAID array: %s, output: %s", err, string(output))
	}
	
	return nil
}

func (rm *RAIDManager) RemoveFailedDevice(arrayName string, device string) error {
	cmd := exec.Command("mdadm", "--manage", arrayName, "--fail", device)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to mark device as failed: %s, output: %s", err, string(output))
	}
	
	cmd = exec.Command("mdadm", "--manage", arrayName, "--remove", device)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove device from RAID array: %s, output: %s", err, string(output))
	}
	
	return nil
}

func (rm *RAIDManager) GetArrayInfo(arrayName string) (map[string]interface{}, error) {
	cmd := exec.Command("mdadm", "--detail", arrayName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get array details: %w", err)
	}
	
	info := make(map[string]interface{})
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				info[key] = value
			}
		}
	}
	
	return info, nil
}

func (rm *RAIDManager) Close() error {
	return nil
}
