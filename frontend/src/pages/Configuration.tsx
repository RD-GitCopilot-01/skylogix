import React, { useState } from 'react'
import { 
  Save, 
  Upload, 
  Download, 
  Settings,
  CheckCircle
} from 'lucide-react'

interface ConfigSection {
  name: string
  settings: ConfigSetting[]
}

interface ConfigSetting {
  key: string
  label: string
  value: string | number | boolean
  type: 'text' | 'number' | 'boolean' | 'select'
  options?: string[]
  description?: string
}

export function Configuration() {
  const [config, setConfig] = useState<ConfigSection[]>([
    {
      name: 'Storage Settings',
      settings: [
        {
          key: 'default_backend',
          label: 'Default Storage Backend',
          value: 'local',
          type: 'select',
          options: ['local', 'nas', 'san', 'object'],
          description: 'Primary storage backend for new data'
        },
        {
          key: 'encryption_enabled',
          label: 'Enable Encryption',
          value: true,
          type: 'boolean',
          description: 'Enable AES-256 encryption for all storage backends'
        },
        {
          key: 'retention_days',
          label: 'Default Retention (Days)',
          value: 30,
          type: 'number',
          description: 'Default data retention period in days'
        }
      ]
    },
    {
      name: 'RAID Configuration',
      settings: [
        {
          key: 'raid_level',
          label: 'RAID Level',
          value: '6',
          type: 'select',
          options: ['0', '1', '5', '6', '10'],
          description: 'RAID level for local storage'
        },
        {
          key: 'chunk_size',
          label: 'Chunk Size (KB)',
          value: 64,
          type: 'number',
          description: 'RAID chunk size in kilobytes'
        },
        {
          key: 'auto_rebuild',
          label: 'Auto Rebuild',
          value: true,
          type: 'boolean',
          description: 'Automatically rebuild RAID on disk failure'
        }
      ]
    },
    {
      name: 'Monitoring & Alerts',
      settings: [
        {
          key: 'snmp_enabled',
          label: 'Enable SNMP',
          value: true,
          type: 'boolean',
          description: 'Enable SNMP monitoring agent'
        },
        {
          key: 'alert_threshold',
          label: 'Storage Alert Threshold (%)',
          value: 85,
          type: 'number',
          description: 'Storage usage percentage to trigger alerts'
        },
        {
          key: 'health_check_interval',
          label: 'Health Check Interval (seconds)',
          value: 60,
          type: 'number',
          description: 'Interval between storage health checks'
        }
      ]
    }
  ])

  const [saving, setSaving] = useState(false)
  const [lastSaved, setLastSaved] = useState<Date | null>(null)

  const updateSetting = (sectionIndex: number, settingIndex: number, value: any) => {
    setConfig(prev => {
      const newConfig = [...prev]
      newConfig[sectionIndex].settings[settingIndex].value = value
      return newConfig
    })
  }

  const saveConfiguration = async () => {
    setSaving(true)
    try {
      const response = await fetch('/api/v1/config/', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(config),
      })
      
      if (response.ok) {
        setLastSaved(new Date())
      }
    } catch (error) {
      console.error('Failed to save configuration:', error)
    } finally {
      setSaving(false)
    }
  }

  const backupConfiguration = async () => {
    try {
      const response = await fetch('/api/v1/config/backup', {
        method: 'POST',
      })
      
      if (response.ok) {
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `skylogix-config-${new Date().toISOString().split('T')[0]}.json`
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
      }
    } catch (error) {
      console.error('Failed to backup configuration:', error)
    }
  }

  const restoreConfiguration = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    const reader = new FileReader()
    reader.onload = async (e) => {
      try {
        const configData = JSON.parse(e.target?.result as string)
        
        const response = await fetch('/api/v1/config/restore', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(configData),
        })
        
        if (response.ok) {
          setConfig(configData)
          setLastSaved(new Date())
        }
      } catch (error) {
        console.error('Failed to restore configuration:', error)
      }
    }
    reader.readAsText(file)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Configuration</h1>
        <div className="flex items-center space-x-3">
          {lastSaved && (
            <div className="flex items-center text-sm text-gray-500">
              <CheckCircle className="w-4 h-4 mr-1 text-green-500" />
              Last saved: {lastSaved.toLocaleTimeString()}
            </div>
          )}
          <input
            type="file"
            accept=".json"
            onChange={restoreConfiguration}
            className="hidden"
            id="restore-config"
          />
          <label
            htmlFor="restore-config"
            className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 cursor-pointer"
          >
            <Upload className="w-4 h-4 mr-2" />
            Restore
          </label>
          <button
            onClick={backupConfiguration}
            className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
          >
            <Download className="w-4 h-4 mr-2" />
            Backup
          </button>
          <button
            onClick={saveConfiguration}
            disabled={saving}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
          >
            {saving ? (
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
            ) : (
              <Save className="w-4 h-4 mr-2" />
            )}
            Save Changes
          </button>
        </div>
      </div>

      <div className="space-y-6">
        {config.map((section, sectionIndex) => (
          <div key={section.name} className="bg-white shadow rounded-lg">
            <div className="px-6 py-4 border-b border-gray-200">
              <h3 className="text-lg font-medium text-gray-900 flex items-center">
                <Settings className="w-5 h-5 mr-2 text-gray-400" />
                {section.name}
              </h3>
            </div>
            <div className="px-6 py-4">
              <div className="space-y-6">
                {section.settings.map((setting, settingIndex) => (
                  <div key={setting.key} className="grid grid-cols-1 md:grid-cols-3 gap-4 items-start">
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        {setting.label}
                      </label>
                      {setting.description && (
                        <p className="mt-1 text-sm text-gray-500">{setting.description}</p>
                      )}
                    </div>
                    <div className="md:col-span-2">
                      {setting.type === 'text' && (
                        <input
                          type="text"
                          value={setting.value as string}
                          onChange={(e) => updateSetting(sectionIndex, settingIndex, e.target.value)}
                          className="block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                        />
                      )}
                      {setting.type === 'number' && (
                        <input
                          type="number"
                          value={setting.value as number}
                          onChange={(e) => updateSetting(sectionIndex, settingIndex, parseInt(e.target.value))}
                          className="block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                        />
                      )}
                      {setting.type === 'boolean' && (
                        <div className="flex items-center">
                          <input
                            type="checkbox"
                            checked={setting.value as boolean}
                            onChange={(e) => updateSetting(sectionIndex, settingIndex, e.target.checked)}
                            className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                          />
                          <span className="ml-2 text-sm text-gray-700">
                            {setting.value ? 'Enabled' : 'Disabled'}
                          </span>
                        </div>
                      )}
                      {setting.type === 'select' && (
                        <select
                          value={setting.value as string}
                          onChange={(e) => updateSetting(sectionIndex, settingIndex, e.target.value)}
                          className="block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                        >
                          {setting.options?.map((option) => (
                            <option key={option} value={option}>
                              {option.toUpperCase()}
                            </option>
                          ))}
                        </select>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
