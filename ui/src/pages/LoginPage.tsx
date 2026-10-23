import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { useTheme } from '@/contexts/ThemeContext';
import { Sun, Moon, Zap, Wifi, WifiOff } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

const LoginPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login, isLoading, connectionState, version, isAuthEnabled } = useAuth();
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();

  // If auth is disabled, redirect to dashboard
  useEffect(() => {
    if (!isAuthEnabled) {
      navigate('/dashboard');
    }
  }, [isAuthEnabled, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!username || !password) {
      setError('Please enter both username and password');
      return;
    }

    const success = await login(username, password);
    if (success) {
      navigate('/dashboard');
    } else {
      setError('Invalid credentials. Please try again.');
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-background">
      {/* Header with theme toggle */}
      <header className="flex justify-between items-center p-6 border-b">
        <div className="flex items-center space-x-2">
          <Zap className="h-8 w-8 text-accent-orange" />
          <h1 className="text-2xl font-bold">ChargePi</h1>
        </div>
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleTheme}
        >
          {theme === 'dark' ? (
            <Sun className="h-5 w-5" />
          ) : (
            <Moon className="h-5 w-5" />
          )}
        </Button>
      </header>

      {/* Hero Section */}
      <main className="flex-1 flex items-center justify-center px-6">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <CardTitle className="text-3xl">Welcome Back</CardTitle>
            <p className="text-muted-foreground">
              Sign in to manage your charge point
            </p>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="space-y-2">
                <label htmlFor="username" className="text-sm font-medium">
                  Username
                </label>
                <Input
                  id="username"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="Enter your username"
                  disabled={isLoading}
                />
              </div>

              <div className="space-y-2">
                <label htmlFor="password" className="text-sm font-medium">
                  Password
                </label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Enter your password"
                  disabled={isLoading}
                />
              </div>

              {error && (
                <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-3">
                  <p className="text-destructive text-sm">{error}</p>
                </div>
              )}

              <Button
                type="submit"
                variant="accent"
                className="w-full"
                disabled={isLoading}
              >
                {isLoading ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    Signing in...
                  </>
                ) : (
                  'Sign In'
                )}
              </Button>
            </form>
          </CardContent>
        </Card>
      </main>

      {/* Footer */}
      <footer className="p-6 border-t">
        <div className="flex justify-between items-center text-sm text-muted-foreground">
          <div className="flex items-center space-x-2">
            {connectionState === 'connected' ? (
              <Wifi className="h-4 w-4 text-green-500" />
            ) : (
              <WifiOff className="h-4 w-4 text-red-500" />
            )}
            <span className="capitalize">{connectionState}</span>
          </div>
          <span>v{version}</span>
        </div>
      </footer>
    </div>
  );
};

export default LoginPage; 