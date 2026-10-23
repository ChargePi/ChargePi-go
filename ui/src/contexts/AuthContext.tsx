import React, { createContext, useContext, useState, useEffect } from 'react';
import { api } from '@/services/api';
import { authConfig } from '@/config/auth';

interface User {
  id: string;
  username: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
  isLoading: boolean;
  connectionState: 'connected' | 'disconnected' | 'connecting';
  version: string;
  isAuthEnabled: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: React.ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [connectionState, setConnectionState] = useState<'connected' | 'disconnected' | 'connecting'>('disconnected');
  const [version] = useState('1.0.0');
  const [isAuthEnabled] = useState(authConfig.enabled);

  useEffect(() => {
    // If auth is disabled, automatically set user
    if (!isAuthEnabled) {
      setUser(authConfig.mockUser);
      setConnectionState('connected');
      return;
    }

    // Check for existing session
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      try {
        setUser(JSON.parse(savedUser));
        setConnectionState('connected');
      } catch (error) {
        localStorage.removeItem('user');
      }
    }
  }, [isAuthEnabled]);

  const login = async (username: string, password: string): Promise<boolean> => {
    setIsLoading(true);
    setConnectionState('connecting');
    
    try {
      const response = await api.login(username, password);
      setUser(response.user);
      localStorage.setItem('user', JSON.stringify(response.user));
      setConnectionState('connected');
      return true;
    } catch (error) {
      setConnectionState('disconnected');
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('user');
    setConnectionState('disconnected');
  };

  return (
    <AuthContext.Provider value={{ 
      user, 
      login, 
      logout, 
      isLoading, 
      connectionState, 
      version,
      isAuthEnabled
    }}>
      {children}
    </AuthContext.Provider>
  );
}; 