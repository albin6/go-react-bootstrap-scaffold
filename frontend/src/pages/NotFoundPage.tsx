import React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/Button';

export const NotFoundPage: React.FC = () => {
  return (
    <div className="py-20 text-center max-w-md mx-auto space-y-4">
      <h1 className="text-6xl font-extrabold text-slate-300">404</h1>
      <h2 className="text-2xl font-bold text-slate-900">Page Not Found</h2>
      <p className="text-sm text-slate-600">
        The requested URL was not found on this application route.
      </p>
      <div className="pt-4">
        <Link to="/">
          <Button variant="primary">Return Home</Button>
        </Link>
      </div>
    </div>
  );
};

export default NotFoundPage;
