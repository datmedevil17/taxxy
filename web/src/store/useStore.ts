import { create } from 'zustand';

// Simple store for global UI state
// All API data is handled in page components with local state

interface AppState {
  // Global loading state
  globalLoading: boolean;
  setGlobalLoading: (loading: boolean) => void;
  
  // Toast/notification state
  notification: { message: string; type: 'success' | 'error' | 'info' } | null;
  showNotification: (message: string, type: 'success' | 'error' | 'info') => void;
  clearNotification: () => void;
}

export const useStore = create<AppState>((set) => ({
  globalLoading: false,
  setGlobalLoading: (loading) => set({ globalLoading: loading }),
  
  notification: null,
  showNotification: (message, type) => set({ notification: { message, type } }),
  clearNotification: () => set({ notification: null }),
}));

// Re-export types for backward compatibility
export type VehicleType = 'sedan' | 'suv' | 'van' | 'luxury';
export type RideStatus = 'idle' | 'requesting' | 'waiting' | 'accepted' | 'in_progress' | 'completed';
