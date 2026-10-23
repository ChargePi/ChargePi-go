import React from 'react';
import { X, Zap, Clock, User, Battery } from 'lucide-react';

interface CurrentSessionWidgetProps {
  title: string;
  onRemove: () => void;
}

const CurrentSessionWidget: React.FC<CurrentSessionWidgetProps> = ({ title, onRemove }) => {
  const currentSession = {
    id: 1,
    user: 'Alice Johnson',
    connector: 'Connector 2',
    startTime: '2024-01-07 16:00',
    duration: '1h 23m',
    energyDelivered: 28.5,
    currentPower: 22.4,
    targetEnergy: 50.0,
    progress: 57
  };

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

      {currentSession ? (
        <div className="space-y-4">
          {/* User Info */}
          <div className="flex items-center space-x-3 p-3 bg-accent-orange/10 rounded-lg">
            <User className="h-5 w-5 text-accent-orange" />
            <div>
              <p className="font-medium text-gray-900 dark:text-white">
                {currentSession.user}
              </p>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                {currentSession.connector}
              </p>
            </div>
          </div>

          {/* Progress Bar */}
          <div>
            <div className="flex justify-between text-sm mb-2">
              <span className="text-gray-600 dark:text-gray-400">Progress</span>
              <span className="font-medium text-gray-900 dark:text-white">
                {currentSession.progress}%
              </span>
            </div>
            <div className="w-full bg-gray-200 dark:bg-anthracite-700 rounded-full h-3">
              <div
                className="bg-gradient-to-r from-accent-orange to-accent-orange/80 h-3 rounded-full transition-all duration-300"
                style={{ width: `${currentSession.progress}%` }}
              />
            </div>
          </div>

          {/* Stats Grid */}
          <div className="grid grid-cols-2 gap-4">
            <div className="text-center p-3 bg-gray-50 dark:bg-anthracite-800 rounded-lg">
              <div className="flex items-center justify-center space-x-1 mb-1">
                <Zap className="h-4 w-4 text-accent-orange" />
                <span className="text-sm text-gray-600 dark:text-gray-400">Energy</span>
              </div>
              <p className="text-lg font-bold text-gray-900 dark:text-white">
                {currentSession.energyDelivered.toFixed(1)} kWh
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-500">
                of {currentSession.targetEnergy} kWh
              </p>
            </div>

            <div className="text-center p-3 bg-gray-50 dark:bg-anthracite-800 rounded-lg">
              <div className="flex items-center justify-center space-x-1 mb-1">
                <Battery className="h-4 w-4 text-blue-500" />
                <span className="text-sm text-gray-600 dark:text-gray-400">Power</span>
              </div>
              <p className="text-lg font-bold text-gray-900 dark:text-white">
                {currentSession.currentPower} kW
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-500">
                Current
              </p>
            </div>
          </div>

          {/* Session Info */}
          <div className="flex items-center justify-between text-sm text-gray-600 dark:text-gray-400">
            <div className="flex items-center space-x-1">
              <Clock className="h-4 w-4" />
              <span>Duration: {currentSession.duration}</span>
            </div>
            <span>Started: {currentSession.startTime}</span>
          </div>
        </div>
      ) : (
        <div className="text-center py-8">
          <Zap className="h-12 w-12 text-gray-400 mx-auto mb-3" />
          <p className="text-gray-600 dark:text-gray-400">No active session</p>
        </div>
      )}
    </div>
  );
};

export default CurrentSessionWidget; 