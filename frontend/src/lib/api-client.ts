import axios, { AxiosError, AxiosInstance, AxiosResponse } from 'axios';
import { env } from '@/config/env';
import { ApiErrorInfo, ApiErrorResponse } from '@/types/api.types';

export class ApiError extends Error {
  readonly code: string;
  readonly status: number;
  readonly details?: unknown[];

  constructor(code: string, message: string, status = 500, details?: unknown[]) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.status = status;
    this.details = details;
  }
}

export const apiClient: AxiosInstance = axios.create({
  baseURL: env.apiBaseUrl,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
});

// Response interceptor: Unwrap and normalize backend responses & errors
apiClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: AxiosError<ApiErrorResponse>) => {
    if (error.response) {
      const data = error.response.data;
      if (data && !data.success && data.error) {
        const errorInfo: ApiErrorInfo = data.error;
        return Promise.reject(
          new ApiError(errorInfo.code, errorInfo.message, error.response.status, errorInfo.details),
        );
      }
      return Promise.reject(
        new ApiError(
          'HTTP_ERROR',
          error.message || `Request failed with status ${error.response.status}`,
          error.response.status,
        ),
      );
    }

    if (error.code === 'ECONNABORTED') {
      return Promise.reject(new ApiError('TIMEOUT', 'Request timed out. Please try again.', 408));
    }

    return Promise.reject(
      new ApiError('NETWORK_ERROR', 'Network connection error or server is unreachable.', 0),
    );
  },
);

export default apiClient;
