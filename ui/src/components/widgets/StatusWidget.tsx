import React, { useState, useEffect } from 'react';
import { X, Zap, CheckCircle, XCircle, Clock } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { api } from '@/services/api';

interface StatusWidgetProps {
  title: string;
  onRemove: () => void;
}

const StatusWidget: React.FC<StatusWidgetProps> = ({ title, onRemove }) => {
  const [status, setStatus] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const data = await api.getStatus();
        setStatus(data);
      } catch (error) {
        console.error('Failed to fetch status:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchStatus();
  }, []);

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'available':
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      case 'charging':
        return <Zap className="h-4 w-4 text-accent-orange" />;
      case 'faulted':
        return <XCircle className="h-4 w-4 text-accent-red" />;
      default:
        return <Clock className="h-4 w-4 text-muted-foreground" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'available':
        return 'text-green-600 dark:text-green-400';
      case 'charging':
        return 'text-accent-orange';
      case 'faulted':
        return 'text-accent-red';
      default:
        return 'text-muted-foreground';
    }
  };

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{title}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-accent-orange"></div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-lg">{title}</CardTitle>
        <Button
          variant="ghost"
          size="icon"
          onClick={onRemove}
        >
          <X className="h-4 w-4" />
        </Button>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {status?.connectors?.map((connector: any) => (
            <div
              key={connector.id}
              className="flex items-center justify-between p-3 bg-muted rounded-lg"
            >
              <div className="flex items-center space-x-3">
                {getStatusIcon(connector.status)}
                <span className="font-medium">
                  {connector.name}
                </span>
              </div>
              <span className={`text-sm font-medium capitalize ${getStatusColor(connector.status)}`}>
                {connector.status}
              </span>
            </div>
          ))}
        </div>

        <div className="mt-4 pt-4 border-t">
          <div className="flex justify-between text-sm text-muted-foreground">
            <span>Total: {status?.totalConnectors || 0}</span>
            <span>Available: {status?.availableConnectors || 0}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default StatusWidget; 