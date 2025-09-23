package snmp

import (
	"fmt"
	"net"
	"time"
)

type Agent struct {
	config   Config
	listener net.PacketConn
	mib      *MIB
	running  bool
}

type Config struct {
	ListenAddress string            `json:"listen_address"`
	Port          int               `json:"port"`
	Community     string            `json:"community"`
	Version       string            `json:"version"`
	Users         []SNMPUser        `json:"users"`
	TrapTargets   []TrapTarget      `json:"trap_targets"`
	MIBPath       string            `json:"mib_path"`
	Options       map[string]string `json:"options"`
}

type SNMPUser struct {
	Username     string `json:"username"`
	AuthProtocol string `json:"auth_protocol"`
	AuthKey      string `json:"auth_key"`
	PrivProtocol string `json:"priv_protocol"`
	PrivKey      string `json:"priv_key"`
}

type TrapTarget struct {
	Address   string `json:"address"`
	Port      int    `json:"port"`
	Community string `json:"community"`
	Version   string `json:"version"`
}

func NewAgent() *Agent {
	return &Agent{
		config: Config{
			ListenAddress: "0.0.0.0",
			Port:          161,
			Community:     "public",
			Version:       "2c",
		},
		mib: NewMIB(),
	}
}

func (a *Agent) Configure(config Config) error {
	a.config = config
	
	if config.MIBPath != "" {
		if err := a.mib.LoadFromFile(config.MIBPath); err != nil {
			return fmt.Errorf("failed to load MIB: %w", err)
		}
	}
	
	return nil
}

func (a *Agent) Start() error {
	if a.running {
		return fmt.Errorf("agent already running")
	}
	
	addr := fmt.Sprintf("%s:%d", a.config.ListenAddress, a.config.Port)
	listener, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	
	a.listener = listener
	a.running = true
	
	go a.handleRequests()
	
	fmt.Printf("SNMP agent started on %s\n", addr)
	return nil
}

func (a *Agent) Stop() error {
	if !a.running {
		return nil
	}
	
	a.running = false
	
	if a.listener != nil {
		return a.listener.Close()
	}
	
	return nil
}

func (a *Agent) handleRequests() {
	buffer := make([]byte, 1500)
	
	for a.running {
		n, addr, err := a.listener.ReadFrom(buffer)
		if err != nil {
			if a.running {
				fmt.Printf("Error reading SNMP request: %v\n", err)
			}
			continue
		}
		
		go a.processRequest(buffer[:n], addr)
	}
}

func (a *Agent) processRequest(data []byte, addr net.Addr) {
	message, err := a.parseMessage(data)
	if err != nil {
		fmt.Printf("Failed to parse SNMP message from %s: %v\n", addr, err)
		return
	}
	
	var response []byte
	switch message.PDUType {
	case "GetRequest":
		response, err = a.handleGetRequest(message)
	case "GetNextRequest":
		response, err = a.handleGetNextRequest(message)
	case "GetBulkRequest":
		response, err = a.handleGetBulkRequest(message)
	case "SetRequest":
		response, err = a.handleSetRequest(message)
	default:
		fmt.Printf("Unsupported PDU type: %s\n", message.PDUType)
		return
	}
	
	if err != nil {
		fmt.Printf("Error handling SNMP request: %v\n", err)
		return
	}
	
	if response != nil {
		_, err = a.listener.WriteTo(response, addr)
		if err != nil {
			fmt.Printf("Error sending SNMP response: %v\n", err)
		}
	}
}

type SNMPMessage struct {
	Version   int
	Community string
	PDUType   string
	RequestID int
	VarBinds  []VarBind
}

type VarBind struct {
	OID   string
	Type  string
	Value interface{}
}

func (a *Agent) parseMessage(data []byte) (*SNMPMessage, error) {
	
	message := &SNMPMessage{
		Version:   2, // SNMPv2c
		Community: a.config.Community,
		PDUType:   "GetRequest",
		RequestID: 1,
		VarBinds: []VarBind{
			{
				OID:  "1.3.6.1.4.1.12345.1.1.1.0", // Example OID
				Type: "OctetString",
			},
		},
	}
	
	return message, nil
}

func (a *Agent) handleGetRequest(message *SNMPMessage) ([]byte, error) {
	response := &SNMPMessage{
		Version:   message.Version,
		Community: message.Community,
		PDUType:   "GetResponse",
		RequestID: message.RequestID,
		VarBinds:  make([]VarBind, len(message.VarBinds)),
	}
	
	for i, vb := range message.VarBinds {
		value, err := a.mib.GetValue(vb.OID)
		if err != nil {
			response.VarBinds[i] = VarBind{
				OID:   vb.OID,
				Type:  "NoSuchObject",
				Value: nil,
			}
		} else {
			response.VarBinds[i] = VarBind{
				OID:   vb.OID,
				Type:  value.Type,
				Value: value.Value,
			}
		}
	}
	
	return a.encodeMessage(response)
}

