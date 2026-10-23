import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const { user, isAuthEnabled } = useAuth();

  // If auth is disabled, allow access
  if (!isAuthEnabled) {
    return <>{children}</>;
  }

  // If auth is enabled but no user, redirect to login
  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute; 