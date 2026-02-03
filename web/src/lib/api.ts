import { authService } from '@/services/auth';

export async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = authService.getToken();
  
  const response = await fetch(`/api${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    },
  });

  if (!response.ok) {
    if (response.status === 401) {
      // Unauthorized - clear auth and redirect to login
      authService.logout();
      window.location.href = '/login';
      throw new Error('Session expired. Please login again.');
    }
    
    const error = await response.json().catch(() => ({ error: response.statusText }));
    throw new Error(error.error || `API error: ${response.statusText}`);
  }

  return response.json();
}
