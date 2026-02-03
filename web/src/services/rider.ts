import { apiRequest } from '@/lib/api';
import type {
  RiderProfile,
  CreateRiderProfileRequest,
  RequestRideRequest,
  RequestRideResponse,
  TripSummary,
} from '@/types';

export const riderService = {
  async createProfile(data: CreateRiderProfileRequest): Promise<{ rider_id: string; status: string }> {
    return apiRequest('/rider/profile', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async getProfile(): Promise<RiderProfile> {
    return apiRequest('/rider/profile');
  },

  async requestRide(data: RequestRideRequest): Promise<RequestRideResponse> {
    return apiRequest('/rider/request-ride', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async getTrips(): Promise<{ trips: TripSummary[]; total: number }> {
    return apiRequest('/rider/trips');
  },
};
