import React, { useState } from 'react';
import { Home, Activity, Zap, Settings, AlertTriangle, Wifi, WifiOff, RotateCcw, Stethoscope, Cpu, Monitor, CreditCard, Users, Key, ChevronLeft, ChevronRight } from 'lucide-react';
import { useAuth } from '@/contexts/AuthContext';
import ConfirmModal from './confirm-modal';

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ isOpen, onClose }) => {
  const { connectionState, version } = useAuth();
  const [showResetModal, setShowResetModal] = useState(false);
  const [componentsOpen, setComponentsOpen] = useState(false);
  const [collapsed, setCollapsed] = useState(false);

  const handleReset = () => {
    // TODO: Replace with real reset logic
    console.log('Reset confirmed!');
    setShowResetModal(false);
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Overlay for mobile */}
      <div
        className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
        onClick={onClose}
      />

      {/* Sidebar */}
      <div
        className={`fixed inset-y-0 left-0 z-50 ${collapsed ? 'w-20' : 'w-64'} bg-card border-r border-border flex flex-col h-screen transition-all duration-200`}
      >
        {/* Header: Only close button for mobile */}
        <div className="sticky top-0 z-20 bg-card flex items-center justify-end p-6 border-b border-border lg:hidden">
          <button
            onClick={onClose}
            className="p-2 rounded-lg hover:bg-muted transition-colors"
          >
            <span className="sr-only">Close menu</span>
            <div className="h-5 w-5 text-muted-foreground">×</div>
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 overflow-y-auto p-4">
          <ul className="space-y-2">
            <li>
              <a href="/dashboard" className="flex items-center px-3 py-2 rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <Home className="h-5 w-5" />
                {!collapsed && <span className="ml-3">Dashboard</span>}
              </a>
            </li>
            <li>
              <a href="/sessions" className="flex items-center px-3 py-2 rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <Zap className="h-5 w-5" />
                {!collapsed && <span className="ml-3">Sessions</span>}
              </a>
            </li>
            <li>
              <a href="/diagnostics" className="flex items-center px-3 py-2 rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <Stethoscope className="h-5 w-5" />
                {!collapsed && <span className="ml-3">Diagnostics</span>}
              </a>
            </li>
            {/* Hardware Menu Entry (was Components) */}
            <li>
              <a href="/hardware" className="flex items-center px-3 py-2 rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <Cpu className="h-5 w-5" />
                {!collapsed && <span className="ml-3">Hardware</span>}
              </a>
            </li>
            <li>
              <a href="/users" className="flex items-center px-3 py-2 rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <Users className="h-5 w-5" />
                {!collapsed && <span className="ml-3">Users</span>}
              </a>
            </li>
            <li>
              <a href="/rfid-tags" className="flex items-center px-3 py-2 rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
                <CreditCard className="h-5 w-5" />
                {!collapsed && <span className="ml-3">RFID Tags</span>}
              </a>
            </li>
          </ul>
        </nav>

        {/* Sticky Footer - Always visible */}
        <div className="border-t border-border bg-card px-3 py-2 flex items-center justify-between gap-2 relative">
          {/* Left: Connection Status and Version */}
          {!collapsed && (
            <div className="flex items-center gap-3">
              <div className="flex flex-col items-center">
              <Wifi className={`h-4 w-4 ${connectionState === 'connected' ? 'text-green-500' : 'text-red-500'}`} />
                <span className={`text-[10px] font-medium capitalize mt-1 ${
                connectionState === 'connected' ? 'text-green-500' : 
                connectionState === 'connecting' ? 'text-yellow-500' : 'text-red-500'
              }`}>
                {connectionState}
              </span>
              </div>
              <span className="text-xs font-semibold text-accent-orange ml-2">v{version}</span>
            </div>
          )}
          {/* Right: Footer Action Icons + Collapse Button */}
          <div className="flex items-center gap-2 ml-auto">
            {!collapsed && (
              <a
                href="#"
                onClick={e => { e.preventDefault(); setShowResetModal(true); }}
                className="flex items-center justify-center rounded-md p-2 w-9 h-9 hover:bg-muted transition-colors focus:outline-none focus:ring-2 focus:ring-primary"
                aria-label="Reset"
                title="Reset"
              >
                <RotateCcw className="h-5 w-5 text-accent-orange" />
              </a>
            )}
            <a
              href="/settings"
              className="flex items-center justify-center rounded-md p-2 hover:bg-muted transition-colors focus:outline-none focus:ring-2 focus:ring-primary"
              aria-label="Settings"
              title="Settings"
            >
              <Settings className="h-5 w-5 text-accent-orange" />
            </a>
            <button
              className="flex items-center justify-center w-9 h-9 rounded-md bg-card border border-card transition-colors focus:outline-none focus:ring-2 focus:ring-primary"
              onClick={() => setCollapsed((c) => !c)}
              aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
            >
              {collapsed ? <ChevronRight className="h-5 w-5" /> : <ChevronLeft className="h-5 w-5" />}
            </button>
          </div>
        </div>
      </div>
      {/* Modal for reset confirmation */}
      <ConfirmModal
        open={showResetModal}
        title="Confirm Reset"
        description="Are you sure you want to reset? This action cannot be undone."
        onConfirm={handleReset}
        onCancel={() => setShowResetModal(false)}
      />
    </>
  );
};

export default Sidebar; 