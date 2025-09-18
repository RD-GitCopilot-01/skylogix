# Product Backlog - Skylogix SDR VHF/UHF System

## Tổng quan Dự án (Project Overview)

Dự án Skylogix là hệ thống Software Defined Radio (SDR) VHF/UHF dựa trên platform XCZU-2EG + AD9364, được thiết kế để cung cấp các chức năng radio chuyên nghiệp cho ứng dụng hàng không và hàng hải.

### Thông tin Dự án
- **Tên dự án**: Skylogix SDR VHF/UHF System
- **Platform**: XCZU-2EG MPSoC + AD9364 RF Transceiver
- **Băng tần**: VHF (108-174 MHz), UHF (225-512 MHz)
- **Chuẩn tuân thủ**: ED-137, ICAO Annex 10, ITU-R recommendations

### Kiến trúc Phần mềm Tổng thể
```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │ Web Management  │  │   SNMP Agent    │  │ Radio Control│ │
│  │   Interface     │  │                 │  │ Application │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                    Middleware Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   VoIP Stack    │  │ Message Broker  │  │  Database   │ │
│  │   (ED-137)      │  │                 │  │ Services    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                    OS Kernel Layer                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │ Linux RT Kernel │  │ Network Stack   │  │   Drivers   │ │
│  │   (5.15 LTS)    │  │                 │  │             │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                    Hardware Layer                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   XCZU-2EG      │  │    AD9364       │  │    FPGA     │ │
│  │    MPSoC        │  │ RF Transceiver  │  │ Processing  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## EPIC 1: Core Radio Software Development

### Epic Description
Phát triển phần mềm radio cốt lõi bao gồm các chức năng cơ bản của transceiver, xử lý tín hiệu số, và giao tiếp với hardware.

### User Stories

#### US1.1: RF Transceiver Control
**Là một** operator  
**Tôi muốn** điều khiển các thông số RF cơ bản (tần số, công suất, mode)  
**Để** có thể vận hành radio theo yêu cầu

**Acceptance Criteria:**
- [ ] Có thể thiết lập tần số trong dải VHF (108-174 MHz) và UHF (225-512 MHz)
- [ ] Điều chỉnh công suất phát từ 1W đến 50W
- [ ] Hỗ trợ các mode: AM, FM, USB, LSB
- [ ] Thời gian chuyển đổi tần số < 100ms
- [ ] Độ chính xác tần số ±1 ppm

#### US1.2: Channel Management System
**Là một** operator  
**Tôi muốn** quản lý danh sách kênh được lập trình sẵn  
**Để** có thể chuyển đổi nhanh giữa các kênh thường dùng

**Acceptance Criteria:**
- [ ] Lưu trữ tối thiểu 999 kênh
- [ ] Mỗi kênh bao gồm: tần số, mode, công suất, tên, mô tả
- [ ] Chuyển đổi kênh trong < 50ms
- [ ] Import/export danh sách kênh qua file CSV
- [ ] Tìm kiếm kênh theo tên hoặc tần số

#### US1.3: PTT (Push-to-Talk) Control
**Là một** operator  
**Tôi muốn** điều khiển PTT để chuyển đổi giữa chế độ thu và phát  
**Để** có thể giao tiếp radio

**Acceptance Criteria:**
- [ ] Phát hiện tín hiệu PTT trong < 1ms
- [ ] Chuyển đổi T/R trong < 10ms
- [ ] Hỗ trợ PTT qua hardware pin và software command
- [ ] Hiển thị trạng thái TX/RX trên giao diện
- [ ] Timeout protection (tự động ngắt sau 5 phút)

#### US1.4: Audio Processing
**Là một** operator  
**Tôi muốn** xử lý audio với chất lượng cao  
**Để** đảm bảo chất lượng giao tiếp

**Acceptance Criteria:**
- [ ] Sample rate: 8kHz, 16kHz, 48kHz
- [ ] Bit depth: 16-bit, 24-bit
- [ ] AGC (Automatic Gain Control) với thời gian phản hồi < 100ms
- [ ] Noise reduction và echo cancellation
- [ ] Audio level monitoring và adjustment

---

## EPIC 2: VoIP Integration (ED-137)

### Epic Description
Tích hợp stack VoIP tuân thủ chuẩn ED-137 cho ứng dụng Air Traffic Control, bao gồm SIP signaling và RTP media streaming.

### User Stories

#### US2.1: SIP Protocol Implementation
**Là một** system administrator  
**Tôi muốn** cấu hình SIP client tuân thủ ED-137  
**Để** tích hợp với hệ thống ATC hiện có

**Acceptance Criteria:**
- [ ] SIP client tuân thủ RFC 3261 và ED-137
- [ ] Hỗ trợ SIP registration, INVITE, BYE, CANCEL
- [ ] Authentication qua SIP digest
- [ ] Heartbeat mechanism cho connection monitoring
- [ ] Failover tự động khi mất kết nối

#### US2.2: RTP Media Streaming
**Là một** operator  
**Tôi muốn** truyền/nhận audio qua RTP  
**Để** thực hiện giao tiếp VoIP chất lượng cao

**Acceptance Criteria:**
- [ ] RTP payload types: PCMA (G.711a), PCMU (G.711u)
- [ ] Jitter buffer adaptive (20-90ms)
- [ ] Packet loss detection và concealment
- [ ] RTCP reports cho quality monitoring
- [ ] End-to-end latency < 150ms

#### US2.3: Multi-frequency Operation
**Là một** ATC operator  
**Tôi muốn** vận hành đồng thời nhiều tần số  
**Để** quản lý nhiều sector cùng lúc

**Acceptance Criteria:**
- [ ] Hỗ trợ tối thiểu 4 kênh đồng thời
- [ ] Independent volume control cho mỗi kênh
- [ ] Priority-based audio mixing
- [ ] Individual mute/unmute controls
- [ ] Cross-channel interference < -60dB

---

## EPIC 3: Web Management Interface

### Epic Description
Phát triển giao diện web hiện đại để quản lý, giám sát và cấu hình hệ thống SDR từ xa.

### User Stories

#### US3.1: System Dashboard
**Là một** administrator  
**Tôi muốn** xem tổng quan trạng thái hệ thống  
**Để** giám sát hoạt động real-time

**Acceptance Criteria:**
- [ ] Hiển thị trạng thái RF (frequency, power, VSWR)
- [ ] System metrics (CPU, memory, temperature)
- [ ] Network status và throughput
- [ ] Active connections và call statistics
- [ ] Real-time updates qua WebSocket

#### US3.2: Configuration Management
**Là một** administrator  
**Tôi muốn** cấu hình hệ thống qua web interface  
**Để** không cần truy cập trực tiếp thiết bị

**Acceptance Criteria:**
- [ ] RF parameters configuration
- [ ] Network settings (IP, VLAN, routing)
- [ ] VoIP settings (SIP server, codecs)
- [ ] User management và permissions
- [ ] Configuration backup/restore

#### US3.3: Monitoring và Logging
**Là một** operator  
**Tôi muốn** xem logs và alerts  
**Để** troubleshoot các vấn đề

**Acceptance Criteria:**
- [ ] Real-time log viewer với filtering
- [ ] Alert notifications (email, SNMP traps)
- [ ] Historical data charts
- [ ] Export logs to CSV/PDF
- [ ] Log rotation và archiving

#### US3.4: Remote Control Interface
**Là một** operator  
**Tôi muốn** điều khiển radio từ xa qua web  
**Để** vận hành không cần có mặt tại thiết bị

**Acceptance Criteria:**
- [ ] Virtual PTT button
- [ ] Frequency và channel selection
- [ ] Volume controls
- [ ] Mode selection (AM/FM/USB/LSB)
- [ ] Power level adjustment

---

## EPIC 4: SNMP Management Agent

### Epic Description
Triển khai SNMP agent tuân thủ chuẩn để tích hợp với các hệ thống network management hiện có.

### User Stories

#### US4.1: SNMP v3 Implementation
**Là một** network administrator  
**Tôi muốn** quản lý thiết bị qua SNMP v3  
**Để** tích hợp với NMS hiện có

**Acceptance Criteria:**
- [ ] SNMPv3 với authentication và encryption
- [ ] Support cho GET, SET, GETNEXT, GETBULK
- [ ] User-based security model (USM)
- [ ] View-based access control (VACM)
- [ ] MIB-II standard compliance

#### US4.2: Custom MIB Development
**Là một** system integrator  
**Tôi muốn** truy cập các thông số đặc thù của SDR  
**Để** giám sát chi tiết hệ thống

**Acceptance Criteria:**
- [ ] Custom MIB cho RF parameters
- [ ] VoIP call statistics MIB objects
- [ ] System health indicators
- [ ] Configuration parameters access
- [ ] Trap definitions cho critical events

#### US4.3: Trap Generation
**Là một** NOC operator  
**Tôi muốn** nhận alerts tự động  
**Để** phản ứng nhanh với các sự cố

**Acceptance Criteria:**
- [ ] Critical alarms: over-temperature, high VSWR
- [ ] Warning notifications: high CPU, memory usage
- [ ] Informational traps: configuration changes
- [ ] Configurable trap destinations
- [ ] Trap rate limiting

---

## EPIC 5: Database Services

### Epic Description
Thiết kế và triển khai hệ thống database để lưu trữ cấu hình, logs, và dữ liệu vận hành.

### User Stories

#### US5.1: Configuration Database
**Là một** system  
**Tôi muốn** lưu trữ cấu hình persistent  
**Để** khôi phục sau khi reboot

**Acceptance Criteria:**
- [ ] SQLite database cho embedded system
- [ ] Configuration versioning
- [ ] Atomic transactions cho config updates
- [ ] Database backup/restore
- [ ] Migration scripts cho schema updates

#### US5.2: Operational Data Storage
**Là một** administrator  
**Tôi muốn** lưu trữ dữ liệu vận hành  
**Để** phân tích và báo cáo

**Acceptance Criteria:**
- [ ] Call detail records (CDR)
- [ ] Performance metrics history
- [ ] Event logs với timestamp
- [ ] Data retention policies
- [ ] Efficient indexing cho queries

#### US5.3: Statistics và Reporting
**Là một** manager  
**Tôi muốn** tạo báo cáo thống kê  
**Để** đánh giá hiệu suất hệ thống

**Acceptance Criteria:**
- [ ] Daily/weekly/monthly reports
- [ ] Usage statistics by channel
- [ ] System uptime reports
- [ ] Performance trend analysis
- [ ] Export to PDF/Excel formats

---

## EPIC 6: FPGA Signal Processing

### Epic Description
Phát triển các module xử lý tín hiệu trong FPGA để tăng tốc các thuật toán DSP quan trọng.

### User Stories

#### US6.1: DDC/DUC Implementation
**Là một** system  
**Tôi muốn** xử lý tín hiệu hiệu quả  
**Để** đạt hiệu năng real-time

**Acceptance Criteria:**
- [ ] Digital Down Converter (DDC) với decimation
- [ ] Digital Up Converter (DUC) với interpolation
- [ ] Configurable filter coefficients
- [ ] Resource utilization < 50% DSP slices
- [ ] Processing latency < 10μs

#### US6.2: Digital Filters
**Là một** system  
**Tôi muốn** lọc tín hiệu với độ chính xác cao  
**Để** cải thiện chất lượng audio

**Acceptance Criteria:**
- [ ] FIR filter bank với multiple taps
- [ ] IIR filters cho specific applications
- [ ] Adaptive filtering capabilities
- [ ] Coefficient updates without interruption
- [ ] Filter response verification

#### US6.3: Modulation/Demodulation
**Là một** system  
**Tôi muốn** xử lý modulation trong hardware  
**Để** giảm tải cho CPU

**Acceptance Criteria:**
- [ ] AM modulation/demodulation
- [ ] FM modulation/demodulation
- [ ] SSB generation và detection
- [ ] Carrier recovery algorithms
- [ ] Symbol timing recovery

---

## EPIC 7: System Integration và Testing

### Epic Description
Tích hợp tất cả các component và thực hiện testing toàn diện để đảm bảo chất lượng hệ thống.

### User Stories

#### US7.1: Hardware-Software Integration
**Là một** developer  
**Tôi muốn** tích hợp software với hardware  
**Để** tạo ra hệ thống hoàn chỉnh

**Acceptance Criteria:**
- [ ] AD9364 driver integration
- [ ] FPGA-PS communication qua AXI
- [ ] GPIO control cho external devices
- [ ] Interrupt handling optimization
- [ ] DMA transfers cho high-speed data

#### US7.2: Performance Testing
**Là một** QA engineer  
**Tôi muốn** verify hiệu năng hệ thống  
**Để** đảm bảo đáp ứng requirements

**Acceptance Criteria:**
- [ ] RF performance measurements
- [ ] Audio quality testing (THD, SNR)
- [ ] Latency measurements end-to-end
- [ ] Stress testing với maximum load
- [ ] Long-term stability testing

#### US7.3: Compliance Testing
**Là một** certification engineer  
**Tôi muốn** verify tuân thủ standards  
**Để** đạt certification

**Acceptance Criteria:**
- [ ] ED-137 compliance testing
- [ ] EMC/EMI testing
- [ ] Safety standards verification
- [ ] Environmental testing
- [ ] Documentation compliance

---

## EPIC 8: Documentation và Training

### Epic Description
Tạo documentation đầy đủ và materials training cho users và maintainers.

### User Stories

#### US8.1: Technical Documentation
**Là một** maintainer  
**Tôi muốn** có documentation chi tiết  
**Để** maintain và troubleshoot hệ thống

**Acceptance Criteria:**
- [ ] System architecture documentation
- [ ] API documentation với examples
- [ ] Configuration guides
- [ ] Troubleshooting procedures
- [ ] Maintenance schedules

#### US8.2: User Manuals
**Là một** end user  
**Tôi muốn** có hướng dẫn sử dụng  
**Để** vận hành hệ thống hiệu quả

**Acceptance Criteria:**
- [ ] Quick start guide
- [ ] Detailed operation manual
- [ ] Safety procedures
- [ ] Emergency procedures
- [ ] Multi-language support

#### US8.3: Training Materials
**Là một** trainer  
**Tôi muốn** có materials training  
**Để** đào tạo operators

**Acceptance Criteria:**
- [ ] Training presentations
- [ ] Hands-on exercises
- [ ] Video tutorials
- [ ] Assessment materials
- [ ] Certification programs

---

## Priority Matrix

### High Priority (Must Have)
- US1.1: RF Transceiver Control
- US1.3: PTT Control
- US2.1: SIP Protocol Implementation
- US3.1: System Dashboard
- US6.1: DDC/DUC Implementation

### Medium Priority (Should Have)
- US1.2: Channel Management System
- US2.2: RTP Media Streaming
- US3.2: Configuration Management
- US4.1: SNMP v3 Implementation
- US5.1: Configuration Database

### Low Priority (Could Have)
- US2.3: Multi-frequency Operation
- US3.4: Remote Control Interface
- US4.3: Trap Generation
- US6.3: Modulation/Demodulation
- US8.3: Training Materials

---

## Release Planning

### Release 1.0 (MVP) - 6 months
- Core radio functionality
- Basic web interface
- VoIP integration
- Essential FPGA processing

### Release 1.1 - 8 months
- Advanced web features
- SNMP management
- Database services
- Performance optimization

### Release 2.0 - 12 months
- Multi-channel operation
- Advanced signal processing
- Compliance certification
- Complete documentation

---

## Technical Constraints

### Hardware Constraints
- XCZU-2EG MPSoC với 240 DSP slices
- 2GB DDR4 memory
- Single AD9364 transceiver
- Limited I/O pins

### Software Constraints
- Linux RT kernel requirement
- Real-time performance < 20ms latency
- ED-137 compliance mandatory
- Memory footprint < 1GB

### Environmental Constraints
- Operating temperature: -40°C to +70°C
- Humidity: 5% to 95% non-condensing
- Vibration resistance per MIL-STD-810
- EMC compliance per CISPR 25

---

## Definition of Done

Một user story được coi là "Done" khi:
- [ ] Code được implement và reviewed
- [ ] Unit tests pass với coverage > 80%
- [ ] Integration tests pass
- [ ] Performance requirements được verify
- [ ] Documentation được update
- [ ] Security review completed
- [ ] Acceptance criteria được fulfill
- [ ] Stakeholder sign-off
