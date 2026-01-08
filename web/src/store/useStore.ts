import { create } from 'zustand'

export type RideStatus = 'IDLE' | 'REQUESTING' | 'WAITING' | 'ACCEPTED' | 'COMPLETED'
export type VehicleType = 'sedan' | 'suv' | 'van' | 'luxury'

interface Trip {
  id: string
  pickup: { lat: number; lng: number } | null
  dropoff: { lat: number; lng: number } | null
  vehicle: VehicleType
  price: number
  distance: string
  time: string
}

interface TaxxyState {
  // Rider State
  status: RideStatus
  currentTrip: Trip | null
  
  // Actions
  requestRide: (trip: Trip) => void
  acceptRide: () => void
  resetRide: () => void
  setPickup: (lat: number, lng: number) => void
  setDropoff: (lat: number, lng: number) => void
}

export const useStore = create<TaxxyState>((set) => ({
  status: 'IDLE',
  currentTrip: null,

  requestRide: (trip) => set({ status: 'WAITING', currentTrip: trip }),
  acceptRide: () => set({ status: 'ACCEPTED' }), // Driver accepts
  resetRide: () => set({ status: 'IDLE', currentTrip: null }),
  
  setPickup: (lat, lng) => set((state) => ({ 
    currentTrip: { ...state.currentTrip, pickup: { lat, lng } } as Trip 
  })),
  setDropoff: (lat, lng) => set((state) => ({ 
    currentTrip: { ...state.currentTrip, dropoff: { lat, lng } } as Trip 
  })),
}))
