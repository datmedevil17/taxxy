import { Navigate } from 'react-router-dom';
import { authService } from '@/services/auth';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: 'rider' | 'driver';
}

export default function ProtectedRoute({ children, requiredRole }: ProtectedRouteProps) {
  const isAuthenticated = authService.isAuthenticated();
  const userRole = authService.getRole();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRole && userRole !== requiredRole) {
    // Redirect to appropriate dashboard if wrong role
    return <Navigate to={userRole === 'rider' ? '/rider' : '/driver'} replace />;
  }

  return <>{children}</>;
}
