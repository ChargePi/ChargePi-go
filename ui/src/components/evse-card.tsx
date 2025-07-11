import React from 'react';
import { Settings, Power, Activity, MapPin } from 'lucide-react';
import { Button } from './ui/button';

interface EVSE {
  id: string;
  name: string;
  status: 'online' | 'offline' | 'charging' | 'available';
  location: string;
  power?: string;
  currentSession?: string;
}

interface EVSECardProps {
  evse: EVSE;
  onClick: (evse: EVSE) => void;
}

const EVSECard: React.FC<EVSECardProps> = ({ evse, onClick }) => {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'online':
        return 'text-green-600 bg-green-100 dark:bg-green-900/20';
      case 'offline':
        return 'text-red-600 bg-red-100 dark:bg-red-900/20';
      case 'charging':
        return 'text-blue-600 bg-blue-100 dark:bg-blue-900/20';
      case 'available':
        return 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900/20';
      default:
        return 'text-gray-600 bg-gray-100 dark:bg-gray-900/20';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online':
        return <Activity className="h-4 w-4" />;
      case 'offline':
        return <Power className="h-4 w-4" />;
      case 'charging':
        return <Activity className="h-4 w-4" />;
      case 'available':
        return <Power className="h-4 w-4" />;
      default:
        return <Power className="h-4 w-4" />;
    }
  };

  return (
    <div 
      className="bg-card border border-border rounded-lg p-6 flex flex-col justify-between h-full hover:shadow-lg transition-all duration-200 cursor-pointer hover:border-accent-orange/50"
      onClick={() => onClick(evse)}
    >
      {/* Header with Icon and Status */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-accent-orange/10 rounded-lg flex items-center justify-center">
            <Power className="h-6 w-6 text-accent-orange" />
          </div>
          <div className="flex flex-col justify-center">
            <h3 className="font-semibold text-lg leading-tight">{evse.name}</h3>
            <p className="text-sm text-muted-foreground leading-tight">ID: {evse.id}</p>
          </div>
        </div>
        <div className={`flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(evse.status)}`}>
          {getStatusIcon(evse.status)}
          <span className="capitalize">{evse.status}</span>
        </div>
      </div>

      {/* Info Section */}
      <div className="flex flex-col gap-2 mb-6 text-left">
        <div className="flex items-center gap-2 text-sm">
          <MapPin className="h-4 w-4 text-muted-foreground" />
          <span className="text-muted-foreground">Location:</span>
          <span className="font-medium">{evse.location}</span>
        </div>
        {evse.power && (
          <div className="flex items-center gap-2 text-sm">
            <Power className="h-4 w-4 text-muted-foreground" />
            <span className="text-muted-foreground">Power:</span>
            <span className="font-medium">{evse.power}</span>
          </div>
        )}
        {evse.currentSession && (
          <div className="flex items-center gap-2 text-sm">
            <Activity className="h-4 w-4 text-muted-foreground" />
            <span className="text-muted-foreground">Session:</span>
            <span className="font-medium">{evse.currentSession}</span>
          </div>
        )}
      </div>

      {/* Actions Section */}
      <div className="border-t border-border pt-4 mt-auto">
        <div className="flex items-center justify-between w-full">
          <Button 
            variant="outline" 
            size="sm"
            className="flex items-center gap-2"
            onClick={(e) => {
              e.stopPropagation();
              // Handle settings action
            }}
          >
            <Settings className="h-4 w-4" />
            Settings
          </Button>
          <div className="flex items-center gap-2">
            <Button 
              variant="ghost" 
              size="sm"
              className="text-xs"
              onClick={(e) => {
                e.stopPropagation();
                // Handle status action
              }}
            >
              Status
            </Button>
            <Button 
              variant="ghost" 
              size="sm"
              className="text-xs"
              onClick={(e) => {
                e.stopPropagation();
                // Handle logs action
              }}
            >
              Logs
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default EVSECard; 