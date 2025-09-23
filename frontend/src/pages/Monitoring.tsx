import { useState, useEffect } from 'react'
import { 
  Activity, 
  AlertTriangle, 
  CheckCircle, 
  Clock,
  TrendingUp,
  TrendingDown,
  Eye,
  EyeOff
} from 'lucide-react'
import { XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from 'recharts'

interface Alert {
  id: string
  type: 'warning' | 'critical' | 'info'
  title: string
  message: string
  timestamp: Date
  acknowledged: boolean
  source: string
}

interface MetricData {
  timestamp: string
  storage_usage: number
  iops: number
  throughput: number
  latency: number
}

export function Monitoring() {
  const [alerts, setAlerts] = useState<Alert[]>([
    {
      id: '1',
      type: 'warning',
      title: 'High Storage Usage',
      message: 'SAN iSCSI storage usage is at 85%',
      timestamp: new Date(Date.now() - 5 * 60 * 1000),
      acknowledged: false,
      source: 'SAN iSCSI'
    },
    {
      id: '2',
      type: 'critical',
      title: 'RAID Disk Failure',
      message: 'Disk /dev/sdb has failed in RAID-6 array',
      timestamp: new Date(Date.now() - 15 * 60 * 1000),
      acknowledged: false,
      source: 'Local RAID-6'
    },
    {
      id: '3',
      type: 'info',
      title: 'Backup Completed',
      message: 'Daily backup to object storage completed successfully',
      timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000),
      acknowledged: true,
      source: 'MinIO Object'
    }
  ])

  const [metrics, setMetrics] = useState<MetricData[]>([
    { timestamp: '00:00', storage_usage: 45, iops: 1200, throughput: 85, latency: 2.1 },
    { timestamp: '04:00', storage_usage: 48, iops: 1350, throughput: 92, latency: 2.3 },
    { timestamp: '08:00', storage_usage: 52, iops: 1800, throughput: 125, latency: 1.8 },
    { timestamp: '12:00', storage_usage: 55, iops: 2100, throughput: 145, latency: 1.9 },
    { timestamp: '16:00', storage_usage: 58, iops: 1950, throughput: 135, latency: 2.0 },
    { timestamp: '20:00', storage_usage: 60, iops: 1600, throughput: 110, latency: 2.2 },
    { timestamp: '24:00', storage_usage: 62, iops: 1400, throughput: 95, latency: 2.4 }
  ])

  const [selectedMetric, setSelectedMetric] = useState<'storage_usage' | 'iops' | 'throughput' | 'latency'>('storage_usage')
  const [showAcknowledged, setShowAcknowledged] = useState(false)

  const acknowledgeAlert = async (alertId: string) => {
    try {
      const response = await fetch(`/api/v1/monitoring/alerts/${alertId}/acknowledge`, {
        method: 'POST',
      })
      
      if (response.ok) {
        setAlerts(prev => prev.map(alert => 
          alert.id === alertId ? { ...alert, acknowledged: true } : alert
        ))
      }
    } catch (error) {
      console.error('Failed to acknowledge alert:', error)
    }
  }

  const getAlertIcon = (type: string) => {
    switch (type) {
      case 'critical': return AlertTriangle
      case 'warning': return AlertTriangle
      case 'info': return CheckCircle
      default: return Activity
    }
  }

  const getAlertColor = (type: string) => {
    switch (type) {
      case 'critical': return 'text-red-500 bg-red-50 border-red-200'
      case 'warning': return 'text-yellow-500 bg-yellow-50 border-yellow-200'
      case 'info': return 'text-blue-500 bg-blue-50 border-blue-200'
      default: return 'text-gray-500 bg-gray-50 border-gray-200'
    }
  }

  const getMetricColor = (metric: string) => {
    switch (metric) {
      case 'storage_usage': return '#3B82F6'
      case 'iops': return '#10B981'
      case 'throughput': return '#F59E0B'
      case 'latency': return '#EF4444'
      default: return '#6B7280'
    }
  }

  const getMetricLabel = (metric: string) => {
    switch (metric) {
      case 'storage_usage': return 'Storage Usage (%)'
      case 'iops': return 'IOPS'
      case 'throughput': return 'Throughput (MB/s)'
      case 'latency': return 'Latency (ms)'
      default: return metric
    }
  }

  const filteredAlerts = showAcknowledged ? alerts : alerts.filter(alert => !alert.acknowledged)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [alertsResponse, metricsResponse] = await Promise.all([
          fetch('/api/v1/monitoring/alerts'),
          fetch('/api/v1/monitoring/metrics/storage')
        ])
        
        if (alertsResponse.ok) {
          const alertsData = await alertsResponse.json()
          setAlerts(alertsData)
        }
        
        if (metricsResponse.ok) {
          const metricsData = await metricsResponse.json()
          setMetrics(metricsData)
        }
      } catch (error) {
        console.error('Failed to fetch monitoring data:', error)
      }
    }

    fetchData()
    const interval = setInterval(fetchData, 30000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Monitoring & Alerts</h1>
        <div className="flex items-center space-x-2">
          <div className="flex items-center text-sm text-gray-500">
            <Activity className="w-4 h-4 mr-1" />
            Real-time monitoring
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          <div className="bg-white shadow rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-medium text-gray-900">Performance Metrics</h3>
              <div className="flex space-x-2">
                {(['storage_usage', 'iops', 'throughput', 'latency'] as const).map((metric) => (
                  <button
                    key={metric}
                    onClick={() => setSelectedMetric(metric)}
                    className={`px-3 py-1 text-xs font-medium rounded-full ${
                      selectedMetric === metric
                        ? 'bg-blue-100 text-blue-800'
                        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                    }`}
                  >
                    {getMetricLabel(metric)}
                  </button>
                ))}
              </div>
            </div>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={metrics}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="timestamp" />
                <YAxis />
                <Tooltip />
                <Area
                  type="monotone"
                  dataKey={selectedMetric}
                  stroke={getMetricColor(selectedMetric)}
                  fill={getMetricColor(selectedMetric)}
                  fillOpacity={0.3}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white shadow rounded-lg">
          <div className="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">Active Alerts</h3>
            <button
              onClick={() => setShowAcknowledged(!showAcknowledged)}
              className="flex items-center text-sm text-gray-500 hover:text-gray-700"
            >
              {showAcknowledged ? <EyeOff className="w-4 h-4 mr-1" /> : <Eye className="w-4 h-4 mr-1" />}
              {showAcknowledged ? 'Hide' : 'Show'} acknowledged
            </button>
          </div>
          <div className="px-6 py-4">
            <div className="space-y-4">
              {filteredAlerts.length === 0 ? (
                <div className="text-center py-8">
                  <CheckCircle className="w-12 h-12 mx-auto text-green-500 mb-3" />
                  <p className="text-gray-500">No active alerts</p>
                </div>
              ) : (
                filteredAlerts.map((alert) => {
                  const Icon = getAlertIcon(alert.type)
                  return (
                    <div
                      key={alert.id}
                      className={`p-4 rounded-lg border ${getAlertColor(alert.type)} ${
                        alert.acknowledged ? 'opacity-60' : ''
                      }`}
                    >
                      <div className="flex items-start">
                        <Icon className="w-5 h-5 mt-0.5 mr-3" />
                        <div className="flex-1 min-w-0">
                          <h4 className="text-sm font-medium text-gray-900">{alert.title}</h4>
                          <p className="text-sm text-gray-600 mt-1">{alert.message}</p>
                          <div className="flex items-center justify-between mt-2">
                            <div className="flex items-center text-xs text-gray-500">
                              <Clock className="w-3 h-3 mr-1" />
                              {alert.timestamp.toLocaleTimeString()}
                            </div>
                            <div className="text-xs text-gray-500">
                              {alert.source}
                            </div>
                          </div>
                          {!alert.acknowledged && (
                            <button
                              onClick={() => acknowledgeAlert(alert.id)}
                              className="mt-2 text-xs text-blue-600 hover:text-blue-800"
                            >
                              Acknowledge
                            </button>
                          )}
                        </div>
                      </div>
                    </div>
                  )
                })
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <TrendingUp className="h-6 w-6 text-green-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Avg IOPS
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    1,650
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
                <TrendingUp className="h-6 w-6 text-blue-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Throughput
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    115 MB/s
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
                <TrendingDown className="h-6 w-6 text-yellow-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Avg Latency
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    2.1 ms
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
                <AlertTriangle className="h-6 w-6 text-red-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Active Alerts
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {alerts.filter(a => !a.acknowledged).length}
                  </dd>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
