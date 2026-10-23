import React, { useState, useMemo } from 'react';
import Sidebar from '@/components/Sidebar';

import ViewToggle from '@/components/view-toggle';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Zap, Clock, BatteryCharging, CreditCard } from 'lucide-react';
import GenericTable, { TableColumn } from '@/components/generic-table';
import GenericCardView from '@/components/generic-card-view';
import { Button } from '@/components/ui/button';
import { formatDate, getUserPrefs } from '@/utils/dateUtils';

interface Session {
  id: string;
  evse: string;
  status: 'active' | 'completed' | 'error';
  consumption: number; // kWh
  duration: number; // seconds
  rfidTag: string;
  start: string; // ISO
  end?: string; // ISO or undefined
}

const mockSessions: Session[] = [
  {
    id: 'sess-1',
    evse: 'EVSE-001',
    status: 'active',
    consumption: 12.4,
    duration: 3600,
    rfidTag: '123456',
    start: '2024-06-01T10:00:00Z',
  },
  {
    id: 'sess-2',
    evse: 'EVSE-002',
    status: 'completed',
    consumption: 8.2,
    duration: 2700,
    rfidTag: '654321',
    start: '2024-06-01T09:00:00Z',
    end: '2024-06-01T09:45:00Z',
  },
  {
    id: 'sess-3',
    evse: 'EVSE-001',
    status: 'completed',
    consumption: 15.1,
    duration: 5400,
    rfidTag: '123456',
    start: '2024-05-31T18:00:00Z',
    end: '2024-05-31T19:30:00Z',
  },
];


function formatEnergy(value: number, unit: string) {
  if (unit === 'kWh') return `${value.toFixed(2)} kWh`;
  if (unit === 'Wh') return `${(value * 1000).toFixed(0)} Wh`;
  return `${value} kWh`;
}

const Sessions: React.FC = () => {
  const [view, setView] = useState<'list' | 'cards'>('list');
  const [sessions] = useState<Session[]>(mockSessions);
  const userPrefs = getUserPrefs();

  const columns: TableColumn<Session>[] = [
    { key: 'id', label: 'Session ID', filterable: true, editable: false, required: false, type: 'text' },
    { key: 'evse', label: 'EVSE', filterable: true, required: true, type: 'text' },
    { key: 'status', label: 'Status', filterable: true, required: true, type: 'text', render: v => v.charAt(0).toUpperCase() + v.slice(1) },
    { key: 'consumption', label: `Consumption (${userPrefs.energyUnit})`, filterable: false, required: true, type: 'number', render: v => formatEnergy(v, userPrefs.energyUnit) },
    { key: 'duration', label: 'Duration', filterable: false, required: true, type: 'number', render: v => `${Math.floor(v/60)}m ${v%60}s` },
    { key: 'rfidTag', label: 'RFID Tag', filterable: true, required: true, type: 'text' },
    { key: 'start', label: 'Start', filterable: false, required: true, type: 'date', render: v => formatDate(v, userPrefs.timezone) },
    { key: 'end', label: 'End', filterable: false, required: false, type: 'date', render: v => v ? formatDate(v, userPrefs.timezone) : '-' },
  ];

  // Statistics
  const stats = useMemo(() => {
    const activeSessions = sessions.filter(s => s.status === 'active');
    const totalConsumption = activeSessions.reduce((sum, s) => sum + s.consumption, 0);
    const chargingTimePerEvse: Record<string, number> = {};
    activeSessions.forEach(s => {
      chargingTimePerEvse[s.evse] = (chargingTimePerEvse[s.evse] || 0) + s.duration;
    });
    return {
      activeCount: activeSessions.length,
      totalConsumption,
      chargingTimePerEvse,
    };
  }, [sessions]);

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <Zap className="h-8 w-8 text-accent-orange" />
            <h1 className="text-xl font-semibold">Sessions</h1>
          </div>
          <div className="flex items-center space-x-4">
            
          </div>
        </div>
      </header>
      <main className="flex-1 p-6 min-w-0">
        {/* Statistics */}
        <section className="mb-8 grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2"><BatteryCharging className="h-5 w-5 text-accent-orange" />Active Sessions</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold">{stats.activeCount}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2"><Zap className="h-5 w-5 text-accent-orange" />Current Consumption</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold">{formatEnergy(stats.totalConsumption, userPrefs.energyUnit)}</div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2"><Clock className="h-5 w-5 text-accent-orange" />Charging Time per EVSE</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="space-y-1">
                {Object.keys(stats.chargingTimePerEvse).map((evse) => {
                  const time = stats.chargingTimePerEvse[evse];
                  return (
                    <li key={evse} className="flex justify-between text-sm">
                      <span>{evse}</span>
                      <span>{Math.floor(time/60)}m {time%60}s</span>
                    </li>
                  );
                })}
                {Object.keys(stats.chargingTimePerEvse).length === 0 && <li className="text-muted-foreground">No active sessions</li>}
              </ul>
            </CardContent>
          </Card>
        </section>
        {/* List/Card view toggle and session list */}
        <section>
          <div className="flex justify-end items-center mb-4">
            <ViewToggle view={view} onViewChange={setView} />
          </div>
          {view === 'list' ? (
            <GenericTable<Session>
              columns={columns}
              data={sessions}
              filterInputClassName="w-1/4"
              className="min-w-full text-sm"
              columnVisibilityEnabled={true}
              columnSearchEnabled={true}
            />
          ) : (
            <GenericCardView<Session>
              data={sessions}
              columns={columns}
              className="min-w-full"
            />
          )}
        </section>
      </main>
    </div>
  );
};

export default Sessions; 