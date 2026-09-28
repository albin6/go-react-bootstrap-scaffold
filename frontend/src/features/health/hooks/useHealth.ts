import { useState, useEffect, useCallback } from 'react';
import { healthService } from '../services/health.service';
import { HealthData } from '../types/health.types';
import { ApiError } from '@/lib/api-client';

export interface UseHealthReturn {
  health: HealthData | null;
  isLoading: boolean;
  error: ApiError | null;
  refetch: () => Promise<void>;
}

export function useHealth(autoFetch = true): UseHealthReturn {
  const [health, setHealth] = useState<HealthData | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(autoFetch);
  const [error, setError] = useState<ApiError | null>(null);

  const fetchHealth = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await healthService.getHealth();
      setHealth(data);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err);
      } else {
        setError(new ApiError('UNKNOWN_ERROR', 'An unexpected error occurred.'));
      }
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    if (autoFetch) {
      void fetchHealth();
    }
  }, [autoFetch, fetchHealth]);

  return {
    health,
    isLoading,
    error,
    refetch: fetchHealth,
  };
}

export default useHealth;
