import { useState } from "react"
import { useStore, type VehicleType } from "@/store/useStore"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { ArrowLeft, CheckCircle2, MapPin, User, XCircle } from "lucide-react"
import { Link } from "react-router-dom"
import { cn } from "@/lib/utils"

export default function Driver() {
  const { currentTrip, status, acceptRide, resetRide } = useStore()
  const [myVehicle, setMyVehicle] = useState<VehicleType>('sedan')

  const isMatchingVehicle = currentTrip?.vehicle === myVehicle

  return (
    <div className="min-h-screen bg-zinc-950 text-white p-4 font-sans">
      <div className="max-w-md mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center gap-4 py-4">
          <Link to="/">
            <Button variant="ghost" size="icon" className="text-white hover:text-white/80 hover:bg-white/10">
              <ArrowLeft className="h-6 w-6" />
            </Button>
          </Link>
          <h1 className="text-xl font-bold">Driver Portal</h1>
        </div>

        {/* Vehicle Selection */}
        <Card className="bg-zinc-900 border-zinc-800 text-white">
          <CardContent className="p-4">
             <label className="text-sm text-zinc-400 mb-2 block">Your Vehicle Type</label>
             <div className="grid grid-cols-4 gap-2">
               {['sedan', 'suv', 'van', 'luxury'].map((type) => (
                 <button
                   key={type}
                   onClick={() => setMyVehicle(type as VehicleType)}
                   className={cn(
                     "p-2 rounded text-xs font-medium capitalize transition-colors",
                     myVehicle === type 
                       ? "bg-white text-black" 
                       : "bg-zinc-800 text-zinc-400 hover:bg-zinc-700"
                   )}
                 >
                   {type}
                 </button>
               ))}
             </div>
          </CardContent>
        </Card>

        {/* Status Area */}
        <div className="space-y-4">
           <h2 className="font-semibold text-lg flex items-center gap-2">
             Requests
             {status === 'WAITING' && isMatchingVehicle && (
               <span className="flex h-2 w-2 rounded-full bg-green-500 animate-pulse" />
             )}
           </h2>

           {status === 'WAITING' && isMatchingVehicle ? (
             <Card className="bg-white text-black border-none animate-in slide-in-from-bottom-2">
               <CardContent className="p-6 space-y-4">
                  <div className="flex justify-between items-start">
                    <div className="flex items-center gap-3">
                       <div className="h-12 w-12 bg-zinc-100 rounded-full flex items-center justify-center">
                         <User className="h-6 w-6 text-zinc-600" />
                       </div>
                       <div>
                         <p className="font-bold text-lg">New Trip Request</p>
                         <p className="text-sm text-zinc-500">{currentTrip?.distance} • {currentTrip?.time}</p>
                       </div>
                    </div>
                    <div className="font-bold text-xl">${currentTrip?.price}</div>
                  </div>
                  
                  <div className="space-y-4 py-2">
                     <div className="flex items-center gap-3 text-sm">
                        <MapPin className="h-4 w-4 text-green-600 shrink-0" />
                        <span className="truncate">Pickup Location (Lat: {currentTrip?.pickup?.lat.toFixed(4)})</span>
                     </div>
                     <div className="pl-2 ml-[5px] border-l-2 border-zinc-100 h-4" />
                     <div className="flex items-center gap-3 text-sm">
                        <MapPin className="h-4 w-4 text-red-600 shrink-0" />
                        <span className="truncate">Dropoff Location (Lat: {currentTrip?.dropoff?.lat.toFixed(4)})</span>
                     </div>
                  </div>

                  <div className="grid grid-cols-2 gap-3 pt-2">
                     <Button 
                       variant="outline" 
                       onClick={resetRide}
                       className="border-zinc-200 hover:bg-zinc-50 text-black gap-2"
                     >
                       <XCircle className="h-4 w-4" /> Reject
                     </Button>
                     <Button 
                       onClick={acceptRide}
                       className="bg-black text-white hover:bg-zinc-800 gap-2"
                     >
                        <CheckCircle2 className="h-4 w-4" /> Accept
                     </Button>
                  </div>
               </CardContent>
             </Card>
           ) : status === 'ACCEPTED' && isMatchingVehicle ? (
              <Card className="bg-green-500 text-white border-none">
                <CardContent className="p-6 flex flex-col items-center justify-center space-y-2">
                   <CheckCircle2 className="h-12 w-12" />
                   <h3 className="font-bold text-xl">Trip Accepted!</h3>
                   <p className="text-white/80 text-sm">Head to the pickup location.</p>
                </CardContent>
              </Card>
           ) : (
             <div className="text-center py-12 text-zinc-500">
                <p>No requests nearby.</p>
                <p className="text-xs mt-1 text-zinc-600">Requests for "{myVehicle}" will appear here.</p>
             </div>
           )}
        </div>
      </div>
    </div>
  )
}
