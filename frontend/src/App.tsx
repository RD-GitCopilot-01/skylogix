import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { Layout } from './components/Layout'
import { Dashboard } from './pages/Dashboard'
import { StorageBackends } from './pages/StorageBackends'
import { Configuration } from './pages/Configuration'
import { Monitoring } from './pages/Monitoring'
import { SNMP } from './pages/SNMP'
import { Database } from './pages/Database'
import './App.css'

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/storage" element={<StorageBackends />} />
          <Route path="/config" element={<Configuration />} />
          <Route path="/monitoring" element={<Monitoring />} />
          <Route path="/snmp" element={<SNMP />} />
          <Route path="/database" element={<Database />} />
        </Routes>
      </Layout>
    </Router>
  )
}

export default App
