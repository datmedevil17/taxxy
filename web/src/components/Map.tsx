import { MapContainer, TileLayer, Marker, Polyline, useMapEvents } from 'react-leaflet'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import { cn } from '@/lib/utils'

// Fix for default marker icons in React Leaflet
delete (L.Icon.Default.prototype as any)._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
})

interface MapProps {
  className?: string
  pickup?: { lat: number; lng: number } | null
  dropoff?: { lat: number; lng: number } | null
  onMapClick?: (lat: number, lng: number) => void
}

function LocationMarker({ onMapClick }: { onMapClick?: (lat: number, lng: number) => void }) {
  useMapEvents({
    click(e) {
      onMapClick?.(e.latlng.lat, e.latlng.lng)
    },
  })
  return null
}

export default function Map({ className, pickup, dropoff, onMapClick }: MapProps) {
  // Mumbai coordinates
  const center = { lat: 19.0760, lng: 72.8777 }

  return (
    <div className={cn("relative h-full w-full overflow-hidden rounded-lg bg-zinc-100", className)}>
      <MapContainer 
        center={center} 
        zoom={13} 
        scrollWheelZoom={true} 
        className="h-full w-full z-0 grayscale"
        // Force grayscale via CSS on the tiles
        style={{ filter: 'grayscale(100%)' }}
      >
         <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
        />
        <LocationMarker onMapClick={onMapClick} />
        {pickup && <Marker position={pickup} />}
        {dropoff && <Marker position={dropoff} />}
        {pickup && dropoff && (
          <Polyline positions={[pickup, dropoff]} color="black" dashArray="5, 10" />
        )}
      </MapContainer>
    </div>
  )
}
