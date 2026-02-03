import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { driverService } from '@/services/driver';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Car, Check, ArrowLeft } from 'lucide-react';
import type { VehicleType } from '@/types';

const VEHICLE_TYPES: { id: VehicleType; name: string; desc: string }[] = [
  { id: 'sedan', name: 'Sedan', desc: 'Standard 4-door car' },
  { id: 'suv', name: 'SUV', desc: 'Sport utility vehicle' },
  { id: 'van', name: 'Van', desc: 'Passenger van' },
  { id: 'luxury', name: 'Luxury', desc: 'Premium vehicle' },
];

export default function DriverSetup() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    phone: '',
    vehicle_type: 'sedan' as VehicleType,
    vehicle_plate: '',
    vehicle_brand: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await driverService.createProfile(formData);
      navigate('/driver');
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
                <Car className="h-6 w-6 text-white" />
              </div>
              <div>
                <h1 className="text-2xl font-bold text-white">Driver Profile</h1>
                <p className="text-zinc-400 text-sm">Set up your vehicle details</p>
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
                  placeholder="Your name"
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
                <label className="block text-sm font-medium text-zinc-300 mb-2">Vehicle Type</label>
                <div className="grid grid-cols-2 gap-2">
                  {VEHICLE_TYPES.map((v) => (
                    <button
                      key={v.id}
                      type="button"
                      onClick={() => setFormData({ ...formData, vehicle_type: v.id })}
                      className={`p-3 rounded-lg border text-left transition-all ${
                        formData.vehicle_type === v.id
                          ? 'border-white bg-white/10 text-white'
                          : 'border-white/20 text-zinc-400 hover:border-white/40'
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <span className="font-medium text-sm">{v.name}</span>
                        {formData.vehicle_type === v.id && <Check className="h-4 w-4" />}
                      </div>
                      <span className="text-xs text-zinc-500">{v.desc}</span>
                    </button>
                  ))}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-zinc-300 mb-2">Vehicle Brand</label>
                <input
                  type="text"
                  value={formData.vehicle_brand}
                  onChange={(e) => setFormData({ ...formData, vehicle_brand: e.target.value })}
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-white/50"
                  placeholder="Toyota, Honda, etc."
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-zinc-300 mb-2">License Plate</label>
                <input
                  type="text"
                  value={formData.vehicle_plate}
                  onChange={(e) => setFormData({ ...formData, vehicle_plate: e.target.value.toUpperCase() })}
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-white/50 uppercase"
                  placeholder="ABC123"
                  required
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
                {loading ? 'Setting up...' : 'Start Driving'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
