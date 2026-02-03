import { apiRequest } from '@/lib/api';
import type { Trip } from '@/types';

export const tripService = {
  async getTrip(tripId: string): Promise<Trip> {
    return apiRequest(`/trip/${tripId}`);
  },
};
