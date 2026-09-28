export interface LivenessData {
  status: string;
}

export interface ReadinessData {
  status: string;
  database: string;
  details?: Record<string, string>;
}

export interface HealthData {
  status: 'healthy' | 'degraded' | 'unhealthy';
  uptime: string;
  timestamp: string;
  checks: {
    process: string;
    database: string;
    [key: string]: string;
  };
}
