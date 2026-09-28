import React from 'react';
import { Link } from 'react-router-dom';
import { Badge } from '@/components/ui/Badge';
import { env } from '@/config/env';

export const Header: React.FC = () => {
  return (
    <header className="bg-slate-900 border-b border-slate-800 text-white sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Link to="/" className="flex items-center space-x-3 group">
            <div className="w-8 h-8 rounded-lg bg-sky-500 flex items-center justify-center font-bold text-slate-950 text-sm shadow-md transition-transform group-hover:scale-105">
              Go
            </div>
            <span className="font-semibold text-lg text-slate-100 tracking-tight">
              {env.appTitle}
            </span>
          </Link>
          <Badge
            variant="info"
            className="hidden sm:inline-flex bg-slate-800 text-sky-400 border-slate-700"
          >
            Go 1.23+ / React 18 / PostgreSQL
          </Badge>
        </div>

        <nav className="flex items-center space-x-6 text-sm font-medium">
          <Link to="/" className="text-slate-300 hover:text-white transition-colors">
            System Status
          </Link>
          <a
            href="https://github.com"
            target="_blank"
            rel="noreferrer"
            className="text-slate-400 hover:text-slate-200 transition-colors"
          >
            Docs
          </a>
        </nav>
      </div>
    </header>
  );
};
