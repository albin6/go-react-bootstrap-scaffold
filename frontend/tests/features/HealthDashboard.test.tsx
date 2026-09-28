import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { HealthDashboard } from '@/features/health/components/HealthDashboard';
import { healthService } from '@/features/health/services/health.service';
import { HealthData } from '@/features/health/types/health.types';

describe('HealthDashboard Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state initially', () => {
    vi.spyOn(healthService, 'getHealth').mockReturnValue(new Promise(() => {}));
    render(<HealthDashboard />);
    expect(screen.getByText(/querying backend \/api\/v1\/health probe/i)).toBeInTheDocument();
  });

  it('renders healthy status and telemetry metrics upon success', async () => {
    const mockHealth: HealthData = {
      status: 'healthy',
      uptime: '2h 15m',
      timestamp: '2026-09-28T20:00:00Z',
      checks: {
        process: 'alive',
        database: 'up',
      },
    };

    vi.spyOn(healthService, 'getHealth').mockResolvedValue(mockHealth);

    render(<HealthDashboard />);

    await waitFor(() => {
      expect(screen.getByText('HEALTHY')).toBeInTheDocument();
    });

    expect(screen.getByText('Connected (pgxpool)')).toBeInTheDocument();
    expect(screen.getByText('2h 15m')).toBeInTheDocument();
  });
});
