import type { AuthResponse, LoginRequest, RegisterRequest } from '@/types';

const TOKEN_KEY = 'taxxy_token';
const USER_KEY = 'taxxy_user';

export const authService = {
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await fetch('/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Registration failed');
    }

    const result: AuthResponse = await response.json();
    this.setToken(result.token);
    this.setUser({ user_id: result.user_id, email: data.email, role: result.role });
    return result;
  },

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Login failed');
    }

    const result: AuthResponse = await response.json();
    this.setToken(result.token);
    this.setUser({ user_id: result.user_id, email: data.email, role: result.role });
    return result;
  },

  logout(): void {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  },

  getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY);
  },

  setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token);
  },

  getUser(): { user_id: string; email: string; role: 'rider' | 'driver' } | null {
    const user = localStorage.getItem(USER_KEY);
    return user ? JSON.parse(user) : null;
  },

  setUser(user: { user_id: string; email: string; role: 'rider' | 'driver' }): void {
    localStorage.setItem(USER_KEY, JSON.stringify(user));
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  },

  getRole(): 'rider' | 'driver' | null {
    const user = this.getUser();
    return user?.role || null;
  },
};
