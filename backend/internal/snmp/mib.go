package snmp

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type MIB struct {
	objects map[string]MIBValue
	mutex   sync.RWMutex
}

type MIBValue struct {
	Type        string      `json:"type"`
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
	Access      string      `json:"access"`
}

func NewMIB() *MIB {
	mib := &MIB{
		objects: make(map[string]MIBValue),
	}
	
	mib.registerSystemObjects()
	
	return mib
}

func (m *MIB) registerSystemObjects() {
	m.Register("1.3.6.1.2.1.1.1.0", MIBValue{
		Type:        "OctetString",
		Value:       "Skylogix Storage System v1.0",
		Description: "A textual description of the entity",
		Access:      "read-only",
	})
	
	m.Register("1.3.6.1.2.1.1.2.0", MIBValue{
		Type:        "ObjectIdentifier",
		Value:       "1.3.6.1.4.1.12345",
		Description: "The vendor's authoritative identification",
		Access:      "read-only",
	})
	
	m.Register("1.3.6.1.2.1.1.3.0", MIBValue{
		Type:        "TimeTicks",
		Value:       uint32(0),
		Description: "Time since the network management portion of the system was last re-initialized",
		Access:      "read-only",
	})
	
	m.Register("1.3.6.1.2.1.1.4.0", MIBValue{
		Type:        "OctetString",
		Value:       "admin@skylogix.com",
		Description: "The textual identification of the contact person",
		Access:      "read-write",
	})
	
	m.Register("1.3.6.1.2.1.1.5.0", MIBValue{
		Type:        "OctetString",
		Value:       "skylogix-storage",
		Description: "An administratively-assigned name for this managed node",
		Access:      "read-write",
	})
	
	m.Register("1.3.6.1.2.1.1.6.0", MIBValue{
		Type:        "OctetString",
		Value:       "Data Center",
		Description: "The physical location of this node",
		Access:      "read-write",
	})
}

func (m *MIB) Register(oid string, value MIBValue) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.objects[oid] = value
}

func (m *MIB) GetValue(oid string) (MIBValue, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	value, exists := m.objects[oid]
	if !exists {
		return MIBValue{}, fmt.Errorf("OID not found: %s", oid)
	}
	
	return value, nil
}

func (m *MIB) SetValue(oid string, value MIBValue) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	existing, exists := m.objects[oid]
	if !exists {
		return fmt.Errorf("OID not found: %s", oid)
	}
	
	if existing.Access == "read-only" {
		return fmt.Errorf("OID is read-only: %s", oid)
	}
	
	if value.Description == "" {
		value.Description = existing.Description
	}
	if value.Access == "" {
		value.Access = existing.Access
	}
	
	m.objects[oid] = value
	return nil
}

func (m *MIB) GetNextValue(oid string) (string, MIBValue, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	oids := make([]string, 0, len(m.objects))
	for k := range m.objects {
		oids = append(oids, k)
	}
	sort.Strings(oids)
	
	for _, nextOID := range oids {
		if compareOIDs(oid, nextOID) < 0 {
			return nextOID, m.objects[nextOID], nil
		}
	}
	
	return "", MIBValue{}, fmt.Errorf("no next OID found")
}

func compareOIDs(oid1, oid2 string) int {
	parts1 := strings.Split(oid1, ".")
	parts2 := strings.Split(oid2, ".")
	
	minLen := len(parts1)
	if len(parts2) < minLen {
		minLen = len(parts2)
	}
	
	for i := 0; i < minLen; i++ {
		if parts1[i] < parts2[i] {
			return -1
		} else if parts1[i] > parts2[i] {
			return 1
		}
	}
	
	if len(parts1) < len(parts2) {
		return -1
	} else if len(parts1) > len(parts2) {
		return 1
	}
	
	return 0
}

func (m *MIB) Walk(startOID string, callback func(oid string, value MIBValue) bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	oids := make([]string, 0, len(m.objects))
	for k := range m.objects {
		oids = append(oids, k)
	}
	sort.Strings(oids)
	
	for _, oid := range oids {
		if compareOIDs(startOID, oid) <= 0 {
			if !callback(oid, m.objects[oid]) {
				break
			}
		}
	}
}

func (m *MIB) LoadFromFile(filename string) error {
	fmt.Printf("Loading MIB from file: %s\n", filename)
	return nil
}

func (m *MIB) GetSubtree(rootOID string) map[string]MIBValue {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	subtree := make(map[string]MIBValue)
	
	for oid, value := range m.objects {
		if strings.HasPrefix(oid, rootOID) {
			subtree[oid] = value
		}
	}
	
	return subtree
}

func (m *MIB) GetStatistics() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	stats := map[string]interface{}{
		"total_objects":    len(m.objects),
		"read_only_count":  0,
		"read_write_count": 0,
	}
	
	for _, value := range m.objects {
		if value.Access == "read-only" {
			stats["read_only_count"] = stats["read_only_count"].(int) + 1
		} else if value.Access == "read-write" {
			stats["read_write_count"] = stats["read_write_count"].(int) + 1
		}
	}
	
	return stats
}
