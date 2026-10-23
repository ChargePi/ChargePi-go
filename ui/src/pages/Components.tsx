import React, { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useTheme } from '@/contexts/ThemeContext';
import { Menu, X, LogOut, Sun, Moon, Settings, Zap, Wifi, Cpu, HardDrive, Power, AlertTriangle, CheckCircle, XCircle, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import Sidebar from '@/components/Sidebar';
import { formatDate, getUserPrefs } from '@/utils/dateUtils';

interface Component {
  id: string;
  name: string;
  type: 'connector' | 'controller' | 'network' | 'power' | 'storage' | 'sensor';
  status: 'online' | 'offline' | 'error' | 'maintenance';
  health: number; // 0-100
  lastSeen: string;
  version: string;
  description: string;
  location: string;
  powerConsumption?: number;
  temperature?: number;
  uptime?: string;
}

interface ComponentsProps {
  type?: 'evse' | 'display' | 'reader' | 'all';
}

const typeMap: Record<string, string> = {
  evse: 'connector',
  display: 'controller',
  reader: 'sensor',
};


function formatPower(value: number, unit: string) {
  if (unit === 'kW') return `${(value/1000).toFixed(2)} kW`;
  if (unit === 'W') return `${value.toFixed(0)} W`;
  return `${value} W`;
}

const Components: React.FC<ComponentsProps> = ({ type = 'all' }) => {
  const { user, logout, isAuthEnabled } = useAuth();
  const { theme, toggleTheme } = useTheme();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [components, setComponents] = useState<Component[]>([]);
  const [filter, setFilter] = useState<'all' | 'online' | 'offline' | 'error'>('all');
  const [typeFilter, setTypeFilter] = useState<'all' | 'connector' | 'controller' | 'network' | 'power' | 'storage' | 'sensor'>(type === 'all' ? 'all' : (typeMap[type] as 'connector' | 'controller' | 'sensor' | 'all'));
  const [isRefreshing, setIsRefreshing] = useState(false);

  useEffect(() => {
    if (type !== 'all') {
      setTypeFilter(typeMap[type] as 'connector' | 'controller' | 'sensor' | 'all');
    }
  }, [type]);

  // Mock components data
  useEffect(() => {
    const mockComponents: Component[] = [
      {
        id: '1',
        name: 'Connector 1',
        type: 'connector',
        status: 'online',
        health: 95,
        lastSeen: '2024-01-07 16:30',
        version: '2.1.0',
        description: 'Type 2 charging connector',
        location: 'Bay A',
        powerConsumption: 22.4,
        temperature: 45,
        uptime: '15d 8h 32m'
      },
      {
        id: '2',
        name: 'Connector 2',
        type: 'connector',
        status: 'online',
        health: 88,
        lastSeen: '2024-01-07 16:30',
        version: '2.1.0',
        description: 'Type 2 charging connector',
        location: 'Bay B',
        powerConsumption: 0,
        temperature: 42,
        uptime: '15d 8h 32m'
      },
      {
        id: '3',
        name: 'Connector 3',
        type: 'connector',
        status: 'error',
        health: 23,
        lastSeen: '2024-01-07 15:30',
        version: '2.1.0',
        description: 'Type 2 charging connector',
        location: 'Bay C',
        powerConsumption: 0,
        temperature: 65,
        uptime: '15d 8h 32m'
      },
      {
        id: '4',
        name: 'Main Controller',
        type: 'controller',
        status: 'online',
        health: 98,
        lastSeen: '2024-01-07 16:30',
        version: '1.5.2',
        description: 'Central control unit',
        location: 'Control Room',
        powerConsumption: 15.2,
        temperature: 38,
        uptime: '45d 12h 15m'
      },
      {
        id: '5',
        name: 'Network Switch',
        type: 'network',
        status: 'online',
        health: 92,
        lastSeen: '2024-01-07 16:30',
        version: '3.0.1',
        description: 'Ethernet network switch',
        location: 'Control Room',
        powerConsumption: 8.5,
        temperature: 41,
        uptime: '30d 6h 45m'
      },
      {
        id: '6',
        name: 'Power Supply Unit',
        type: 'power',
        status: 'online',
        health: 87,
        lastSeen: '2024-01-07 16:30',
        version: '1.2.0',
        description: 'Main power distribution unit',
        location: 'Electrical Room',
        powerConsumption: 125.8,
        temperature: 52,
        uptime: '60d 3h 20m'
      },
      {
        id: '7',
        name: 'Temperature Sensor 1',
        type: 'sensor',
        status: 'maintenance',
        health: 45,
        lastSeen: '2024-01-07 14:15',
        version: '1.0.3',
        description: 'Temperature monitoring sensor',
        location: 'Bay A',
        temperature: 48,
        uptime: '20d 9h 12m'
      },
      {
        id: '8',
        name: 'Data Storage',
        type: 'storage',
        status: 'online',
        health: 96,
        lastSeen: '2024-01-07 16:30',
        version: '2.0.1',
        description: 'Local data storage system',
        location: 'Control Room',
        powerConsumption: 12.3,
        temperature: 35,
        uptime: '25d 18h 30m'
      }
    ];
    setComponents(mockComponents);
  }, []);

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'connector':
        return <Zap className="h-5 w-5" />;
      case 'controller':
        return <Cpu className="h-5 w-5" />;
      case 'network':
        return <Wifi className="h-5 w-5" />;
      case 'power':
        return <Power className="h-5 w-5" />;
      case 'storage':
        return <HardDrive className="h-5 w-5" />;
      case 'sensor':
        return <Settings className="h-5 w-5" />;
      default:
        return <Settings className="h-5 w-5" />;
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online':
        return <CheckCircle className="h-5 w-5 text-green-500" />;
      case 'offline':
        return <XCircle className="h-5 w-5 text-gray-500" />;
      case 'error':
        return <XCircle className="h-5 w-5 text-red-500" />;
      case 'maintenance':
        return <AlertTriangle className="h-5 w-5 text-yellow-500" />;
      default:
        return <XCircle className="h-5 w-5 text-gray-500" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online':
        return 'bg-green-100 dark:bg-green-900/20 text-green-800 dark:text-green-400';
      case 'offline':
        return 'bg-gray-100 dark:bg-gray-900/20 text-gray-800 dark:text-gray-400';
      case 'error':
        return 'bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-400';
      case 'maintenance':
        return 'bg-yellow-100 dark:bg-yellow-900/20 text-yellow-800 dark:text-yellow-400';
      default:
        return 'bg-gray-100 dark:bg-gray-900/20 text-gray-800 dark:text-gray-400';
    }
  };

  const getHealthColor = (health: number) => {
    if (health >= 80) return 'text-green-500';
    if (health >= 60) return 'text-yellow-500';
    return 'text-red-500';
  };

  const filteredComponents = components.filter(component => {
    const matchesFilter = filter === 'all' || component.status === filter;
    const matchesType = typeFilter === 'all' || component.type === typeFilter;
    return matchesFilter && matchesType;
  });

  const handleLogout = () => {
    logout();
  };

  const refreshComponents = async () => {
    setIsRefreshing(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000));
    setIsRefreshing(false);
  };

  const onlineCount = components.filter(c => c.status === 'online').length;
  const errorCount = components.filter(c => c.status === 'error').length;
  const maintenanceCount = components.filter(c => c.status === 'maintenance').length;

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      
      {/* Sticky Header */}
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">Components</h1>
          </div>
          <div className="flex items-center space-x-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={refreshComponents}
              disabled={isRefreshing}
            >
              <RefreshCw className={`h-5 w-5 ${isRefreshing ? 'animate-spin' : ''}`} />
            </Button>
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
          <div className="bg-card/80 border-2 border-primary/30 shadow-lg rounded-xl p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Total Components</p>
                <p className="text-2xl font-bold">{components.length}</p>
              </div>
              <Settings className="h-8 w-8 text-muted-foreground" />
            </div>
          </div>
          <div className="bg-card/80 border-2 border-primary/30 shadow-lg rounded-xl p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Online</p>
                <p className="text-2xl font-bold text-green-600">{onlineCount}</p>
              </div>
              <CheckCircle className="h-8 w-8 text-green-500" />
            </div>
          </div>
          <div className="bg-card/80 border-2 border-primary/30 shadow-lg rounded-xl p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Errors</p>
                <p className="text-2xl font-bold text-red-600">{errorCount}</p>
              </div>
              <XCircle className="h-8 w-8 text-red-500" />
            </div>
          </div>
          <div className="bg-card/80 border-2 border-primary/30 shadow-lg rounded-xl p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Maintenance</p>
                <p className="text-2xl font-bold text-yellow-600">{maintenanceCount}</p>
              </div>
              <AlertTriangle className="h-8 w-8 text-yellow-500" />
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
              <option value="online">Online</option>
              <option value="offline">Offline</option>
              <option value="error">Error</option>
            </select>
          </div>
          <div className="flex items-center space-x-2">
            <label className="text-sm font-medium">Type:</label>
            <select
              value={typeFilter}
              onChange={(e) => setTypeFilter(e.target.value as any)}
              className="px-3 py-1 border border-border rounded-md bg-background text-sm"
            >
              <option value="all">All</option>
              <option value="connector">Connector</option>
              <option value="controller">Controller</option>
              <option value="network">Network</option>
              <option value="power">Power</option>
              <option value="storage">Storage</option>
              <option value="sensor">Sensor</option>
            </select>
          </div>
        </div>

        {/* Components Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredComponents.map((component) => (
            <div
              key={component.id}
              className="bg-card/80 border-2 border-primary/30 shadow-lg rounded-xl p-4 hover:shadow-xl transition-shadow"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center space-x-2">
                  {getTypeIcon(component.type)}
                  <div>
                    <h3 className="font-semibold text-foreground">{component.name}</h3>
                    <p className="text-xs text-muted-foreground capitalize">{component.type}</p>
                  </div>
                </div>
                {getStatusIcon(component.status)}
              </div>

              <div className="space-y-2 mb-4">
                <p className="text-sm text-muted-foreground">{component.description}</p>
                <div className="flex items-center justify-between text-xs text-muted-foreground">
                  <span>Location: {component.location}</span>
                  <span>v{component.version}</span>
                </div>
              </div>

              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-xs text-muted-foreground">Health</span>
                  <span className={`text-sm font-medium ${getHealthColor(component.health)}`}>
                    {component.health}%
                  </span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div
                    className={`h-2 rounded-full ${getHealthColor(component.health).replace('text-', 'bg-')}`}
                    style={{ width: `${component.health}%` }}
                  />
                </div>
              </div>

              {component.powerConsumption !== undefined && (
                <div className="mt-3 pt-3 border-t border-border">
                  <div className="grid grid-cols-2 gap-2 text-xs">
                    <div>
                      <span className="text-muted-foreground">Power:</span>
                      <span className="ml-1 font-medium">{formatPower(component.powerConsumption, getUserPrefs().powerUnit)}</span>
                    </div>
                    {component.temperature !== undefined && (
                      <div>
                        <span className="text-muted-foreground">Temp:</span>
                        <span className="ml-1 font-medium">{component.temperature}°C</span>
                      </div>
                    )}
                  </div>
                </div>
              )}

              <div className="mt-3 pt-3 border-t border-border">
                <div className="flex items-center justify-between text-xs text-muted-foreground">
                  <span>Last seen: {formatDate(component.lastSeen, getUserPrefs().timezone)}</span>
                  {component.uptime && (
                    <span>Uptime: {component.uptime}</span>
                  )}
                </div>
              </div>

              <div className="mt-3 flex items-center justify-between">
                <span className={`text-xs px-2 py-1 rounded-full ${getStatusColor(component.status)}`}>
                  {component.status}
                </span>
                <Button size="sm" variant="outline">
                  Details
                </Button>
              </div>
            </div>
          ))}
        </div>

        {filteredComponents.length === 0 && (
          <div className="text-center py-12">
            <Settings className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
            <p className="text-muted-foreground">No components found with current filters</p>
          </div>
        )}
      </main>
    </div>
  );
};

export default Components; 