import { apiClient } from '@/lib/api-client';
import { ApiResponse } from '@/types/api.types';
import { HealthData, LivenessData, ReadinessData } from '../types/health.types';

export const healthService = {
  async getHealth(): Promise<HealthData> {
    const response = await apiClient.get<ApiResponse<HealthData>>('/health');
    return response.data.data;
  },

  async getLiveness(): Promise<LivenessData> {
    const response = await apiClient.get<ApiResponse<LivenessData>>('/health/live');
    return response.data.data;
  },

  async getReadiness(): Promise<ReadinessData> {
    const response = await apiClient.get<ApiResponse<ReadinessData>>('/health/ready');
    return response.data.data;
  },
};

export default healthService;
