import { useState, useEffect, useCallback } from "react";
import { useNavigate, Link } from "react-router-dom";
import { riderService } from "@/services/rider";
import { paymentService } from "@/services/payment";
import { tripService } from "@/services/trip";
import { authService } from "@/services/auth";
import Map from "@/components/Map";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { ArrowLeft, Car, Clock, CreditCard, LogOut, RefreshCw, MapPin } from "lucide-react";
import { cn } from "@/lib/utils";
import type { VehicleType, RiderProfile, Trip } from "@/types";

const VEHICLES: { id: VehicleType; name: string; time: string }[] = [
  { id: 'sedan', name: 'Sedan', time: '3 min' },
  { id: 'suv', name: 'SUV', time: '5 min' },
  { id: 'van', name: 'Van', time: '8 min' },
  { id: 'luxury', name: 'Luxury', time: '10 min' },
];

export default function Rider() {
  const navigate = useNavigate();
  const [profile, setProfile] = useState<RiderProfile | null>(null);
  const [pickup, setPickup] = useState<{ lat: number; lng: number } | null>(null);
  const [dropoff, setDropoff] = useState<{ lat: number; lng: number } | null>(null);
  const [selectedVehicle, setSelectedVehicle] = useState<VehicleType | null>(null);
  const [estimatedFare, setEstimatedFare] = useState<number | null>(null);
  const [currentTrip, setCurrentTrip] = useState<Trip | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [status, setStatus] = useState<'idle' | 'estimating' | 'requesting' | 'waiting' | 'accepted' | 'in_progress'>('idle');

  // Load profile on mount
  useEffect(() => {
    const loadProfile = async () => {
      try {
        const p = await riderService.getProfile();
        setProfile(p);
      } catch {
        // Profile doesn't exist, redirect to setup
        navigate('/rider/setup');
      }
    };
    loadProfile();
  }, [navigate]);

  // Poll trip status when we have an active trip
  useEffect(() => {
    if (!currentTrip) return;

    const interval = setInterval(async () => {
      try {
        const trip = await tripService.getTrip(currentTrip.trip_id);
        setCurrentTrip(trip);
        
        if (trip.status === 'accepted') setStatus('accepted');
        else if (trip.status === 'in_progress') setStatus('in_progress');
        else if (trip.status === 'completed' || trip.status === 'cancelled') {
          setCurrentTrip(null);
          setStatus('idle');
          resetRide();
        }
      } catch (err) {
        console.error('Failed to poll trip status:', err);
      }
    }, 3000);

    return () => clearInterval(interval);
  }, [currentTrip]);

  const handleMapClick = (lat: number, lng: number) => {
    if (status !== 'idle' && status !== 'estimating') return;

    if (!pickup) {
      setPickup({ lat, lng });
    } else if (!dropoff) {
      setDropoff({ lat, lng });
    } else {
      resetRide();
      setPickup({ lat, lng });
    }
  };

  const resetRide = () => {
    setPickup(null);
    setDropoff(null);
    setSelectedVehicle(null);
    setEstimatedFare(null);
    setCurrentTrip(null);
    setStatus('idle');
    setError("");
  };

  // Calculate fare when locations and vehicle are selected
  const calculateFare = useCallback(async () => {
    if (!pickup || !dropoff || !selectedVehicle) return;
    
    setStatus('estimating');
    try {
      const result = await paymentService.calculateFare({
        pickup_lat: pickup.lat,
        pickup_lng: pickup.lng,
        dropoff_lat: dropoff.lat,
        dropoff_lng: dropoff.lng,
        vehicle_type: selectedVehicle,
      });
      setEstimatedFare(result.total_fare);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to calculate fare');
    }
    setStatus('idle');
  }, [pickup, dropoff, selectedVehicle]);

  useEffect(() => {
    if (pickup && dropoff && selectedVehicle) {
      calculateFare();
    }
  }, [pickup, dropoff, selectedVehicle, calculateFare]);

  const handleRequestRide = async () => {
    if (!pickup || !dropoff || !selectedVehicle) return;

    setLoading(true);
    setError("");
    setStatus('requesting');

    try {
      const result = await riderService.requestRide({
        pickup: { lat: pickup.lat, lng: pickup.lng },
        dropoff: { lat: dropoff.lat, lng: dropoff.lng },
        vehicle_type: selectedVehicle,
      });
      
      // Fetch full trip details
      const trip = await tripService.getTrip(result.trip_id);
      setCurrentTrip(trip);
      setStatus('waiting');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to request ride');
      setStatus('idle');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    authService.logout();
    navigate('/');
  };

  const handlePayment = async () => {
    if (!currentTrip) return;
    
    setLoading(true);
    try {
      await paymentService.processPayment({
        trip_id: currentTrip.trip_id,
        amount: currentTrip.final_fare || currentTrip.estimated_fare,
        payment_method: profile?.payment_method || 'credit_card',
      });
      alert('Payment successful!');
      resetRide();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Payment failed');
    } finally {
      setLoading(false);
    }
  };

  if (!profile) {
    return (
      <div className="h-screen flex items-center justify-center bg-zinc-900">
        <div className="animate-spin h-8 w-8 border-4 border-white border-t-transparent rounded-full" />
      </div>
    );
  }

  return (
    <div className="h-screen flex flex-col md:flex-row bg-white text-black font-sans overflow-hidden">
      {/* Map Area */}
      <div className="flex-1 relative h-[50vh] md:h-full">
        <div className="absolute top-4 left-4 z-[1000] flex gap-2">
          <Link to="/">
            <Button variant="outline" size="icon" className="bg-white shadow">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <Button variant="outline" size="icon" className="bg-white shadow" onClick={handleLogout}>
            <LogOut className="h-4 w-4" />
          </Button>
        </div>
        <Map 
          pickup={pickup} 
          dropoff={dropoff} 
          onMapClick={handleMapClick}
        />
      </div>

      {/* Sidebar */}
      <div className="w-full md:w-[400px] bg-white border-l shadow-xl flex flex-col h-[50vh] md:h-full z-10">
        <div className="p-6 border-b">
          <h2 className="text-2xl font-bold tracking-tight">Request a Ride</h2>
          <p className="text-zinc-500 text-sm mt-1">
            {!pickup ? "Tap map to set pickup" : 
             !dropoff ? "Tap map to set destination" : 
             status === 'waiting' ? "Finding your driver..." :
             status === 'accepted' ? "Driver is on the way!" :
             status === 'in_progress' ? "Enjoy your ride!" :
             "Select your vehicle"}
          </p>
        </div>

        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-3">
              <p className="text-red-600 text-sm">{error}</p>
            </div>
          )}

          {status === 'in_progress' && currentTrip ? (
            <div className="flex flex-col items-center justify-center space-y-6 py-10 animate-in fade-in slide-in-from-bottom-4">
              <div className="bg-green-500 text-white rounded-full p-4">
                <Car className="h-8 w-8" />
              </div>
              <div className="text-center space-y-1">
                <h3 className="text-xl font-bold">Ride in Progress</h3>
                <p className="text-zinc-500">Enjoy your ride!</p>
              </div>
              <div className="w-full bg-zinc-50 p-4 rounded-lg">
                <div className="flex justify-between items-center">
                  <span className="text-zinc-600">Estimated Fare</span>
                  <span className="font-bold text-lg">${currentTrip.estimated_fare.toFixed(2)}</span>
                </div>
              </div>
            </div>
          ) : status === 'accepted' && currentTrip ? (
            <div className="flex flex-col items-center justify-center space-y-6 py-10 animate-in fade-in slide-in-from-bottom-4">
              <div className="bg-black text-white rounded-full p-4">
                <Car className="h-8 w-8" />
              </div>
              <div className="text-center space-y-1">
                <h3 className="text-xl font-bold">Driver Arriving</h3>
                <p className="text-zinc-500">Your {selectedVehicle} is on the way.</p>
              </div>
              <div className="w-full bg-zinc-50 p-4 rounded-lg flex justify-between items-center">
                <div className="flex items-center gap-2">
                  <Clock className="h-4 w-4 text-zinc-500" />
                  <span>~5 min</span>
                </div>
                <div className="font-bold">${currentTrip.estimated_fare.toFixed(2)}</div>
              </div>
              <Button size="lg" className="w-full gap-2" onClick={handlePayment} disabled={loading}>
                <CreditCard className="h-4 w-4" /> Pay Now
              </Button>
            </div>
          ) : status === 'waiting' ? (
            <div className="flex flex-col items-center justify-center py-20 text-center space-y-4">
              <div className="h-16 w-16 bg-zinc-100 rounded-full flex items-center justify-center">
                <RefreshCw className="h-8 w-8 text-zinc-400 animate-spin" />
              </div>
              <div>
                <h3 className="font-medium text-lg">Connecting to drivers...</h3>
                <p className="text-sm text-zinc-500">Please wait while we find you a ride.</p>
              </div>
              <Button variant="outline" onClick={resetRide}>Cancel</Button>
            </div>
          ) : (
            <>
              {/* Location Summary */}
              {(pickup || dropoff) && (
                <Card className="border-zinc-100">
                  <CardContent className="p-4 space-y-2">
                    {pickup && (
                      <div className="flex items-center gap-2 text-sm">
                        <MapPin className="h-4 w-4 text-green-600" />
                        <span>Pickup: {pickup.lat.toFixed(4)}, {pickup.lng.toFixed(4)}</span>
                      </div>
                    )}
                    {dropoff && (
                      <div className="flex items-center gap-2 text-sm">
                        <MapPin className="h-4 w-4 text-red-600" />
                        <span>Dropoff: {dropoff.lat.toFixed(4)}, {dropoff.lng.toFixed(4)}</span>
                      </div>
                    )}
                  </CardContent>
                </Card>
              )}

              {/* Vehicle Selection */}
              <div className="space-y-3">
                {VEHICLES.map((v) => (
                  <Card 
                    key={v.id}
                    onClick={() => setSelectedVehicle(v.id)}
                    className={cn(
                      "cursor-pointer transition-all hover:bg-zinc-50 border-zinc-100 shadow-sm",
                      selectedVehicle === v.id ? "ring-2 ring-black border-transparent" : ""
                    )}
                  >
                    <CardContent className="flex items-center justify-between p-4">
                      <div className="flex items-center gap-4">
                        <div className="h-10 w-10 bg-zinc-100 rounded-full flex items-center justify-center">
                          <Car className="h-5 w-5 text-black" />
                        </div>
                        <div>
                          <div className="font-semibold">{v.name}</div>
                          <div className="text-xs text-zinc-500">{v.time} away</div>
                        </div>
                      </div>
                      {selectedVehicle === v.id && estimatedFare ? (
                        <div className="font-bold">${estimatedFare.toFixed(2)}</div>
                      ) : (
                        <div className="text-zinc-400 text-sm">Select</div>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </div>
            </>
          )}
        </div>

        {status === 'idle' && (
          <div className="p-4 border-t bg-zinc-50/50">
            <Button 
              className="w-full text-lg h-12" 
              disabled={!selectedVehicle || !pickup || !dropoff || loading}
              onClick={handleRequestRide}
            >
              {loading ? 'Requesting...' : `Confirm ${selectedVehicle ? VEHICLES.find(v => v.id === selectedVehicle)?.name : 'Ride'}`}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