func (a *Agent) handleGetNextRequest(message *SNMPMessage) ([]byte, error) {
	response := &SNMPMessage{
		Version:   message.Version,
		Community: message.Community,
		PDUType:   "GetResponse",
		RequestID: message.RequestID,
		VarBinds:  make([]VarBind, len(message.VarBinds)),
	}
	
	for i, vb := range message.VarBinds {
		nextOID, value, err := a.mib.GetNextValue(vb.OID)
		if err != nil {
			response.VarBinds[i] = VarBind{
				OID:   vb.OID,
				Type:  "EndOfMibView",
				Value: nil,
			}
		} else {
			response.VarBinds[i] = VarBind{
				OID:   nextOID,
				Type:  value.Type,
				Value: value.Value,
			}
		}
	}
	
	return a.encodeMessage(response)
}

func (a *Agent) handleGetBulkRequest(message *SNMPMessage) ([]byte, error) {
	return a.handleGetNextRequest(message)
}

func (a *Agent) handleSetRequest(message *SNMPMessage) ([]byte, error) {
	response := &SNMPMessage{
		Version:   message.Version,
		Community: message.Community,
		PDUType:   "GetResponse",
		RequestID: message.RequestID,
		VarBinds:  make([]VarBind, len(message.VarBinds)),
	}
	
	for i, vb := range message.VarBinds {
		err := a.mib.SetValue(vb.OID, MIBValue{
			Type:  vb.Type,
			Value: vb.Value,
		})
		
		if err != nil {
			response.VarBinds[i] = VarBind{
				OID:   vb.OID,
				Type:  "NoAccess",
				Value: nil,
			}
		} else {
			response.VarBinds[i] = vb
		}
	}
	
	return a.encodeMessage(response)
}

func (a *Agent) encodeMessage(message *SNMPMessage) ([]byte, error) {
	
	return []byte("SNMP response"), nil
}

func (a *Agent) SendTrap(trapOID string, varBinds []VarBind) error {
	for _, target := range a.config.TrapTargets {
		if err := a.sendTrapToTarget(target, trapOID, varBinds); err != nil {
			fmt.Printf("Failed to send trap to %s:%d: %v\n", target.Address, target.Port, err)
		}
	}
	
	return nil
}

func (a *Agent) sendTrapToTarget(target TrapTarget, trapOID string, varBinds []VarBind) error {
	trap := &SNMPMessage{
		Version:   2, // SNMPv2c
		Community: target.Community,
		PDUType:   "SNMPv2-Trap",
		RequestID: int(time.Now().Unix()),
		VarBinds:  varBinds,
	}
	
	data, err := a.encodeMessage(trap)
	if err != nil {
		return err
	}
	
	addr := fmt.Sprintf("%s:%d", target.Address, target.Port)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	
	_, err = conn.Write(data)
	return err
}

func (a *Agent) RegisterStorageMetrics() {
	a.mib.Register("1.3.6.1.4.1.12345.1.1.1.0", MIBValue{
		Type:        "OctetString",
		Value:       "Skylogix Storage System",
		Description: "System description",
	})
	
	a.mib.Register("1.3.6.1.4.1.12345.1.1.2.0", MIBValue{
		Type:        "Counter64",
		Value:       uint64(0),
		Description: "Total storage space in bytes",
	})
	
	a.mib.Register("1.3.6.1.4.1.12345.1.1.3.0", MIBValue{
		Type:        "Counter64",
		Value:       uint64(0),
		Description: "Used storage space in bytes",
	})
	
	a.mib.Register("1.3.6.1.4.1.12345.1.1.4.0", MIBValue{
		Type:        "Counter64",
		Value:       uint64(0),
		Description: "Free storage space in bytes",
	})
	
	a.mib.Register("1.3.6.1.4.1.12345.1.1.5.0", MIBValue{
		Type:        "Integer32",
		Value:       int32(1),
		Description: "Storage system health status",
	})
}

func (a *Agent) UpdateStorageMetrics(totalSpace, usedSpace, freeSpace uint64, healthStatus int32) {
	a.mib.SetValue("1.3.6.1.4.1.12345.1.1.2.0", MIBValue{
		Type:  "Counter64",
		Value: totalSpace,
	})
	
	a.mib.SetValue("1.3.6.1.4.1.12345.1.1.3.0", MIBValue{
		Type:  "Counter64",
		Value: usedSpace,
	})
	
	a.mib.SetValue("1.3.6.1.4.1.12345.1.1.4.0", MIBValue{
		Type:  "Counter64",
		Value: freeSpace,
	})
	
	a.mib.SetValue("1.3.6.1.4.1.12345.1.1.5.0", MIBValue{
		Type:  "Integer32",
		Value: healthStatus,
	})
}
