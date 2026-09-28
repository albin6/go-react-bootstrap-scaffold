import React from 'react';

export const Footer: React.FC = () => {
  return (
    <footer className="bg-slate-900 border-t border-slate-800 text-slate-400 py-8 text-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          <p>© {new Date().getFullYear()} Production Scaffold. Free and open architecture.</p>
        </div>
        <div className="flex items-center space-x-6 text-xs text-slate-500">
          <span>Go + net/http</span>
          <span>•</span>
          <span>PostgreSQL + pgx + SQLC</span>
          <span>•</span>
          <span>React + Vite + TypeScript</span>
        </div>
      </div>
    </footer>
  );
};
