import React from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card';
import { Badge } from '@/components/ui/Badge';
import { Button } from '@/components/ui/Button';
import { Spinner } from '@/components/ui/Spinner';
import { useHealth } from '../hooks/useHealth';
import { formatDate } from '@/utils/format';

export const HealthDashboard: React.FC = () => {
  const { health, isLoading, error, refetch } = useHealth();

  return (
    <Card className="max-w-2xl mx-auto border-slate-200">
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle>System Connectivity Probe</CardTitle>
          <p className="text-xs text-slate-500 mt-1">
            Real-time health telemetry from Go API & PostgreSQL
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={() => void refetch()} isLoading={isLoading}>
          Check Status
        </Button>
      </CardHeader>

      <CardContent>
        {isLoading && !health ? (
          <div className="py-12 flex flex-col items-center justify-center space-y-3">
            <Spinner size="lg" />
            <p className="text-sm text-slate-500">Querying backend /api/v1/health probe...</p>
          </div>
        ) : error ? (
          <div className="py-6 px-4 bg-rose-50 border border-rose-200 rounded-lg text-rose-800">
            <div className="flex items-center space-x-2 font-semibold mb-1">
              <svg className="w-5 h-5 text-rose-600" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clipRule="evenodd"
                />
              </svg>
              <span>Backend Connectivity Failed ({error.code})</span>
            </div>
            <p className="text-sm">{error.message}</p>
            <p className="text-xs text-rose-600 mt-2">
              Ensure the backend is running at{' '}
              <code className="bg-rose-100 px-1 py-0.5 rounded">http://localhost:8080</code>.
            </p>
          </div>
        ) : health ? (
          <div className="space-y-6">
            <div className="flex items-center justify-between pb-4 border-b border-slate-100">
              <span className="text-sm font-medium text-slate-700">Overall System Status</span>
              <Badge
                variant={
                  health.status === 'healthy'
                    ? 'success'
                    : health.status === 'degraded'
                      ? 'warning'
                      : 'error'
                }
              >
                {health.status.toUpperCase()}
              </Badge>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="p-3.5 bg-slate-50 rounded-lg border border-slate-200">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider block mb-1">
                  Go API Server
                </span>
                <div className="flex items-center space-x-2">
                  <span className="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
                  <span className="text-sm font-medium text-slate-800">
                    {health.checks.process || 'alive'}
                  </span>
                </div>
              </div>

              <div className="p-3.5 bg-slate-50 rounded-lg border border-slate-200">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider block mb-1">
                  PostgreSQL Pool
                </span>
                <div className="flex items-center space-x-2">
                  <span
                    className={`w-2.5 h-2.5 rounded-full ${
                      health.checks.database === 'up' ? 'bg-emerald-500' : 'bg-rose-500'
                    }`}
                  ></span>
                  <span className="text-sm font-medium text-slate-800">
                    {health.checks.database === 'up' ? 'Connected (pgxpool)' : 'Disconnected'}
                  </span>
                </div>
              </div>

              <div className="p-3.5 bg-slate-50 rounded-lg border border-slate-200">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider block mb-1">
                  Server Uptime
                </span>
                <span className="text-sm font-medium text-slate-800">{health.uptime}</span>
              </div>

              <div className="p-3.5 bg-slate-50 rounded-lg border border-slate-200">
                <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider block mb-1">
                  Last Telemetry Ping
                </span>
                <span className="text-sm font-medium text-slate-800">
                  {formatDate(health.timestamp)}
                </span>
              </div>
            </div>
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
};

export default HealthDashboard;
