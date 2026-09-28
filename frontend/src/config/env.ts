// Frontend Environment Variables Configuration

export const env = {
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  appTitle: import.meta.env.VITE_APP_TITLE || 'Production Scaffold',
  isDev: import.meta.env.DEV,
  isProd: import.meta.env.PROD,
} as const;

export default env;
