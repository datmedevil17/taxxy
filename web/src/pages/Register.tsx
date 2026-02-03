import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { authService } from '@/services/auth';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { ArrowLeft, Car, User } from 'lucide-react';
import { cn } from '@/lib/utils';

export default function Register() {
  const navigate = useNavigate();
  const [role, setRole] = useState<'rider' | 'driver'>('rider');
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    name: '',
    phone: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const result = await authService.register({
        ...formData,
        role,
      });
      // Redirect to profile setup
      navigate(result.role === 'rider' ? '/rider/setup' : '/driver/setup');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-zinc-900 via-black to-zinc-900 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <Link to="/" className="inline-block mb-6">
          <Button variant="ghost" size="icon" className="text-white hover:bg-white/10">
            <ArrowLeft className="h-5 w-5" />
          </Button>
        </Link>

        <Card className="bg-white/5 backdrop-blur-lg border-white/10">
          <CardContent className="p-8">
            <div className="text-center mb-4">
              <img src="/favicon.png" alt="Taxxy" className="h-14 w-auto mx-auto" />
            </div>
            <h1 className="text-3xl font-bold text-white mb-2 text-center">Create Account</h1>
            <p className="text-zinc-400 mb-6 text-center">Join Taxxy today</p>

            {/* Role Selection */}
            <div className="grid grid-cols-2 gap-4 mb-6">
              <button
                type="button"
                onClick={() => setRole('rider')}
                className={cn(
                  'p-4 rounded-lg border-2 transition-all flex flex-col items-center gap-2',
                  role === 'rider'
                    ? 'border-white bg-white/10 text-white'
                    : 'border-white/20 text-zinc-400 hover:border-white/40'
                )}
              >
                <User className="h-6 w-6" />
                <span className="font-medium">Rider</span>
              </button>
              <button
                type="button"
                onClick={() => setRole('driver')}
                className={cn(
                  'p-4 rounded-lg border-2 transition-all flex flex-col items-center gap-2',
                  role === 'driver'
                    ? 'border-white bg-white/10 text-white'
                    : 'border-white/20 text-zinc-400 hover:border-white/40'
                )}
              >
                <Car className="h-6 w-6" />
                <span className="font-medium">Driver</span>
              </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-zinc-300 mb-2">Name</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-white/50"
                  placeholder="Your full name"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-zinc-300 mb-2">Email</label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-white/50"
                  placeholder="you@example.com"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-zinc-300 mb-2">Phone</label>
                <input
                  type="tel"
                  value={formData.phone}
                  onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-white/50"
                  placeholder="+1234567890"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-zinc-300 mb-2">Password</label>
                <input
                  type="password"
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-white/50"
                  placeholder="••••••••"
                  required
                  minLength={6}
                />
              </div>

              {error && (
                <div className="bg-red-500/10 border border-red-500/50 rounded-lg p-3">
                  <p className="text-red-400 text-sm">{error}</p>
                </div>
              )}

              <Button
                type="submit"
                className="w-full bg-white text-black hover:bg-zinc-200 h-12 text-lg font-semibold"
                disabled={loading}
              >
                {loading ? 'Creating Account...' : `Sign Up as ${role === 'rider' ? 'Rider' : 'Driver'}`}
              </Button>
            </form>

            <div className="mt-6 text-center">
              <p className="text-zinc-400 text-sm">
                Already have an account?{' '}
                <Link to="/login" className="text-white hover:underline font-medium">
                  Sign in
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
