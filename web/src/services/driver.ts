import { apiRequest } from '@/lib/api';
import type {
  DriverProfile,
  CreateDriverProfileRequest,
  UpdateAvailabilityRequest,
} from '@/types';

export const driverService = {
  async createProfile(data: CreateDriverProfileRequest): Promise<{ driver_id: string; status: string }> {
    return apiRequest('/driver/profile', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async getProfile(): Promise<DriverProfile> {
    return apiRequest('/driver/profile');
  },

  async updateAvailability(data: UpdateAvailabilityRequest): Promise<{ status: string }> {
    return apiRequest('/driver/availability', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  async acceptRide(tripId: string): Promise<{ trip_id: string; status: string }> {
    return apiRequest('/driver/accept-ride', {
      method: 'POST',
      body: JSON.stringify({ trip_id: tripId }),
    });
  },

  async startRide(tripId: string): Promise<{ trip_id: string; status: string }> {
    return apiRequest('/driver/start-ride', {
      method: 'POST',
      body: JSON.stringify({ trip_id: tripId }),
    });
  },

  async completeRide(tripId: string): Promise<{ trip_id: string; status: string; fare: number }> {
    return apiRequest('/driver/complete-ride', {
      method: 'POST',
      body: JSON.stringify({ trip_id: tripId }),
    });
  },
};
