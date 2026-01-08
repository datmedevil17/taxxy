import { useState } from "react"
import { useStore, type VehicleType } from "@/store/useStore"
import Map from "@/components/Map"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { ArrowLeft, Car, Clock, CreditCard } from "lucide-react"
import { Link } from "react-router-dom"
import { cn } from "@/lib/utils"

const VEHICLES: { id: VehicleType; name: string; price: number; time: string }[] = [
  { id: 'sedan', name: 'Sedan', price: 25, time: '3 min' },
  { id: 'suv', name: 'SUV', price: 35, time: '5 min' },
  { id: 'van', name: 'Van', price: 45, time: '8 min' },
  { id: 'luxury', name: 'Luxury', price: 60, time: '10 min' },
]

export default function Rider() {
  const { currentTrip, setPickup, setDropoff, requestRide, status, resetRide } = useStore()
  const [selectedVehicle, setSelectedVehicle] = useState<VehicleType | null>(null)

  const handleMapClick = (lat: number, lng: number) => {
    if (status !== 'IDLE' && status !== 'REQUESTING') return

    if (!currentTrip?.pickup) {
      setPickup(lat, lng)
    } else if (!currentTrip?.dropoff) {
      setDropoff(lat, lng)
    } else {
      // Reset if both set
      resetRide()
      setPickup(lat, lng)
    }
  }

  const handleRequest = () => {
    if (selectedVehicle && currentTrip?.pickup && currentTrip?.dropoff) {
      const v = VEHICLES.find(v => v.id === selectedVehicle)!
      requestRide({
        ...currentTrip,
        id: Math.random().toString(36).substring(7),
        vehicle: selectedVehicle,
        price: v.price,
        time: v.time,
        distance: '5.2 km' // Mock distance
      })
    }
  }

  return (
    <div className="h-screen flex flex-col md:flex-row bg-white text-black font-sans overflow-hidden">
      {/* Map Area */}
      <div className="flex-1 relative h-[50vh] md:h-full">
        <div className="absolute top-4 left-4 z-[1000]">
          <Link to="/">
            <Button variant="outline" size="icon" className="bg-white shadow">
              <ArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
        </div>
        <Map 
          pickup={currentTrip?.pickup} 
          dropoff={currentTrip?.dropoff} 
          onMapClick={handleMapClick}
        />
      </div>

      {/* Sidebar */}
      <div className="w-full md:w-[400px] bg-white border-l shadow-xl flex flex-col h-[50vh] md:h-full z-10">
        <div className="p-6 border-b">
          <h2 className="text-2xl font-bold tracking-tight">Request a Ride</h2>
          <p className="text-zinc-500 text-sm mt-1">
            {!currentTrip?.pickup ? "Tap map to set pickup" : 
             !currentTrip?.dropoff ? "Tap map to set destination" : 
             "Select your vehicle"}
          </p>
        </div>

        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {status === 'ACCEPTED' ? (
             <div className="flex flex-col items-center justify-center space-y-6 py-10 animate-in fade-in slide-in-from-bottom-4">
                <div className="bg-black text-white rounded-full p-4">
                  <Car className="h-8 w-8" />
                </div>
                <div className="text-center space-y-1">
                  <h3 className="text-xl font-bold">Driver Arriving</h3>
                  <p className="text-zinc-500">Your {currentTrip?.vehicle} is on the way.</p>
                </div>
                <div className="w-full bg-zinc-50 p-4 rounded-lg flex justify-between items-center">
                   <div className="flex items-center gap-2">
                     <Clock className="h-4 w-4 text-zinc-500" />
                     <span>{currentTrip?.time}</span>
                   </div>
                   <div className="font-bold">${currentTrip?.price}</div>
                </div>
                <Button size="lg" className="w-full gap-2">
                  <CreditCard className="h-4 w-4" /> Pay Now
                </Button>
             </div>
          ) : status === 'WAITING' ? (
            <div className="flex flex-col items-center justify-center py-20 text-center space-y-4 animate-pulse">
                <div className="h-16 w-16 bg-zinc-100 rounded-full flex items-center justify-center">
                  <Car className="h-8 w-8 text-zinc-400" />
                </div>
                <div>
                   <h3 className="font-medium text-lg">Connecting to drivers...</h3>
                   <p className="text-sm text-zinc-500">Please wait while we find you a ride.</p>
                </div>
            </div>
          ) : (
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
                    <div className="font-bold">${v.price}</div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>

        {status !== 'ACCEPTED' && status !== 'WAITING' && (
          <div className="p-4 border-t bg-zinc-50/50">
             <Button 
               className="w-full text-lg h-12" 
               disabled={!selectedVehicle || !currentTrip?.pickup || !currentTrip?.dropoff}
               onClick={handleRequest}
             >
               Confirm {selectedVehicle ? VEHICLES.find(v => v.id === selectedVehicle)?.name : 'Ride'}
             </Button>
          </div>
        )}
      </div>
    </div>
  )
}
