import { Routes, Route } from 'react-router-dom';
import Welcome from '@/pages/Welcome';
import Login from '@/pages/Login';
import Register from '@/pages/Register';
import Rider from '@/pages/Rider';
import RiderSetup from '@/pages/RiderSetup';
import Driver from '@/pages/Driver';
import DriverSetup from '@/pages/DriverSetup';
import ProtectedRoute from '@/components/ProtectedRoute';

function App() {
  return (
    <Routes>
      {/* Public Routes */}
      <Route path="/" element={<Welcome />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      {/* Protected Rider Routes */}
      <Route
        path="/rider"
        element={
          <ProtectedRoute requiredRole="rider">
            <Rider />
          </ProtectedRoute>
        }
      />
      <Route
        path="/rider/setup"
        element={
          <ProtectedRoute requiredRole="rider">
            <RiderSetup />
          </ProtectedRoute>
        }
      />

      {/* Protected Driver Routes */}
      <Route
        path="/driver"
        element={
          <ProtectedRoute requiredRole="driver">
            <Driver />
          </ProtectedRoute>
        }
      />
      <Route
        path="/driver/setup"
        element={
          <ProtectedRoute requiredRole="driver">
            <DriverSetup />
          </ProtectedRoute>
        }
      />
    </Routes>
  );
}

export default App;
