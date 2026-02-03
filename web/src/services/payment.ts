import { apiRequest } from '@/lib/api';
import type {
  CalculateFareRequest,
  CalculateFareResponse,
  ProcessPaymentRequest,
  ProcessPaymentResponse,
} from '@/types';

export const paymentService = {
  async calculateFare(data: CalculateFareRequest): Promise<CalculateFareResponse> {
    return apiRequest('/payment/calculate-fare', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async processPayment(data: ProcessPaymentRequest): Promise<ProcessPaymentResponse> {
    return apiRequest('/payment/pay', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },
};
