import React, { useState } from 'react';
import Sidebar from '@/components/Sidebar';
import EVSECard from '@/components/evse-card';
import ViewToggle from '@/components/view-toggle';
import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import AddSidebarButton from '@/components/add-sidebar-button';
import { getUserPrefs } from '@/utils/dateUtils';

import { useNavigate } from 'react-router-dom';

interface EVSE {
  id: string;
  name: string;
  status: 'online' | 'offline' | 'charging' | 'available';
  location: string;
  power?: string;
  currentSession?: string;
}

const initialEVSEs: EVSE[] = [
  { id: '1', name: 'EVSE-001', status: 'available', location: 'Parking A', power: '22 kW' },
  { id: '2', name: 'EVSE-002', status: 'charging', location: 'Parking B', power: '11 kW', currentSession: 'Active - 45 min' },
  { id: '3', name: 'EVSE-003', status: 'offline', location: 'Parking C', power: '22 kW' },
  { id: '4', name: 'EVSE-004', status: 'online', location: 'Parking D', power: '11 kW' },
];


function formatPower(value: string | undefined, unit: string) {
  if (!value) return '-';
  const num = parseFloat(value);
  if (isNaN(num)) return value;
  if (unit === 'kW') return `${num} kW`;
  if (unit === 'W') return `${(num * 1000).toFixed(0)} W`;
  return value;
}

const EVSE: React.FC = () => {
  const [evses, setEvses] = useState<EVSE[]>(initialEVSEs);
  const [view, setView] = useState<'list' | 'cards'>('cards');
  const navigate = useNavigate();

  const handleEVSEClick = (evse: EVSE) => {
    navigate(`/evse/${evse.id}`, { state: { evse } });
  };

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">EVSE Management</h1>
          </div>
          <div className="flex items-center space-x-4">
            
          </div>
        </div>
      </header>
      <main className="flex-1 p-6 min-w-0">
        <div className="bg-card/80 border-2 border-primary/30 shadow-lg rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">EVSEs</h2>
            <div className="flex items-center gap-4">
              <ViewToggle view={view} onViewChange={setView} />
              <Button variant="accent" size="sm">
                <Plus className="h-4 w-4 mr-2" />
                Add EVSE
              </Button>
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {evses.map((evse) => (
              <EVSECard key={evse.id} evse={evse} onClick={handleEVSEClick} />
            ))}
          </div>
        </div>
      </main>
    </div>
  );
};

export default EVSE; 