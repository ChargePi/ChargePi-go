import React from 'react';
import { Activity, Zap, BarChart3, AlertTriangle, Clock } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Card, CardContent } from '@/components/ui/card';

interface Widget {
  id: string;
  type: string;
  title: string;
  size: 'small' | 'medium' | 'large';
  x: number;
  y: number;
  width: number;
  height: number;
}

interface WidgetSelectorProps {
  open: boolean;
  onClose: () => void;
  currentWidgets: Widget[];
  onAddWidget: (type: string) => void;
}

const WidgetSelector: React.FC<WidgetSelectorProps> = ({ open, onClose, currentWidgets, onAddWidget }) => {
  const availableWidgets = [
    {
      id: 'status',
      title: 'Charge Point Status',
      description: 'Current status of all connectors',
      icon: Activity,
      size: 'medium',
    },
    {
      id: 'energy',
      title: 'Energy Consumption',
      description: 'Energy consumption statistics and charts',
      icon: BarChart3,
      size: 'large',
    },
    {
      id: 'sessions',
      title: 'Recent Sessions',
      description: 'List of recent charging sessions',
      icon: Clock,
      size: 'medium',
    },
    {
      id: 'diagnostics',
      title: 'Diagnostics',
      description: 'System diagnostics and health monitoring',
      icon: AlertTriangle,
      size: 'small',
    },
    {
      id: 'current-session',
      title: 'Current Session',
      description: 'Details of ongoing charging session',
      icon: Zap,
      size: 'medium',
    },
  ];

  const handleAddWidget = (widgetId: string) => {
    onAddWidget(widgetId);
  };

  const isWidgetAdded = (widgetType: string) => {
    return currentWidgets.some(widget => widget.type === widgetType);
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Add Widget</DialogTitle>
        </DialogHeader>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {availableWidgets.map((widget) => {
            const isAdded = isWidgetAdded(widget.id);
            return (
            <Card
              key={widget.id}
              className={`p-4 transition-shadow ${
                isAdded 
                  ? 'opacity-50 cursor-not-allowed bg-muted' 
                  : 'hover:shadow-md cursor-pointer'
              }`}
              onClick={() => !isAdded && handleAddWidget(widget.id)}
            >
              <CardContent className="p-0">
                <div className="flex items-start space-x-3">
                  <div className="p-2 bg-accent-orange/10 rounded-lg">
                    <widget.icon className="h-6 w-6 text-accent-orange" />
                  </div>
                  <div className="flex-1">
                    <h3 className="font-medium">
                      {widget.title}
                    </h3>
                    <p className="text-sm text-muted-foreground mt-1">
                      {widget.description}
                    </p>
                    <div className="mt-2 flex items-center space-x-2">
                      <span className="inline-block px-2 py-1 text-xs bg-secondary text-secondary-foreground rounded">
                        {widget.size}
                      </span>
                      {isAdded && (
                        <span className="inline-block px-2 py-1 text-xs bg-green-100 dark:bg-green-900/20 text-green-800 dark:text-green-400 rounded">
                          Added
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
            );
          })}
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default WidgetSelector; 