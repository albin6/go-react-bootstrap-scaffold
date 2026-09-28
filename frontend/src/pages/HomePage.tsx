import React from 'react';
import { HealthDashboard } from '@/features/health/components/HealthDashboard';

export const HomePage: React.FC = () => {
  return (
    <div className="space-y-12 max-w-4xl mx-auto">
      <div className="text-center space-y-4">
        <h1 className="text-3xl sm:text-4xl font-extrabold text-slate-900 tracking-tight">
          Production-Grade Scaffold
        </h1>
        <p className="text-base sm:text-lg text-slate-600 max-w-2xl mx-auto leading-relaxed">
          An idiomatic foundation powered by <strong className="text-slate-900">Go</strong>,{' '}
          <strong className="text-slate-900">React</strong>,{' '}
          <strong className="text-slate-900">TypeScript</strong>, and{' '}
          <strong className="text-slate-900">PostgreSQL</strong> with native{' '}
          <code className="text-xs bg-slate-100 px-1.5 py-0.5 rounded text-slate-800 font-mono">
            pgx
          </code>{' '}
          and compile-safe{' '}
          <code className="text-xs bg-slate-100 px-1.5 py-0.5 rounded text-slate-800 font-mono">
            SQLC
          </code>
          .
        </p>
      </div>

      <HealthDashboard />

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 pt-6">
        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm">
          <div className="w-9 h-9 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center font-bold mb-3 text-sm">
            1
          </div>
          <h4 className="font-semibold text-slate-900 mb-1">Strict Boundaries</h4>
          <p className="text-xs text-slate-600 leading-relaxed">
            Independent frontend and backend workspaces. Zero shared code, zero cross-dependencies.
          </p>
        </div>

        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm">
          <div className="w-9 h-9 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center font-bold mb-3 text-sm">
            2
          </div>
          <h4 className="font-semibold text-slate-900 mb-1">Idiomatic Go Backend</h4>
          <p className="text-xs text-slate-600 leading-relaxed">
            Standard net/http ServeMux routing, structured slog logging, pgxpool, and clean error
            wrapping.
          </p>
        </div>

        <div className="bg-white p-5 rounded-xl border border-slate-200 shadow-sm">
          <div className="w-9 h-9 rounded-lg bg-purple-50 text-purple-600 flex items-center justify-center font-bold mb-3 text-sm">
            3
          </div>
          <h4 className="font-semibold text-slate-900 mb-1">Compile-Safe Persistence</h4>
          <p className="text-xs text-slate-600 leading-relaxed">
            Explicit PostgreSQL migrations and type-safe query execution generated automatically via
            SQLC.
          </p>
        </div>
      </div>
    </div>
  );
};

export default HomePage;
