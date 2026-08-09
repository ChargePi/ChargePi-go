import React from 'react';
import { X, AlertTriangle, AlertCircle, XCircle, CheckCircle, Info } from 'lucide-react';

interface DiagnosticsWidgetProps {
  title: string;
  onRemove: () => void;
}

const DiagnosticsWidget: React.FC<DiagnosticsWidgetProps> = ({ title, onRemove }) => {
  const diagnostics = [
    {
      id: 1,
      type: 'error',
      message: 'Connector 3 communication failure',
      component: 'Connector 3',
      timestamp: '2024-01-07 15:30',
      severity: 'high',
      status: 'active'
    },
    {
      id: 2,
      type: 'warning',
      message: 'Temperature sensor reading high',
      component: 'Connector 1',
      timestamp: '2024-01-07 14:15',
      severity: 'medium',
      status: 'active'
    },
    {
      id: 3,
      type: 'info',
      message: 'System maintenance due next week',
      component: 'System',
      timestamp: '2024-01-07 12:00',
      severity: 'low',
      status: 'acknowledged'
    }
  ];

  const getDiagnosticIcon = (type: string) => {
    switch (type) {
      case 'error':
        return <XCircle className="h-4 w-4 text-red-500" />;
      case 'warning':
        return <AlertTriangle className="h-4 w-4 text-yellow-500" />;
      case 'info':
        return <Info className="h-4 w-4 text-blue-500" />;
      case 'success':
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      default:
        return <AlertCircle className="h-4 w-4 text-gray-500" />;
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'high':
        return 'bg-red-100 dark:bg-red-900/20 text-red-800 dark:text-red-400';
      case 'medium':
        return 'bg-yellow-100 dark:bg-yellow-900/20 text-yellow-800 dark:text-yellow-400';
      case 'low':
        return 'bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-400';
      default:
        return 'bg-gray-100 dark:bg-gray-900/20 text-gray-800 dark:text-gray-400';
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

  const activeCount = diagnostics.filter(d => d.status === 'active').length;
  const highPriorityCount = diagnostics.filter(d => d.severity === 'high' && d.status === 'active').length;

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

      {diagnostics.length === 0 ? (
        <div className="text-center py-8">
          <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-3" />
          <p className="text-gray-600 dark:text-gray-400">No active diagnostics</p>
        </div>
      ) : (
        <div className="space-y-3">
          {diagnostics.slice(0, 3).map((diagnostic) => (
            <div
              key={diagnostic.id}
              className="p-3 bg-gray-50 dark:bg-anthracite-800 rounded-lg border-l-4 border-gray-300"
            >
              <div className="flex items-start space-x-3">
                {getDiagnosticIcon(diagnostic.type)}
                <div className="flex-1">
                  <div className="flex items-center justify-between mb-1">
                    <span className="font-medium text-gray-900 dark:text-white">
                      {diagnostic.component}
                    </span>
                    <div className="flex items-center space-x-1">
                      <span className={`text-xs px-2 py-1 rounded ${getSeverityColor(diagnostic.severity)}`}>
                        {diagnostic.severity}
                      </span>
                      <span className={`text-xs px-2 py-1 rounded ${getStatusColor(diagnostic.status)}`}>
                        {diagnostic.status}
                      </span>
                    </div>
                  </div>
                  <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                    {diagnostic.message}
                  </p>
                  <span className="text-xs text-gray-500 dark:text-gray-500">
                    {diagnostic.timestamp}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="mt-4 pt-4 border-t border-gray-200 dark:border-anthracite-700">
        <div className="flex justify-between text-sm text-gray-600 dark:text-gray-400">
          <span>Active Issues: {activeCount}</span>
          <span>High Priority: {highPriorityCount}</span>
        </div>
        {diagnostics.length > 3 && (
          <div className="mt-2 text-xs text-gray-500 dark:text-gray-500">
            Showing 3 of {diagnostics.length} diagnostics
          </div>
        )}
      </div>
    </div>
  );
};

export default DiagnosticsWidget; 