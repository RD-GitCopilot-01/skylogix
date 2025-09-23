import React, { useState, useEffect } from 'react'
import { 
  Database as DatabaseIcon, 
  Download, 
  Upload, 
  Trash2, 
  RefreshCw,
  HardDrive,
  Activity,
  Clock,
  BarChart3
} from 'lucide-react'

interface DatabaseStats {
  total_size: string
  table_count: number
  record_count: number
  last_backup: Date | null
  uptime: number
  connections: number
}

interface TableInfo {
  name: string
  records: number
  size: string
  last_updated: Date
}

export function Database() {
  const [stats, setStats] = useState<DatabaseStats>({
    total_size: '245 MB',
    table_count: 12,
    record_count: 156789,
    last_backup: new Date(Date.now() - 6 * 60 * 60 * 1000),
    uptime: 432000,
    connections: 3
  })

  const [tables] = useState<TableInfo[]>([
    {
      name: 'storage_metrics',
      records: 89456,
      size: '125 MB',
      last_updated: new Date(Date.now() - 5 * 60 * 1000)
    },
    {
      name: 'storage_configs',
      records: 24,
      size: '2.1 MB',
      last_updated: new Date(Date.now() - 2 * 60 * 60 * 1000)
    },
    {
      name: 'audit_logs',
      records: 45678,
      size: '89 MB',
      last_updated: new Date(Date.now() - 15 * 60 * 1000)
    },
    {
      name: 'alert_history',
      records: 12345,
      size: '15 MB',
      last_updated: new Date(Date.now() - 30 * 60 * 1000)
    },
    {
      name: 'snmp_requests',
      records: 8901,
      size: '8.5 MB',
      last_updated: new Date(Date.now() - 10 * 60 * 1000)
    },
    {
      name: 'backend_health',
      records: 2345,
      size: '3.2 MB',
      last_updated: new Date(Date.now() - 1 * 60 * 1000)
    }
  ])

  const [operations, setOperations] = useState({
    backup: { running: false, progress: 0 },
    restore: { running: false, progress: 0 },
    cleanup: { running: false, progress: 0 }
  })

  const [retentionSettings, setRetentionSettings] = useState({
    metrics_days: 30,
    audit_days: 90,
    alerts_days: 60,
    auto_cleanup: true
  })

  const formatUptime = (seconds: number) => {
    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    return `${days}d ${hours}h ${minutes}m`
  }

  const startBackup = async () => {
    setOperations(prev => ({ ...prev, backup: { running: true, progress: 0 } }))
    
    try {
      const response = await fetch('/api/v1/database/backup', {
        method: 'POST',
      })
      
      if (response.ok) {
        for (let i = 0; i <= 100; i += 10) {
          await new Promise(resolve => setTimeout(resolve, 200))
          setOperations(prev => ({ ...prev, backup: { running: true, progress: i } }))
        }
        
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `skylogix-db-backup-${new Date().toISOString().split('T')[0]}.sql`
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
        
        setStats(prev => ({ ...prev, last_backup: new Date() }))
      }
    } catch (error) {
      console.error('Failed to backup database:', error)
    } finally {
      setOperations(prev => ({ ...prev, backup: { running: false, progress: 0 } }))
    }
  }

  const startRestore = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    setOperations(prev => ({ ...prev, restore: { running: true, progress: 0 } }))

    const formData = new FormData()
    formData.append('backup', file)

    fetch('/api/v1/database/restore', {
      method: 'POST',
      body: formData,
    })
    .then(async response => {
      if (response.ok) {
        for (let i = 0; i <= 100; i += 10) {
          await new Promise(resolve => setTimeout(resolve, 300))
          setOperations(prev => ({ ...prev, restore: { running: true, progress: i } }))
        }
      }
    })
    .catch(error => {
      console.error('Failed to restore database:', error)
    })
    .finally(() => {
      setOperations(prev => ({ ...prev, restore: { running: false, progress: 0 } }))
    })
  }

  const startCleanup = async () => {
    setOperations(prev => ({ ...prev, cleanup: { running: true, progress: 0 } }))
    
    try {
      const response = await fetch('/api/v1/database/cleanup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(retentionSettings),
      })
      
      if (response.ok) {
        for (let i = 0; i <= 100; i += 20) {
          await new Promise(resolve => setTimeout(resolve, 500))
          setOperations(prev => ({ ...prev, cleanup: { running: true, progress: i } }))
        }
      }
    } catch (error) {
      console.error('Failed to cleanup database:', error)
    } finally {
      setOperations(prev => ({ ...prev, cleanup: { running: false, progress: 0 } }))
    }
  }

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await fetch('/api/v1/database/stats')
        if (response.ok) {
          const data = await response.json()
          setStats(data)
        }
      } catch (error) {
        console.error('Failed to fetch database stats:', error)
      }
    }

    fetchStats()
    const interval = setInterval(fetchStats, 30000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Database Management</h1>
        <div className="flex items-center space-x-3">
          <div className="flex items-center text-sm text-gray-500">
            <Activity className="w-4 h-4 mr-1" />
            {stats.connections} active connections
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="p-5">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <HardDrive className="h-6 w-6 text-blue-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Database Size
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {stats.total_size}
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
                <BarChart3 className="h-6 w-6 text-green-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Total Records
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {stats.record_count.toLocaleString()}
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
                <DatabaseIcon className="h-6 w-6 text-purple-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Tables
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {stats.table_count}
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
                <Clock className="h-6 w-6 text-orange-400" />
              </div>
              <div className="ml-5 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    Uptime
                  </dt>
                  <dd className="text-lg font-medium text-gray-900">
                    {formatUptime(stats.uptime)}
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
            <h3 className="text-lg font-medium text-gray-900">Database Operations</h3>
          </div>
          <div className="px-6 py-4">
            <div className="space-y-4">
              <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Backup Database</h4>
                  <p className="text-sm text-gray-500">
                    {stats.last_backup 
                      ? `Last backup: ${stats.last_backup.toLocaleString()}`
                      : 'No backup available'
                    }
                  </p>
                  {operations.backup.running && (
                    <div className="mt-2">
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                          style={{ width: `${operations.backup.progress}%` }}
                        />
                      </div>
                      <p className="text-xs text-gray-500 mt-1">{operations.backup.progress}% complete</p>
                    </div>
                  )}
                </div>
                <button
                  onClick={startBackup}
                  disabled={operations.backup.running}
                  className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
                >
                  {operations.backup.running ? (
                    <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                  ) : (
                    <Download className="w-4 h-4 mr-2" />
                  )}
                  {operations.backup.running ? 'Backing up...' : 'Backup'}
                </button>
              </div>

              <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Restore Database</h4>
                  <p className="text-sm text-gray-500">Restore from backup file</p>
                  {operations.restore.running && (
                    <div className="mt-2">
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-green-600 h-2 rounded-full transition-all duration-300"
                          style={{ width: `${operations.restore.progress}%` }}
                        />
                      </div>
                      <p className="text-xs text-gray-500 mt-1">{operations.restore.progress}% complete</p>
                    </div>
                  )}
                </div>
                <div>
                  <input
                    type="file"
                    accept=".sql,.db"
                    onChange={startRestore}
                    disabled={operations.restore.running}
                    className="hidden"
                    id="restore-file"
                  />
                  <label
                    htmlFor="restore-file"
                    className={`inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 cursor-pointer ${
                      operations.restore.running ? 'opacity-50 cursor-not-allowed' : ''
                    }`}
                  >
                    {operations.restore.running ? (
                      <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                    ) : (
                      <Upload className="w-4 h-4 mr-2" />
                    )}
                    {operations.restore.running ? 'Restoring...' : 'Restore'}
                  </label>
                </div>
              </div>

              <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Cleanup Old Data</h4>
                  <p className="text-sm text-gray-500">Remove old records based on retention policy</p>
                  {operations.cleanup.running && (
                    <div className="mt-2">
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-red-600 h-2 rounded-full transition-all duration-300"
                          style={{ width: `${operations.cleanup.progress}%` }}
                        />
                      </div>
                      <p className="text-xs text-gray-500 mt-1">{operations.cleanup.progress}% complete</p>
                    </div>
                  )}
                </div>
                <button
                  onClick={startCleanup}
                  disabled={operations.cleanup.running}
                  className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
                >
                  {operations.cleanup.running ? (
                    <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                  ) : (
                    <Trash2 className="w-4 h-4 mr-2" />
                  )}
                  {operations.cleanup.running ? 'Cleaning...' : 'Cleanup'}
                </button>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white shadow rounded-lg">
          <div className="px-6 py-4 border-b border-gray-200">
            <h3 className="text-lg font-medium text-gray-900">Table Information</h3>
          </div>
          <div className="px-6 py-4">
            <div className="space-y-3 max-h-96 overflow-y-auto">
              {tables.map((table, index) => (
                <div key={index} className="border border-gray-200 rounded-lg p-3">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-900">{table.name}</span>
                    <span className="text-xs text-gray-500">{table.size}</span>
                  </div>
                  <div className="flex items-center justify-between text-sm text-gray-600">
                    <span>{table.records.toLocaleString()} records</span>
                    <span>Updated {table.last_updated.toLocaleTimeString()}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="bg-white shadow rounded-lg">
        <div className="px-6 py-4 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Retention Settings</h3>
        </div>
        <div className="px-6 py-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Metrics Retention (days)</label>
              <input
                type="number"
                value={retentionSettings.metrics_days}
                onChange={(e) => setRetentionSettings(prev => ({ ...prev, metrics_days: parseInt(e.target.value) }))}
                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Audit Logs Retention (days)</label>
              <input
                type="number"
                value={retentionSettings.audit_days}
                onChange={(e) => setRetentionSettings(prev => ({ ...prev, audit_days: parseInt(e.target.value) }))}
                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Alerts Retention (days)</label>
              <input
                type="number"
                value={retentionSettings.alerts_days}
                onChange={(e) => setRetentionSettings(prev => ({ ...prev, alerts_days: parseInt(e.target.value) }))}
                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </div>
          <div className="mt-4">
            <div className="flex items-center">
              <input
                type="checkbox"
                checked={retentionSettings.auto_cleanup}
                onChange={(e) => setRetentionSettings(prev => ({ ...prev, auto_cleanup: e.target.checked }))}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <label className="ml-2 block text-sm text-gray-900">
                Enable automatic cleanup based on retention settings
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
