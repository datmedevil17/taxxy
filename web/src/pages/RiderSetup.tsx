import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { riderService } from '@/services/rider';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { User, CreditCard, Check, ArrowLeft } from 'lucide-react';

export default function RiderSetup() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    phone: '',
    payment_method: 'credit_card',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await riderService.createProfile(formData);
      navigate('/rider');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Profile creation failed');
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
            <div className="text-center mb-6">
              <img src="/favicon.png" alt="Taxxy" className="h-14 w-auto mx-auto" />
            </div>
            <div className="flex items-center gap-3 mb-6">
              <div className="h-12 w-12 bg-white/10 rounded-full flex items-center justify-center">
                <User className="h-6 w-6 text-white" />
              </div>
              <div>
                <h1 className="text-2xl font-bold text-white">Complete Your Profile</h1>
                <p className="text-zinc-400 text-sm">Set up your rider profile to get started</p>
              </div>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-zinc-300 mb-2">Display Name</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-white/50"
                  placeholder="How should drivers call you?"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-zinc-300 mb-2">Phone Number</label>
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
                <label className="block text-sm font-medium text-zinc-300 mb-2">Preferred Payment</label>
                <div className="grid grid-cols-2 gap-3">
                  {['credit_card', 'cash'].map((method) => (
                    <button
                      key={method}
                      type="button"
                      onClick={() => setFormData({ ...formData, payment_method: method })}
                      className={`p-3 rounded-lg border flex items-center gap-2 transition-all ${
                        formData.payment_method === method
                          ? 'border-white bg-white/10 text-white'
                          : 'border-white/20 text-zinc-400 hover:border-white/40'
                      }`}
                    >
                      <CreditCard className="h-4 w-4" />
                      <span className="text-sm capitalize">{method.replace('_', ' ')}</span>
                      {formData.payment_method === method && <Check className="h-4 w-4 ml-auto" />}
                    </button>
                  ))}
                </div>
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
                {loading ? 'Setting up...' : 'Complete Setup'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
