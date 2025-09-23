import { useState, useEffect } from 'react'
import { 
  HardDrive, 
  Activity, 
  AlertTriangle, 
  CheckCircle,
  Database,
  Server
} from 'lucide-react'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts'

const storageData = [
  { time: '00:00', usage: 45 },
  { time: '04:00', usage: 48 },
  { time: '08:00', usage: 52 },
  { time: '12:00', usage: 55 },
  { time: '16:00', usage: 58 },
  { time: '20:00', usage: 60 },
  { time: '24:00', usage: 62 },
]

const storageBreakdown = [
  { name: 'Local Storage', value: 40, color: '#3B82F6' },
  { name: 'NAS Storage', value: 35, color: '#10B981' },
  { name: 'SAN Storage', value: 20, color: '#F59E0B' },
  { name: 'Object Storage', value: 5, color: '#EF4444' },
]

export function Dashboard() {
  const [metrics, setMetrics] = useState({
    totalStorage: '2.5 TB',
    usedStorage: '1.5 TB',
    freeStorage: '1.0 TB',
    activeBackends: 4,
    healthyBackends: 3,
    alerts: 2
  })

  const [health] = useState({
    overall: 'healthy',
    backends: [
      { name: 'Local RAID-6', status: 'healthy', usage: 65 },
      { name: 'NAS Primary', status: 'healthy', usage: 45 },
      { name: 'SAN iSCSI', status: 'warning', usage: 85 },
      { name: 'MinIO Object', status: 'healthy', usage: 25 },
    ]
  })

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('/api/v1/storage/metrics')
        if (response.ok) {
          const data = await response.json()
          setMetrics(data)
        }
      } catch (error) {
        console.error('Failed to fetch metrics:', error)
      }
    }

    fetchData()
    const interval = setInterval(fetchData, 30000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Storage Dashboard</h1>
        <div className="flex items-center space-x-2">
          <div className="flex items-center text-sm text-gray-500">
            <Activity className="w-4 h-4 mr-1" />
            Last updated: {new Date().toLocaleTimeString()}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <HardDrive className="h-6 w-6 text-gray-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Total Storage
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {metrics.totalStorage}
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
                <Database className="h-6 w-6 text-gray-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Used Storage
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {metrics.usedStorage}
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
                <Server className="h-6 w-6 text-gray-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Active Backends
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {metrics.activeBackends}
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
                <AlertTriangle className="h-6 w-6 text-yellow-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Active Alerts
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {metrics.alerts}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Storage Usage Trend</h3>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={storageData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis />
              <Tooltip />
              <Line type="monotone" dataKey="usage" stroke="#3B82F6" strokeWidth={2} />
            </LineChart>
          </ResponsiveContainer>
        </div>

        <div className="bg-white shadow rounded-lg p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Storage Distribution</h3>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={storageBreakdown}
                cx="50%"
                cy="50%"
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
              >
                {storageBreakdown.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="bg-white shadow rounded-lg">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Backend Health Status</h3>
        </div>
        <div className="px-6 py-4">
          <div className="space-y-4">
            {health.backends.map((backend, index) => (
              <div key={index} className="flex items-center justify-between">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    {backend.status === 'healthy' ? (
                      <CheckCircle className="h-5 w-5 text-green-500" />
                    ) : (
                      <AlertTriangle className="h-5 w-5 text-yellow-500" />
                    )}
                  </div>
                  <div className="ml-3">
                    <p className="text-sm font-medium text-gray-900">{backend.name}</p>
                    <p className="text-sm text-gray-500">
                      Status: {backend.status} • Usage: {backend.usage}%
                    </p>
                  </div>
                </div>
                <div className="flex items-center">
                  <div className="w-32 bg-gray-200 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full ${
                        backend.usage > 80 ? 'bg-red-500' : 
                        backend.usage > 60 ? 'bg-yellow-500' : 'bg-green-500'
                      }`}
                      style={{ width: `${backend.usage}%` }}
                    />
                  </div>
                  <span className="ml-2 text-sm text-gray-500">{backend.usage}%</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
