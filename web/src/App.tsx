import { Routes, Route } from 'react-router-dom'
import Welcome from '@/pages/Welcome'
import Rider from '@/pages/Rider'
import Driver from '@/pages/Driver'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Welcome />} />
      <Route path="/rider" element={<Rider />} />
      <Route path="/driver" element={<Driver />} />
    </Routes>
  )
}

export default App
