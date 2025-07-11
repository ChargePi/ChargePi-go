import React, { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { LogOut, AlertTriangle, AlertCircle, XCircle, CheckCircle, Info } from 'lucide-react';
import { Button } from '@/components/ui/button';
import Sidebar from '@/components/Sidebar';
import { formatDate, getUserPrefs } from '@/utils/dateUtils';

interface DiagnosticItem {
  id: string;
  type: 'error' | 'warning' | 'info' | 'success';
  title: string;
  description: string;
  component: string;
  timestamp: string;
  severity: 'low' | 'medium' | 'high';
  status: 'active' | 'resolved' | 'acknowledged';
}



const Diagnostics: React.FC = () => {
  const { user, logout, isAuthEnabled } = useAuth();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [diagnostics, setDiagnostics] = useState<DiagnosticItem[]>([]);
  const [filter, setFilter] = useState<'all' | 'active' | 'resolved'>('all');
  const [severityFilter, setSeverityFilter] = useState<'all' | 'high' | 'medium' | 'low'>('all');

  // Mock diagnostics data
  useEffect(() => {
    const mockDiagnostics: DiagnosticItem[] = [
      {
        id: '1',
        type: 'error',
        title: 'Connector 3 Communication Failure',
        description: 'Unable to establish communication with connector 3. Check physical connections and power supply.',
        component: 'Connector 3',
        timestamp: '2024-01-07 15:30',
        severity: 'high',
        status: 'active'
      },
      {
        id: '2',
        type: 'warning',
        title: 'Temperature Sensor Reading High',
        description: 'Temperature sensor on connector 1 is reading above normal operating range.',
        component: 'Connector 1',
        timestamp: '2024-01-07 14:15',
        severity: 'medium',
        status: 'active'
      },
      {
        id: '3',
        type: 'info',
        title: 'System Maintenance Due',
        description: 'Routine maintenance is scheduled for next week.',
        component: 'System',
        timestamp: '2024-01-07 12:00',
        severity: 'low',
        status: 'acknowledged'
      },
      {
        id: '4',
        type: 'success',
        title: 'Connector 2 Self-Test Completed',
        description: 'Connector 2 has successfully completed its self-diagnostic test.',
        component: 'Connector 2',
        timestamp: '2024-01-07 11:45',
        severity: 'low',
        status: 'resolved'
      }
    ];
    setDiagnostics(mockDiagnostics);
  }, []);

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'error':
        return <XCircle className="h-5 w-5 text-red-500" />;
      case 'warning':
        return <AlertTriangle className="h-5 w-5 text-yellow-500" />;
      case 'info':
        return <Info className="h-5 w-5 text-blue-500" />;
      case 'success':
        return <CheckCircle className="h-5 w-5 text-green-500" />;
      default:
        return <AlertCircle className="h-5 w-5 text-gray-500" />;
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'high':
        return 'bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-400 border-red-200 dark:border-red-800';
      case 'medium':
        return 'bg-yellow-100 dark:bg-yellow-900/20 text-yellow-800 dark:text-yellow-400 border-yellow-200 dark:border-yellow-800';
      case 'low':
        return 'bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-400 border-blue-200 dark:border-blue-800';
      default:
        return 'bg-gray-100 dark:bg-gray-900/20 text-gray-800 dark:text-gray-400 border-gray-200 dark:border-gray-800';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-400';
      case 'resolved':
        return 'bg-green-100 dark:bg-green-900/20 text-green-800 dark:text-green-400';
      case 'acknowledged':
        return 'bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-400';
      default:
        return 'bg-gray-100 dark:bg-gray-900/20 text-gray-800 dark:text-gray-400';
    }
  };

  const filteredDiagnostics = diagnostics.filter(item => {
    const matchesFilter = filter === 'all' || item.status === filter;
    const matchesSeverity = severityFilter === 'all' || item.severity === severityFilter;
    return matchesFilter && matchesSeverity;
  });

  const handleLogout = () => {
    logout();
  };

  const acknowledgeDiagnostic = (id: string) => {
    setDiagnostics(prev => prev.map(item => 
      item.id === id ? { ...item, status: 'acknowledged' as const } : item
    ));
  };

  const resolveDiagnostic = (id: string) => {
    setDiagnostics(prev => prev.map(item => 
      item.id === id ? { ...item, status: 'resolved' as const } : item
    ));
  };

  const activeCount = diagnostics.filter(d => d.status === 'active').length;
  const highPriorityCount = diagnostics.filter(d => d.severity === 'high' && d.status === 'active').length;

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      
      {/* Sticky Header */}
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">Diagnostics</h1>
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
      <main className="flex-1 p-6 min-w-0">
        {/* Summary Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-card border border-border rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Diagnostics</p>
                <p className="text-2xl font-bold">{diagnostics.length}</p>
              </div>
              <AlertCircle className="h-8 w-8 text-muted-foreground" />
            </div>
          </div>
          <div className="bg-card border border-border rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Active Issues</p>
                <p className="text-2xl font-bold text-red-600">{activeCount}</p>
              </div>
              <AlertTriangle className="h-8 w-8 text-red-500" />
            </div>
          </div>
          <div className="bg-card border border-border rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">High Priority</p>
                <p className="text-2xl font-bold text-red-600">{highPriorityCount}</p>
              </div>
              <XCircle className="h-8 w-8 text-red-500" />
            </div>
          </div>
          <div className="bg-card border border-border rounded-lg p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Resolved</p>
                <p className="text-2xl font-bold text-green-600">
                  {diagnostics.filter(d => d.status === 'resolved').length}
                </p>
              </div>
              <CheckCircle className="h-8 w-8 text-green-500" />
            </div>
          </div>
        </div>

        {/* Filters */}
        <div className="flex flex-wrap gap-4 mb-6">
          <div className="flex items-center space-x-2">
            <label className="text-sm font-medium">Status:</label>
            <select
              value={filter}
              onChange={(e) => setFilter(e.target.value as any)}
              className="px-3 py-1 border border-border rounded-md bg-background text-sm"
            >
              <option value="all">All</option>
              <option value="active">Active</option>
              <option value="resolved">Resolved</option>
            </select>
          </div>
          <div className="flex items-center space-x-2">
            <label className="text-sm font-medium">Severity:</label>
            <select
              value={severityFilter}
              onChange={(e) => setSeverityFilter(e.target.value as any)}
              className="px-3 py-1 border border-border rounded-md bg-background text-sm"
            >
              <option value="all">All</option>
              <option value="high">High</option>
              <option value="medium">Medium</option>
              <option value="low">Low</option>
            </select>
          </div>
        </div>

        {/* Diagnostics List */}
        <div className="space-y-4">
          {filteredDiagnostics.length === 0 ? (
            <div className="text-center py-12">
              <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-4" />
              <p className="text-muted-foreground">No diagnostics found with current filters</p>
            </div>
          ) : (
            filteredDiagnostics.map((item) => (
              <div
                key={item.id}
                className="bg-card border border-border rounded-lg p-4 hover:shadow-md transition-shadow"
              >
                <div className="flex items-start space-x-4">
                  {getTypeIcon(item.type)}
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-2">
                      <h3 className="font-semibold text-foreground">{item.title}</h3>
                      <div className="flex items-center space-x-2">
                        <span className={`text-xs px-2 py-1 rounded-full border ${getSeverityColor(item.severity)}`}>
                          {item.severity}
                        </span>
                        <span className={`text-xs px-2 py-1 rounded-full ${getStatusColor(item.status)}`}>
                          {item.status}
                        </span>
                      </div>
                    </div>
                    <p className="text-sm text-muted-foreground mb-2">{item.description}</p>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-4 text-xs text-muted-foreground">
                        <span>Component: {item.component}</span>
                        <span>Time: {formatDate(item.timestamp, getUserPrefs().timezone)}</span>
                      </div>
                      {/* Action buttons removed as requested */}
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
};

export default Diagnostics; 