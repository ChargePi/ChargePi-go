import React from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import Sidebar from '@/components/Sidebar';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Settings, Activity, Power, MapPin, AlertTriangle } from 'lucide-react';


interface EVSE {
  id: string;
  name: string;
  status: 'online' | 'offline' | 'charging' | 'available';
  location: string;
  power?: string;
  currentSession?: string;
}

const EVSEDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const evse = location.state?.evse as EVSE;

  if (!evse) {
    return (
      <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
        <Sidebar isOpen={true} onClose={() => {}} />
        <header className="sticky top-0 z-30 bg-card border-b">
          <div className="flex items-center justify-between px-6 py-4">
            <div className="flex items-center space-x-3">
              <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
              <h1 className="text-xl font-semibold">EVSE Detail</h1>
            </div>
            <div className="flex items-center space-x-4">
              
            </div>
          </div>
        </header>
        <main className="flex-1 p-6 min-w-0">
          <div className="bg-card/80 border-2 border-primary/30 shadow-lg rounded-xl p-6">
            <div className="text-center py-8">
              <AlertTriangle className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h2 className="text-lg font-semibold mb-2">EVSE Not Found</h2>
              <p className="text-muted-foreground mb-4">The requested EVSE could not be found.</p>
              <Button onClick={() => navigate('/evse')}>
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to EVSEs
              </Button>
            </div>
          </div>
        </main>
      </div>
    );
  }

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

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">{evse.name}</h1>
          </div>
          <div className="flex items-center space-x-4">
            
          </div>
        </div>
      </header>
      <main className="flex-1 p-6 min-w-0">
        <div className="bg-card/80 border-2 border-primary/30 shadow-lg rounded-xl p-6">
          <div className="flex items-center justify-between mb-6">
            <Button variant="outline" onClick={() => navigate('/evse')}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to EVSEs
            </Button>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm">
                <Settings className="h-4 w-4 mr-2" />
                Settings
              </Button>
              <Button variant="outline" size="sm">
                <Activity className="h-4 w-4 mr-2" />
                Logs
              </Button>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Power className="h-5 w-5" />
                  EVSE Information
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-muted-foreground">Status</span>
                  <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(evse.status)}`}>
                    {evse.status.charAt(0).toUpperCase() + evse.status.slice(1)}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-muted-foreground">ID</span>
                  <span className="font-mono text-sm">{evse.id}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-muted-foreground">Location</span>
                  <span className="flex items-center gap-1">
                    <MapPin className="h-4 w-4" />
                    {evse.location}
                  </span>
                </div>
                {evse.power && (
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-muted-foreground">Power</span>
                    <span className="font-medium">{evse.power}</span>
                  </div>
                )}
                {evse.currentSession && (
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-muted-foreground">Current Session</span>
                    <span className="font-medium">{evse.currentSession}</span>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Activity className="h-5 w-5" />
                  Recent Activity
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center justify-between text-sm">
                    <span>Session started</span>
                    <span className="text-muted-foreground">2 hours ago</span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span>Status changed to charging</span>
                    <span className="text-muted-foreground">1 hour ago</span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span>Power output increased</span>
                    <span className="text-muted-foreground">30 min ago</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </main>
    </div>
  );
};

export default EVSEDetail; 