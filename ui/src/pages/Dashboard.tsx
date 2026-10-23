import React, { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { LogOut } from 'lucide-react';
import { Button } from '@/components/ui/button';
import Sidebar from '@/components/Sidebar';

const SIDEBAR_WIDTH = 256; // 64 * 4 (w-64)

const Dashboard: React.FC = () => {
  const { user, logout, isAuthEnabled } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const handleLogout = () => {
    logout();
  };

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      {/* Sticky Header */}
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">Dashboard</h1>
          </div>
          <div className="flex items-center space-x-4">
            {isAuthEnabled && (
              <Button
                variant="ghost"
                size="icon"
                onClick={handleLogout}
              >
                <LogOut className="h-5 w-5" />
              </Button>
            )}
          </div>
        </div>
      </header>
      {/* Main Content */}
      <main className="flex-1 p-6 min-w-0 overflow-hidden">
        {/* Dashboard content goes here */}
      </main>
    </div>
  );
};

export default Dashboard; 