import { useState } from 'react'
import { 
  Plus, 
  Settings, 
  Activity, 
  AlertTriangle, 
  CheckCircle,
  HardDrive,
  Network,
  Database,
  Cloud
} from 'lucide-react'

interface StorageBackend {
  name: string
  type: 'local' | 'nas' | 'san' | 'object'
  status: 'healthy' | 'warning' | 'critical'
  endpoint: string
  totalSpace: string
  usedSpace: string
  freeSpace: string
  encryption: boolean
  raid?: string
}

export function StorageBackends() {
  const [backends] = useState<StorageBackend[]>([
    {
      name: 'Local RAID-6',
      type: 'local',
      status: 'healthy',
      endpoint: '/data/minio',
      totalSpace: '1.0 TB',
      usedSpace: '650 GB',
      freeSpace: '350 GB',
      encryption: true,
      raid: 'RAID-6'
    },
    {
      name: 'NAS Primary',
      type: 'nas',
      status: 'healthy',
      endpoint: 'nfs://nas.local/storage',
      totalSpace: '2.0 TB',
      usedSpace: '900 GB',
      freeSpace: '1.1 TB',
      encryption: true
    },
    {
      name: 'SAN iSCSI',
      type: 'san',
      status: 'warning',
      endpoint: 'iscsi://san.local:3260',
      totalSpace: '5.0 TB',
      usedSpace: '4.2 TB',
      freeSpace: '800 GB',
      encryption: true
    },
    {
      name: 'MinIO Object',
      type: 'object',
      status: 'healthy',
      endpoint: 'minio://minio.local:9000',
      totalSpace: '10.0 TB',
      usedSpace: '2.5 TB',
      freeSpace: '7.5 TB',
      encryption: true
    }
  ])

  const [showCreateModal, setShowCreateModal] = useState(false)

  const getBackendIcon = (type: string) => {
    switch (type) {
      case 'local': return HardDrive
      case 'nas': return Network
      case 'san': return Database
      case 'object': return Cloud
      default: return HardDrive
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'healthy': return 'text-green-500'
      case 'warning': return 'text-yellow-500'
      case 'critical': return 'text-red-500'
      default: return 'text-gray-500'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'healthy': return CheckCircle
      case 'warning': return AlertTriangle
      case 'critical': return AlertTriangle
      default: return Activity
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Storage Backends</h1>
        <button
          onClick={() => setShowCreateModal(true)}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <Plus className="w-4 h-4 mr-2" />
          Add Backend
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {backends.map((backend, index) => {
          const Icon = getBackendIcon(backend.type)
          const StatusIcon = getStatusIcon(backend.status)
          
          return (
            <div key={index} className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <Icon className="h-8 w-8 text-gray-400" />
                    <div className="ml-3">
                      <h3 className="text-lg font-medium text-gray-900">{backend.name}</h3>
                      <p className="text-sm text-gray-500 capitalize">{backend.type} Storage</p>
                    </div>
                  </div>
                  <div className="flex items-center">
                    <StatusIcon className={`h-5 w-5 ${getStatusColor(backend.status)}`} />
                    <button className="ml-2 p-1 rounded-full text-gray-400 hover:text-gray-500">
                      <Settings className="h-4 w-4" />
                    </button>
                  </div>
                </div>

                <div className="mt-4">
                  <div className="text-sm text-gray-500 mb-2">Endpoint</div>
                  <div className="text-sm font-mono text-gray-900 bg-gray-50 p-2 rounded">
                    {backend.endpoint}
                  </div>
                </div>

                <div className="mt-4 grid grid-cols-3 gap-4 text-sm">
                  <div>
                    <div className="text-gray-500">Total</div>
                    <div className="font-medium text-gray-900">{backend.totalSpace}</div>
                  </div>
                  <div>
                    <div className="text-gray-500">Used</div>
                    <div className="font-medium text-gray-900">{backend.usedSpace}</div>
                  </div>
                  <div>
                    <div className="text-gray-500">Free</div>
                    <div className="font-medium text-gray-900">{backend.freeSpace}</div>
                  </div>
                </div>

                <div className="mt-4 flex items-center justify-between text-sm">
                  <div className="flex items-center">
                    <span className="text-gray-500">Encryption:</span>
                    <span className={`ml-1 ${backend.encryption ? 'text-green-600' : 'text-red-600'}`}>
                      {backend.encryption ? 'Enabled' : 'Disabled'}
                    </span>
                  </div>
                  {backend.raid && (
                    <div className="text-gray-500">
                      {backend.raid}
                    </div>
                  )}
                </div>

                <div className="mt-4 flex space-x-2">
                  <button className="flex-1 bg-blue-50 text-blue-700 text-sm font-medium py-2 px-3 rounded hover:bg-blue-100">
                    View Details
                  </button>
                  <button className="flex-1 bg-gray-50 text-gray-700 text-sm font-medium py-2 px-3 rounded hover:bg-gray-100">
                    Test Connection
                  </button>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {showCreateModal && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
            <div className="mt-3">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Add Storage Backend</h3>
              <form className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Name</label>
                  <input
                    type="text"
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Backend name"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Type</label>
                  <select className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500">
                    <option value="local">Local Storage</option>
                    <option value="nas">NAS Storage</option>
                    <option value="san">SAN Storage</option>
                    <option value="object">Object Storage</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Endpoint</label>
                  <input
                    type="text"
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                    placeholder="Storage endpoint URL"
                  />
                </div>
                <div className="flex items-center">
                  <input
                    type="checkbox"
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                  <label className="ml-2 block text-sm text-gray-900">
                    Enable encryption
                  </label>
                </div>
                <div className="flex justify-end space-x-3 pt-4">
                  <button
                    type="button"
                    onClick={() => setShowCreateModal(false)}
                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                  >
                    Create Backend
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
