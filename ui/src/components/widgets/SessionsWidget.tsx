import React from 'react';
import { X, Clock, Zap, User } from 'lucide-react';

interface SessionsWidgetProps {
  title: string;
  onRemove: () => void;
}

const SessionsWidget: React.FC<SessionsWidgetProps> = ({ title, onRemove }) => {
  const sessions = [
    {
      id: 1,
      user: 'John Doe',
      connector: 'Connector 2',
      startTime: '2024-01-07 14:30',
      duration: '2h 15m',
      energy: 45.2,
      status: 'completed'
    },
    {
      id: 2,
      user: 'Jane Smith',
      connector: 'Connector 1',
      startTime: '2024-01-07 10:15',
      duration: '1h 45m',
      energy: 32.8,
      status: 'completed'
    },
    {
      id: 3,
      user: 'Bob Johnson',
      connector: 'Connector 3',
      startTime: '2024-01-07 08:00',
      duration: '3h 30m',
      energy: 58.9,
      status: 'completed'
    }
  ];

  return (
    <div className="card p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
          {title}
        </h3>
        <button
          onClick={onRemove}
          className="p-1 rounded-lg hover:bg-gray-100 dark:hover:bg-anthracite-800 transition-colors"
        >
          <X className="h-4 w-4 text-gray-600 dark:text-gray-400" />
        </button>
      </div>

      <div className="space-y-3">
        {sessions.map((session) => (
          <div
            key={session.id}
            className="p-3 bg-gray-50 dark:bg-anthracite-800 rounded-lg"
          >
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center space-x-2">
                <User className="h-4 w-4 text-gray-600 dark:text-gray-400" />
                <span className="font-medium text-gray-900 dark:text-white">
                  {session.user}
                </span>
              </div>
              <span className="text-xs bg-green-100 dark:bg-green-900/20 text-green-800 dark:text-green-400 px-2 py-1 rounded">
                {session.status}
              </span>
            </div>
            
            <div className="grid grid-cols-2 gap-2 text-sm text-gray-600 dark:text-gray-400">
              <div className="flex items-center space-x-1">
                <Zap className="h-3 w-3" />
                <span>{session.energy} kWh</span>
              </div>
              <div className="flex items-center space-x-1">
                <Clock className="h-3 w-3" />
                <span>{session.duration}</span>
              </div>
            </div>
            
            <div className="mt-2 text-xs text-gray-500 dark:text-gray-500">
              {session.connector} • {session.startTime}
            </div>
          </div>
        ))}
      </div>

      <div className="mt-4 pt-4 border-t border-gray-200 dark:border-anthracite-700">
        <div className="flex justify-between text-sm text-gray-600 dark:text-gray-400">
          <span>Total Sessions: {sessions.length}</span>
          <span>Total Energy: {sessions.reduce((sum, s) => sum + s.energy, 0).toFixed(1)} kWh</span>
        </div>
      </div>
    </div>
  );
};

export default SessionsWidget; 