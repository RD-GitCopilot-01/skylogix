import { useState, useEffect } from 'react'
import { 
  Play, 
  Square, 
  Settings, 
  Network,
  AlertTriangle,
  CheckCircle,
  Activity,
  Send
} from 'lucide-react'

interface SNMPStatus {
  running: boolean
  version: string
  port: number
  community: string
  uptime: number
  requests_received: number
  traps_sent: number
}

interface MIBObject {
  oid: string
  name: string
  type: string
  value: string
  description: string
}

export function SNMP() {
  const [status, setStatus] = useState<SNMPStatus>({
    running: true,
    version: 'v2c',
    port: 161,
    community: 'public',
    uptime: 86400,
    requests_received: 1247,
    traps_sent: 23
  })

  const [mibObjects] = useState<MIBObject[]>([
    {
      oid: '1.3.6.1.4.1.12345.1.1.1',
      name: 'skylogixStorageTotal',
      type: 'Counter64',
      value: '2684354560000',
      description: 'Total storage capacity in bytes'
    },
    {
      oid: '1.3.6.1.4.1.12345.1.1.2',
      name: 'skylogixStorageUsed',
      type: 'Counter64',
      value: '1610612736000',
      description: 'Used storage capacity in bytes'
    },
    {
      oid: '1.3.6.1.4.1.12345.1.1.3',
      name: 'skylogixStorageFree',
      type: 'Counter64',
      value: '1073741824000',
      description: 'Free storage capacity in bytes'
    },
    {
      oid: '1.3.6.1.4.1.12345.1.2.1',
      name: 'skylogixBackendCount',
      type: 'Integer32',
      value: '4',
      description: 'Number of configured storage backends'
    },
    {
      oid: '1.3.6.1.4.1.12345.1.2.2',
      name: 'skylogixHealthyBackends',
      type: 'Integer32',
      value: '3',
      description: 'Number of healthy storage backends'
    },
    {
      oid: '1.3.6.1.4.1.12345.1.3.1',
      name: 'skylogixRAIDStatus',
      type: 'Integer32',
      value: '1',
      description: 'RAID array status (1=healthy, 2=degraded, 3=failed)'
    },
    {
      oid: '1.3.6.1.4.1.12345.1.3.2',
      name: 'skylogixRAIDLevel',
      type: 'OctetString',
      value: 'RAID-6',
      description: 'Current RAID level configuration'
    },
    {
      oid: '1.3.6.1.4.1.12345.1.4.1',
      name: 'skylogixIOPS',
      type: 'Counter32',
      value: '1650',
      description: 'Current I/O operations per second'
    }
  ])

  const [config, setConfig] = useState({
    enabled: true,
    port: 161,
    community: 'public',
    version: 'v2c',
    trap_destination: '192.168.1.100',
    trap_port: 162
  })

  const [testTrap, setTestTrap] = useState({
    type: 'warning',
    message: 'Test SNMP trap from Skylogix Storage'
  })

  const startSNMP = async () => {
    try {
      const response = await fetch('/api/v1/snmp/start', {
        method: 'POST',
      })
      
      if (response.ok) {
        setStatus(prev => ({ ...prev, running: true }))
      }
    } catch (error) {
      console.error('Failed to start SNMP agent:', error)
    }
  }

  const stopSNMP = async () => {
    try {
      const response = await fetch('/api/v1/snmp/stop', {
        method: 'POST',
      })
      
      if (response.ok) {
        setStatus(prev => ({ ...prev, running: false }))
      }
    } catch (error) {
      console.error('Failed to stop SNMP agent:', error)
    }
  }

  const sendTestTrap = async () => {
    try {
      const response = await fetch('/api/v1/snmp/trap', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          type: testTrap.type,
          message: testTrap.message,
          destination: config.trap_destination,
          port: config.trap_port
        }),
      })
      
      if (response.ok) {
        setStatus(prev => ({ ...prev, traps_sent: prev.traps_sent + 1 }))
      }
    } catch (error) {
      console.error('Failed to send test trap:', error)
    }
  }

  const formatUptime = (seconds: number) => {
    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    return `${days}d ${hours}h ${minutes}m`
  }

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const response = await fetch('/api/v1/snmp/status')
        if (response.ok) {
          const data = await response.json()
          setStatus(data)
        }
      } catch (error) {
        console.error('Failed to fetch SNMP status:', error)
      }
    }

    fetchStatus()
    const interval = setInterval(fetchStatus, 30000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">SNMP Management</h1>
        <div className="flex items-center space-x-3">
          <div className={`flex items-center text-sm ${status.running ? 'text-green-600' : 'text-red-600'}`}>
            {status.running ? <CheckCircle className="w-4 h-4 mr-1" /> : <AlertTriangle className="w-4 h-4 mr-1" />}
            {status.running ? 'Running' : 'Stopped'}
          </div>
          {status.running ? (
            <button
              onClick={stopSNMP}
              className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
            >
              <Square className="w-4 h-4 mr-2" />
              Stop Agent
            </button>
          ) : (
            <button
              onClick={startSNMP}
              className="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700"
            >
              <Play className="w-4 h-4 mr-2" />
              Start Agent
            </button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <Network className="h-6 w-6 text-blue-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    SNMP Version
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {status.version.toUpperCase()}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <Activity className="h-6 w-6 text-green-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Uptime
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {formatUptime(status.uptime)}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <CheckCircle className="h-6 w-6 text-purple-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Requests
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {status.requests_received.toLocaleString()}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <Send className="h-6 w-6 text-orange-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Traps Sent
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {status.traps_sent}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white shadow rounded-lg">
          <div className="px-6 py-4 border-b border-gray-200">
            <h3 className="text-lg font-medium text-gray-900 flex items-center">
              <Settings className="w-5 h-5 mr-2 text-gray-400" />
              SNMP Configuration
            </h3>
          </div>
          <div className="px-6 py-4">
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Port</label>
                  <input
                    type="number"
                    value={config.port}
                    onChange={(e) => setConfig(prev => ({ ...prev, port: parseInt(e.target.value) }))}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Version</label>
                  <select
                    value={config.version}
                    onChange={(e) => setConfig(prev => ({ ...prev, version: e.target.value }))}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="v1">SNMPv1</option>
                    <option value="v2c">SNMPv2c</option>
                    <option value="v3">SNMPv3</option>
                  </select>
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Community String</label>
                <input
                  type="text"
                  value={config.community}
                  onChange={(e) => setConfig(prev => ({ ...prev, community: e.target.value }))}
                  className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Trap Destination</label>
                  <input
                    type="text"
                    value={config.trap_destination}
                    onChange={(e) => setConfig(prev => ({ ...prev, trap_destination: e.target.value }))}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Trap Port</label>
                  <input
                    type="number"
                    value={config.trap_port}
                    onChange={(e) => setConfig(prev => ({ ...prev, trap_port: parseInt(e.target.value) }))}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>
              </div>
              <div className="pt-4">
                <h4 className="text-sm font-medium text-gray-700 mb-2">Test Trap</h4>
                <div className="flex space-x-2">
                  <select
                    value={testTrap.type}
                    onChange={(e) => setTestTrap(prev => ({ ...prev, type: e.target.value }))}
                    className="block border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="info">Info</option>
                    <option value="warning">Warning</option>
                    <option value="critical">Critical</option>
                  </select>
                  <input
                    type="text"
                    value={testTrap.message}
                    onChange={(e) => setTestTrap(prev => ({ ...prev, message: e.target.value }))}
                    className="flex-1 block border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Trap message"
                  />
                  <button
                    onClick={sendTestTrap}
                    className="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
                  >
                    <Send className="w-4 h-4 mr-1" />
                    Send
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white shadow rounded-lg">
          <div className="px-6 py-4 border-b border-gray-200">
            <h3 className="text-lg font-medium text-gray-900">MIB Objects</h3>
          </div>
          <div className="px-6 py-4">
            <div className="space-y-3 max-h-96 overflow-y-auto">
              {mibObjects.map((obj, index) => (
                <div key={index} className="border border-gray-200 rounded-lg p-3">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-900">{obj.name}</span>
                    <span className="text-xs text-gray-500">{obj.type}</span>
                  </div>
                  <div className="text-xs text-gray-600 mb-1">OID: {obj.oid}</div>
                  <div className="text-sm text-gray-800 mb-1">Value: {obj.value}</div>
                  <div className="text-xs text-gray-500">{obj.description}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
