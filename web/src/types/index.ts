// API Types
export interface User {
  user_id: string;
  email: string;
  role: 'rider' | 'driver';
}

export interface AuthResponse {
  user_id: string;
  token: string;
  role: 'rider' | 'driver';
}

export interface RegisterRequest {
  email: string;
  password: string;
  role: 'rider' | 'driver';
  name: string;
  phone: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

// Rider Types
export interface RiderProfile {
  rider_id: string;
  user_id: string;
  name: string;
  phone: string;
  payment_method?: string;
  total_rides?: number;
  average_rating?: number;
}

export interface CreateRiderProfileRequest {
  name: string;
  phone: string;
  payment_method?: string;
}

export interface RequestRideRequest {
  pickup: {
    lat: number;
    lng: number;
  };
  dropoff: {
    lat: number;
    lng: number;
  };
  vehicle_type: VehicleType;
}

export interface RequestRideResponse {
  trip_id: string;
  status: string;
  estimated_fare: number;
  eta: string;
}

// Driver Types
export interface DriverProfile {
  driver_id: string;
  user_id: string;
  name: string;
  phone: string;
  vehicle_type: VehicleType;
  vehicle_plate: string;
  vehicle_brand: string;
  status: 'offline' | 'available' | 'busy';
  total_trips?: number;
  rating?: number;
}

export interface CreateDriverProfileRequest {
  name: string;
  phone: string;
  vehicle_type: VehicleType;
  vehicle_plate: string;
  vehicle_brand: string;
}

export interface UpdateAvailabilityRequest {
  status: 'offline' | 'available' | 'busy';
  location?: {
    lat: number;
    lng: number;
  };
}

// Trip Types
export interface Trip {
  trip_id: string;
  rider_id: string;
  driver_id?: string;
  status: TripStatus;
  pickup: Location;
  dropoff: Location;
  vehicle_type: VehicleType;
  estimated_fare: number;
  final_fare?: number;
  distance_km?: number;
  duration_sec?: number;
  created_at: string;
  updated_at: string;
}

export interface Location {
  latitude: number;
  longitude: number;
}

export type TripStatus = 
  | 'requested' 
  | 'finding_driver' 
  | 'accepted' 
  | 'in_progress' 
  | 'completed' 
  | 'cancelled';

export type VehicleType = 'sedan' | 'suv' | 'van' | 'luxury';

export interface TripSummary {
  trip_id: string;
  status: TripStatus;
  fare?: number;
  created_at: string;
}

// Payment Types
export interface CalculateFareRequest {
  pickup_lat: number;
  pickup_lng: number;
  dropoff_lat: number;
  dropoff_lng: number;
  vehicle_type: VehicleType;
}

export interface CalculateFareResponse {
  base_fare: number;
  distance_km: number;
  duration_min: number;
  surge_multiplier: number;
  total_fare: number;
  breakdown: string;
}

export interface ProcessPaymentRequest {
  trip_id: string;
  amount: number;
  payment_method: string;
}

export interface ProcessPaymentResponse {
  payment_id: string;
  status: string;
  amount: number;
}

// API Error
export interface ApiError {
  error: string;
}
