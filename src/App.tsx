import { useState, useEffect } from 'react'
import { 
  Play, 
  Pause, 
  Square, 
  SkipBack, 
  SkipForward, 
  Volume2, 
  Settings, 
  Bell, 
  BarChart3, 
  Shield, 
  RotateCcw,
  ChevronDown,
  ChevronRight,
  Radio,
  Camera,
  Radar,
  FileText,
  Headphones
} from 'lucide-react'
import './App.css'

function App() {
  const [currentTime, setCurrentTime] = useState(new Date())
  const [isPlaying, setIsPlaying] = useState(false)
  const [expandedNodes, setExpandedNodes] = useState<string[]>(['audio', 'h264', 'h265'])

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date())
    }, 1000)
    return () => clearInterval(timer)
  }, [])

  const toggleNode = (nodeId: string) => {
    setExpandedNodes(prev => 
      prev.includes(nodeId) 
        ? prev.filter(id => id !== nodeId)
        : [...prev, nodeId]
    )
  }

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('en-GB', { 
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  }

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('en-GB', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric'
    })
  }

  return (
    <div className="h-screen bg-gray-100 flex flex-col">
      {/* Top Navigation Bar */}
      <div className="bg-blue-600 text-white px-4 py-2 flex items-center justify-between">
        <div className="flex items-center space-x-6">
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-white rounded flex items-center justify-center">
              <span className="text-blue-600 font-bold text-sm">R</span>
            </div>
            <span className="font-semibold">Ricochet Catch - (v4.12) - Logged in as: admin on FAA Server B Licensed to: Ricochet</span>
          </div>
        </div>
        <div className="flex items-center space-x-4">
          <span className="text-sm">{formatDate(currentTime)}</span>
          <span className="text-xl font-mono">{formatTime(currentTime)}</span>
          <div className="text-sm">JOTRON</div>
        </div>
      </div>

      {/* Secondary Navigation */}
      <div className="bg-blue-500 text-white px-4 py-2 flex items-center space-x-6">
        <div className="flex items-center space-x-2 bg-blue-600 px-3 py-1 rounded">
          <RotateCcw size={16} />
          <span>Replay</span>
        </div>
        <div className="flex items-center space-x-2 hover:bg-blue-600 px-3 py-1 rounded cursor-pointer">
          <Shield size={16} />
          <span>Admin</span>
        </div>
        <div className="flex items-center space-x-2 hover:bg-blue-600 px-3 py-1 rounded cursor-pointer">
          <BarChart3 size={16} />
          <span>Reports</span>
        </div>
        <div className="flex items-center space-x-2 hover:bg-blue-600 px-3 py-1 rounded cursor-pointer">
          <Bell size={16} />
          <span>Alarms</span>
        </div>
        <div className="flex items-center space-x-2 hover:bg-blue-600 px-3 py-1 rounded cursor-pointer">
          <Settings size={16} />
          <span>Settings</span>
        </div>
      </div>

      <div className="flex flex-1 overflow-hidden">
        {/* Left Sidebar */}
        <div className="w-80 bg-gray-200 border-r border-gray-300 flex flex-col">
          <div className="p-2 border-b border-gray-300">
            <button className="text-sm text-blue-600 hover:underline">Remove all</button>
          </div>
          
          <div className="flex-1 overflow-y-auto p-2">
            {/* Audio Section */}
            <div className="mb-2">
              <div 
                className="flex items-center space-x-1 cursor-pointer hover:bg-gray-300 p-1 rounded"
                onClick={() => toggleNode('audio')}
              >
                {expandedNodes.includes('audio') ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                <Headphones size={16} />
                <span className="text-sm font-medium">Audio</span>
                <span className="text-xs text-gray-500 ml-auto">- 1 +</span>
              </div>
              {expandedNodes.includes('audio') && (
                <div className="ml-6 mt-1">
                  <div className="flex items-center space-x-1 p-1 hover:bg-gray-300 rounded cursor-pointer">
                    {expandedNodes.includes('speaking-clock') ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                    <Radio size={14} />
                    <span className="text-sm">(DID 35) Speaking Clock</span>
                    <span className="text-xs text-gray-500 ml-auto">- 1 +</span>
                  </div>
                  <div className="ml-4 p-1 text-sm text-gray-600">Speaking Clock</div>
                </div>
              )}
            </div>

            {/* H264 Section */}
            <div className="mb-2">
              <div 
                className="flex items-center space-x-1 cursor-pointer hover:bg-gray-300 p-1 rounded"
                onClick={() => toggleNode('h264')}
              >
                {expandedNodes.includes('h264') ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                <Camera size={16} />
                <span className="text-sm font-medium">H264</span>
                <span className="text-xs text-gray-500 ml-auto">- 7 +</span>
              </div>
              {expandedNodes.includes('h264') && (
                <div className="ml-6 mt-1 space-y-1">
                  <div className="text-sm text-gray-600 p-1">(DID 81) H264</div>
                  <div className="text-sm text-gray-600 p-1">H265</div>
                  <div className="text-sm text-gray-600 p-1">(DID 85) H265</div>
                </div>
              )}
            </div>

            {/* JVG801 Section */}
            <div className="mb-2">
              <div className="flex items-center space-x-1 cursor-pointer hover:bg-gray-300 p-1 rounded">
                <ChevronRight size={16} />
                <Radio size={16} />
                <span className="text-sm font-medium">JVG801</span>
                <span className="text-xs text-gray-500 ml-auto">- 1 +</span>
              </div>
            </div>

            {/* Radar Section */}
            <div className="mb-2">
              <div className="flex items-center space-x-1 cursor-pointer hover:bg-gray-300 p-1 rounded">
                <ChevronDown size={16} />
                <Radar size={16} />
                <span className="text-sm font-medium">Radar</span>
                <span className="text-xs text-gray-500 ml-auto">- 19 +</span>
              </div>
              <div className="ml-6 mt-1">
                <div className="text-sm text-gray-600 p-1">(DID 55) EthernetRAW - Radar</div>
              </div>
            </div>

            {/* Raw Section */}
            <div className="mb-2">
              <div className="flex items-center space-x-1 cursor-pointer hover:bg-gray-300 p-1 rounded">
                <ChevronDown size={16} />
                <FileText size={16} />
                <span className="text-sm font-medium">Raw</span>
                <span className="text-xs text-gray-500 ml-auto">- 3 +</span>
              </div>
              <div className="ml-6 mt-1">
                <div className="text-sm text-gray-600 p-1">(DID 57) EthernetRAW - RAW</div>
              </div>
            </div>

            {/* ReVue Section */}
            <div className="mb-2">
              <div className="flex items-center space-x-1 cursor-pointer hover:bg-gray-300 p-1 rounded">
                <ChevronDown size={16} />
                <Camera size={16} />
                <span className="text-sm font-medium">ReVue</span>
                <span className="text-xs text-gray-500 ml-auto">- 2 +</span>
              </div>
              <div className="ml-6 mt-1">
                <div className="text-sm text-gray-600 p-1">(DID 75) ReVue</div>
              </div>
            </div>

            {/* RTPAudio Section */}
            <div className="mb-2">
              <div className="flex items-center space-x-1 cursor-pointer hover:bg-gray-300 p-1 rounded">
                <ChevronDown size={16} />
                <Headphones size={16} />
                <span className="text-sm font-medium">RTPAudio</span>
                <span className="text-xs text-gray-500 ml-auto">- 213 +</span>
              </div>
              <div className="ml-6 mt-1">
                <div className="text-sm text-gray-600 p-1">(DID 21) EDJ37 version B</div>
                <div className="bg-gray-400 text-white text-sm p-2 rounded mt-1">
                  Channel 1 (UDP) [ CWP001 ]
                </div>
                <div className="text-sm text-gray-600 p-1 mt-1">Channel 2 (UDP) [ CWP002 ]</div>
              </div>
            </div>
          </div>
        </div>

        {/* Main Content Area */}
        <div className="flex-1 flex flex-col">
          {/* Control Bar */}
          <div className="bg-blue-400 text-white p-2 flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2">
                <button className="p-1 hover:bg-blue-500 rounded">
                  <SkipBack size={20} />
                </button>
                <button 
                  className="p-1 hover:bg-blue-500 rounded"
                  onClick={() => setIsPlaying(!isPlaying)}
                >
                  {isPlaying ? <Pause size={20} /> : <Play size={20} />}
                </button>
                <button className="p-1 hover:bg-blue-500 rounded">
                  <Square size={20} />
                </button>
                <button className="p-1 hover:bg-blue-500 rounded">
                  <SkipForward size={20} />
                </button>
              </div>
              <div className="text-sm">
                {formatDate(currentTime)} <span className="font-mono text-lg">{formatTime(currentTime)}</span> UTC
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <span className="bg-green-500 px-2 py-1 rounded text-sm">Online</span>
              <div className="text-sm">
                <span>Replay Speed: 1x</span>
                <span className="ml-4">View Range: 10 minutes</span>
              </div>
            </div>
          </div>

          {/* Tab Navigation */}
          <div className="bg-gray-300 border-b border-gray-400 flex">
            <div className="bg-white border-r border-gray-400 px-4 py-2 text-sm font-medium">Source panel</div>
            <div className="bg-gray-200 border-r border-gray-400 px-4 py-2 text-sm hover:bg-gray-100 cursor-pointer">Presenter panel</div>
            <div className="bg-gray-200 px-4 py-2 text-sm hover:bg-gray-100 cursor-pointer">Metadata panel</div>
          </div>

          {/* Main Panel Content */}
          <div className="flex-1 flex">
            {/* Audio Waveform Section */}
            <div className="flex-1 bg-gradient-to-b from-slate-500 via-slate-600 to-slate-700 p-4">
              <div className="mb-4">
                <div className="flex items-center space-x-2 mb-2">
                  <Volume2 size={16} className="text-white" />
                  <span className="text-sm font-medium text-white">Speaking Clock</span>
                  <span className="text-xs text-gray-300">00</span>
                </div>
                <div className="bg-slate-700 h-20 border border-slate-600 relative overflow-hidden">
                  {/* Waveform background */}
                  <div className="absolute inset-0 bg-slate-800"></div>
                  {/* More realistic waveform visualization */}
                  <div className="absolute top-0 left-0 w-full h-full flex items-center">
                    {Array.from({length: 300}).map((_, i) => (
                      <div 
                        key={i} 
                        className="bg-white opacity-90" 
                        style={{
                          width: '1px',
                          height: `${Math.abs(Math.sin(i * 0.05) * Math.cos(i * 0.02)) * 60 + 20}%`,
                          marginLeft: '0.5px'
                        }}
                      />
                    ))}
                  </div>
                  {/* Green timeline indicator */}
                  <div className="absolute left-1/2 top-0 w-0.5 h-full bg-green-400"></div>
                </div>
              </div>

              <div className="mb-4">
                <div className="flex items-center space-x-2 mb-2">
                  <Volume2 size={16} className="text-white" />
                  <span className="text-sm font-medium text-white">Channel 1 (UDP)</span>
                </div>
                <div className="bg-slate-700 h-20 border border-slate-600 relative overflow-hidden">
                  {/* Waveform background */}
                  <div className="absolute inset-0 bg-slate-800"></div>
                  {/* Audio waveform with circular peaks like original */}
                  <div className="absolute top-0 left-0 w-full h-full flex items-center justify-around">
                    {Array.from({length: 8}).map((_, i) => (
                      <div key={i} className="flex items-center justify-center">
                        <div 
                          className="bg-white rounded-full opacity-95 shadow-lg" 
                          style={{
                            width: `${30 + Math.random() * 15}px`,
                            height: `${30 + Math.random() * 15}px`,
                            background: 'radial-gradient(circle, rgba(255,255,255,0.95) 0%, rgba(220,220,220,0.8) 100%)'
                          }}
                        ></div>
                      </div>
                    ))}
                  </div>
                  {/* Green timeline indicator */}
                  <div className="absolute left-1/2 top-0 w-0.5 h-full bg-green-400"></div>
                </div>
              </div>

              {/* Video Thumbnails - More realistic placeholders */}
              <div className="grid grid-cols-4 gap-1 mt-4">
                {Array.from({length: 8}).map((_, i) => (
                  <div key={i} className="bg-slate-500 aspect-video border border-slate-600 flex items-center justify-center relative">
                    <div className="absolute inset-0 bg-gradient-to-br from-slate-400 to-slate-600"></div>
                    <div className="absolute inset-1 bg-slate-600 rounded-sm"></div>
                    <Camera size={16} className="text-slate-400 relative z-10" />
                  </div>
                ))}
              </div>

              {/* Timeline */}
              <div className="mt-4 bg-slate-400 p-2 border border-slate-500">
                <div className="flex justify-between text-xs text-slate-800 mb-1 font-mono">
                  <span>06:28:00</span>
                  <span>06:32:17</span>
                  <span>06:38:00</span>
                </div>
                <div className="h-3 bg-slate-500 relative border border-slate-600">
                  <div className="absolute left-1/2 top-0 w-0.5 h-full bg-green-500"></div>
                </div>
              </div>
            </div>

            {/* Right Panel */}
            <div className="w-64 bg-slate-500 border-l border-slate-600 p-4">
              <div className="space-y-6">
                <div>
                  <label className="text-sm font-medium block mb-2 text-white">Volume</label>
                  <div className="bg-slate-600 p-2 rounded border border-slate-700">
                    <input type="range" className="w-full accent-blue-400" defaultValue="75" />
                  </div>
                  <div className="text-xs text-slate-300 mt-1">Channel</div>
                </div>
                <div>
                  <label className="text-sm font-medium block mb-2 text-white">Channel</label>
                  <select className="w-full p-2 bg-slate-600 border border-slate-700 rounded text-sm text-white">
                    <option>Channel 1</option>
                    <option>Channel 2</option>
                  </select>
                </div>
                <div>
                  <label className="text-sm font-medium block mb-2 text-white">AVC</label>
                  <div className="bg-slate-600 p-2 rounded border border-slate-700">
                    <input type="range" className="w-full accent-blue-400" defaultValue="60" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App
