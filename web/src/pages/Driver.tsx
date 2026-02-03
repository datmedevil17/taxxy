import { useState, useEffect } from "react";
import { useNavigate, Link } from "react-router-dom";
import { driverService } from "@/services/driver";
import { tripService } from "@/services/trip";
import { authService } from "@/services/auth";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { ArrowLeft, CheckCircle2, MapPin, User, XCircle, Play, LogOut, Power } from "lucide-react";
import { cn } from "@/lib/utils";
import type { DriverProfile, Trip } from "@/types";

export default function Driver() {
  const navigate = useNavigate();
  const [profile, setProfile] = useState<DriverProfile | null>(null);
  const [availability, setAvailability] = useState<'offline' | 'available' | 'busy'>('offline');
  const [currentTrip, setCurrentTrip] = useState<Trip | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // Load profile on mount
  useEffect(() => {
    const loadProfile = async () => {
      try {
        const p = await driverService.getProfile();
        setProfile(p);
        setAvailability(p.status || 'offline');
      } catch {
        navigate('/driver/setup');
      }
    };
    loadProfile();
  }, [navigate]);

  // Poll for available trips when available
  useEffect(() => {
    if (availability !== 'available' || currentTrip) return;

    // In a real app, this would be WebSocket or server-sent events
    // For now, we'll rely on the trip being assigned by the backend
    const checkForTrips = async () => {
      // This is a placeholder - the backend would push available trips
      // For demo, trips appear when a rider requests within the same system
    };

    const interval = setInterval(checkForTrips, 5000);
    return () => clearInterval(interval);
  }, [availability, currentTrip]);

  // Poll current trip status
  useEffect(() => {
    if (!currentTrip) return;

    const interval = setInterval(async () => {
      try {
        const trip = await tripService.getTrip(currentTrip.trip_id);
        setCurrentTrip(trip);
        
        if (trip.status === 'completed' || trip.status === 'cancelled') {
          setCurrentTrip(null);
          setAvailability('available');
        }
      } catch (err) {
        console.error('Failed to poll trip:', err);
      }
    }, 3000);

    return () => clearInterval(interval);
  }, [currentTrip]);

  const toggleAvailability = async () => {
    if (currentTrip) return; // Can't toggle during active trip
    
    setLoading(true);
    setError("");
    
    try {
      const newStatus = availability === 'available' ? 'offline' : 'available';
      // Get current location (mock for now)
      const location = { lat: 37.7749, lng: -122.4194 };
      
      await driverService.updateAvailability({
        status: newStatus,
        location,
      });
      
      setAvailability(newStatus);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update availability');
    } finally {
      setLoading(false);
    }
  };

  const acceptTrip = async (tripId: string) => {
    setLoading(true);
    setError("");
    
    try {
      await driverService.acceptRide(tripId);
      const trip = await tripService.getTrip(tripId);
      setCurrentTrip(trip);
      setAvailability('busy');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to accept trip');
    } finally {
      setLoading(false);
    }
  };

  const startTrip = async () => {
    if (!currentTrip) return;
    
    setLoading(true);
    setError("");
    
    try {
      await driverService.startRide(currentTrip.trip_id);
      const trip = await tripService.getTrip(currentTrip.trip_id);
      setCurrentTrip(trip);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to start trip');
    } finally {
      setLoading(false);
    }
  };

  const completeTrip = async () => {
    if (!currentTrip) return;
    
    setLoading(true);
    setError("");
    
    try {
      const result = await driverService.completeRide(currentTrip.trip_id);
      alert(`Trip completed! You earned $${result.fare.toFixed(2)}`);
      setCurrentTrip(null);
      setAvailability('available');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to complete trip');
    } finally {
      setLoading(false);
    }
  };

  const cancelTrip = () => {
    setCurrentTrip(null);
    setAvailability('available');
  };

  const handleLogout = () => {
    authService.logout();
    navigate('/');
  };

  if (!profile) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-zinc-950">
        <div className="animate-spin h-8 w-8 border-4 border-white border-t-transparent rounded-full" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-950 text-white p-4 font-sans">
      <div className="max-w-md mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between py-4">
          <div className="flex items-center gap-4">
            <Link to="/">
              <Button variant="ghost" size="icon" className="text-white hover:text-white/80 hover:bg-white/10">
                <ArrowLeft className="h-6 w-6" />
              </Button>
            </Link>
            <h1 className="text-xl font-bold">Driver Portal</h1>
          </div>
          <Button variant="ghost" size="icon" className="text-white hover:text-white/80 hover:bg-white/10" onClick={handleLogout}>
            <LogOut className="h-5 w-5" />
          </Button>
        </div>

        {/* Driver Info */}
        <Card className="bg-zinc-900 border-zinc-800 text-white">
          <CardContent className="p-4">
            <div className="flex items-center gap-4">
              <div className="h-12 w-12 bg-zinc-800 rounded-full flex items-center justify-center">
                <User className="h-6 w-6 text-zinc-400" />
              </div>
              <div className="flex-1">
                <p className="font-bold">{profile.name}</p>
                <p className="text-sm text-zinc-400">{profile.vehicle_brand} • {profile.vehicle_plate}</p>
              </div>
              <div className={cn(
                "px-3 py-1 rounded-full text-xs font-medium",
                availability === 'available' ? "bg-green-500/20 text-green-400" :
                availability === 'busy' ? "bg-yellow-500/20 text-yellow-400" :
                "bg-zinc-800 text-zinc-400"
              )}>
                {availability}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Availability Toggle */}
        {!currentTrip && (
          <Button
            onClick={toggleAvailability}
            disabled={loading}
            className={cn(
              "w-full h-14 text-lg font-semibold gap-3",
              availability === 'available' 
                ? "bg-red-500 hover:bg-red-600" 
                : "bg-green-500 hover:bg-green-600"
            )}
          >
            <Power className="h-5 w-5" />
            {availability === 'available' ? 'Go Offline' : 'Go Online'}
          </Button>
        )}

        {/* Error Display */}
        {error && (
          <div className="bg-red-500/10 border border-red-500/50 rounded-lg p-3">
            <p className="text-red-400 text-sm">{error}</p>
          </div>
        )}

        {/* Active Trip */}
        {currentTrip && (
          <Card className="bg-white text-black border-none">
            <CardContent className="p-6 space-y-4">
              <div className="flex justify-between items-start">
                <div>
                  <p className="font-bold text-lg">
                    {currentTrip.status === 'accepted' ? 'Pickup Rider' : 
                     currentTrip.status === 'in_progress' ? 'Trip in Progress' : 
                     'Active Trip'}
                  </p>
                  <p className="text-sm text-zinc-500">{currentTrip.vehicle_type}</p>
                </div>
                <div className="font-bold text-xl">${currentTrip.estimated_fare.toFixed(2)}</div>
              </div>
              
              <div className="space-y-4 py-2">
                <div className="flex items-center gap-3 text-sm">
                  <MapPin className="h-4 w-4 text-green-600 shrink-0" />
                  <span className="truncate">Pickup: {currentTrip.pickup.latitude.toFixed(4)}, {currentTrip.pickup.longitude.toFixed(4)}</span>
                </div>
                <div className="pl-2 ml-[5px] border-l-2 border-zinc-100 h-4" />
                <div className="flex items-center gap-3 text-sm">
                  <MapPin className="h-4 w-4 text-red-600 shrink-0" />
                  <span className="truncate">Dropoff: {currentTrip.dropoff.latitude.toFixed(4)}, {currentTrip.dropoff.longitude.toFixed(4)}</span>
                </div>
              </div>

              {currentTrip.status === 'accepted' && (
                <div className="grid grid-cols-2 gap-3 pt-2">
                  <Button 
                    variant="outline" 
                    onClick={cancelTrip}
                    className="border-zinc-200 hover:bg-zinc-50 text-black gap-2"
                  >
                    <XCircle className="h-4 w-4" /> Cancel
                  </Button>
                  <Button 
                    onClick={startTrip}
                    disabled={loading}
                    className="bg-green-500 text-white hover:bg-green-600 gap-2"
                  >
                    <Play className="h-4 w-4" /> Start Ride
                  </Button>
                </div>
              )}

              {currentTrip.status === 'in_progress' && (
                <Button 
                  onClick={completeTrip}
                  disabled={loading}
                  className="w-full bg-black text-white hover:bg-zinc-800 gap-2"
                >
                  <CheckCircle2 className="h-4 w-4" /> Complete Trip
                </Button>
              )}
            </CardContent>
          </Card>
        )}

        {/* No Active Trip - Waiting for requests */}
        {!currentTrip && availability === 'available' && (
          <Card className="bg-zinc-900 border-zinc-800">
            <CardContent className="p-8 text-center">
              <div className="h-16 w-16 bg-zinc-800 rounded-full flex items-center justify-center mx-auto mb-4">
                <div className="h-3 w-3 bg-green-500 rounded-full animate-pulse" />
              </div>
              <h3 className="text-white font-medium text-lg">Waiting for ride requests</h3>
              <p className="text-zinc-500 text-sm mt-1">You'll see new requests here</p>
              <p className="text-zinc-600 text-xs mt-4">
                Note: For demo, use the Rider app to request a ride, 
                then manually enter the trip ID to accept
              </p>
              
              {/* Demo: Manual trip acceptance input */}
              <div className="mt-4 flex gap-2">
                <input
                  type="text"
                  placeholder="Enter Trip ID"
                  className="flex-1 px-3 py-2 bg-zinc-800 border border-zinc-700 rounded text-white text-sm"
                  id="tripIdInput"
                />
                <Button
                  size="sm"
                  onClick={() => {
                    const input = document.getElementById('tripIdInput') as HTMLInputElement;
                    if (input?.value) acceptTrip(input.value);
                  }}
                  disabled={loading}
                >
                  Accept
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Offline State */}
        {!currentTrip && availability === 'offline' && (
          <div className="text-center py-12 text-zinc-500">
            <p>You're currently offline.</p>
            <p className="text-xs mt-1 text-zinc-600">Go online to start receiving ride requests.</p>
          </div>
        )}
      </div>
    </div>
  );
}
