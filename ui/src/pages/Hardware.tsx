import React, { useState } from 'react';
import Sidebar from '@/components/Sidebar';
import GenericTable, { TableColumn } from '@/components/generic-table';
import GenericCardView from '@/components/generic-card-view';
import ViewToggle from '@/components/view-toggle';
import { Copy, Zap, Monitor, CreditCard } from 'lucide-react';

import EVSECard from '@/components/evse-card';
import { Button } from '@/components/ui/button';
import AddSidebarButton from '@/components/add-sidebar-button';

// Data and columns for Display
interface Display {
  id: string;
  name: string;
  status: 'online' | 'offline';
  location: string;
}
const initialDisplays: Display[] = [
  { id: '1', name: 'Display A', status: 'online', location: 'Lobby' },
  { id: '2', name: 'Display B', status: 'offline', location: 'Entrance' },
];

// Data and columns for Reader
interface Reader {
  id: string;
  name: string;
  status: 'online' | 'offline';
  location: string;
}
const initialReaders: Reader[] = [
  { id: '1', name: 'Reader A', status: 'online', location: 'Lobby' },
  { id: '2', name: 'Reader B', status: 'offline', location: 'Entrance' },
];

// Data for EVSE
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

const Hardware: React.FC = () => {
  const [tab, setTab] = useState<'evse' | 'display' | 'reader'>('evse');
  const [displays, setDisplays] = useState<Display[]>(initialDisplays);
  const [readers, setReaders] = useState<Reader[]>(initialReaders);
  const [evses, setEvses] = useState<EVSE[]>(initialEVSEs);
  const [displayView, setDisplayView] = useState<'list' | 'cards'>('list');
  const [readerView, setReaderView] = useState<'list' | 'cards'>('list');
  const [evseView, setEvseView] = useState<'cards'>('cards');
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const displayColumns: TableColumn<Display>[] = [
    {
      key: 'id',
      label: 'ID',
      filterable: true,
      editable: false,
      required: false,
      type: 'text',
      render: (id: string) => (
        <span className="flex items-center gap-2 relative">
          <span>{id}</span>
          <button
            className="p-1 rounded hover:bg-muted relative"
            onClick={e => {
              e.stopPropagation();
              navigator.clipboard.writeText(id);
              setCopiedId(id);
              setTimeout(() => setCopiedId(current => (current === id ? null : current)), 2000);
            }}
            title="Copy ID"
          >
            <Copy className="h-4 w-4" />
            {copiedId === id && (
              <span className="absolute left-1/2 -translate-x-1/2 top-8 bg-black text-white text-xs rounded px-2 py-1 shadow z-10 animate-fade-in">
                Copied!
              </span>
            )}
          </button>
        </span>
      ),
    },
    { key: 'name', label: 'Name', filterable: true, required: true, type: 'text' },
    { key: 'status', label: 'Status', filterable: true, required: true, type: 'text' },
    { key: 'location', label: 'Location', filterable: true, required: true, type: 'text' },
  ];

  const readerColumns: TableColumn<Reader>[] = [
    {
      key: 'id',
      label: 'ID',
      filterable: true,
      editable: false,
      required: false,
      type: 'text',
      render: (id: string) => (
        <span className="flex items-center gap-2 relative">
          <span>{id}</span>
          <button
            className="p-1 rounded hover:bg-muted relative"
            onClick={e => {
              e.stopPropagation();
              navigator.clipboard.writeText(id);
              setCopiedId(id);
              setTimeout(() => setCopiedId(current => (current === id ? null : current)), 2000);
            }}
            title="Copy ID"
          >
            <Copy className="h-4 w-4" />
            {copiedId === id && (
              <span className="absolute left-1/2 -translate-x-1/2 top-8 bg-black text-white text-xs rounded px-2 py-1 shadow z-10 animate-fade-in">
                Copied!
              </span>
            )}
          </button>
        </span>
      ),
    },
    { key: 'name', label: 'Name', filterable: true, required: true, type: 'text' },
    { key: 'status', label: 'Status', filterable: true, required: true, type: 'text' },
    { key: 'location', label: 'Location', filterable: true, required: true, type: 'text' },
  ];

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">Hardware</h1>
          </div>
          <div className="flex items-center space-x-4">
            
          </div>
        </div>
      </header>
      <main className="flex-1 p-6 min-w-0">
        <nav className="flex items-center justify-between mb-6" aria-label="Hardware sections">
          <div className="flex gap-4 items-center">
            <button
              className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${tab === 'evse' ? 'bg-accent-orange/20 text-accent-orange' : 'hover:bg-muted text-muted-foreground'}`}
              onClick={() => setTab('evse')}
              aria-current={tab === 'evse'}
            >
              <Zap className="h-5 w-5" /> EVSE
            </button>
            <button
              className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${tab === 'display' ? 'bg-accent-orange/20 text-accent-orange' : 'hover:bg-muted text-muted-foreground'}`}
              onClick={() => setTab('display')}
              aria-current={tab === 'display'}
            >
              <Monitor className="h-5 w-5" /> Display
            </button>
            <button
              className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${tab === 'reader' ? 'bg-accent-orange/20 text-accent-orange' : 'hover:bg-muted text-muted-foreground'}`}
              onClick={() => setTab('reader')}
              aria-current={tab === 'reader'}
            >
              <CreditCard className="h-5 w-5" /> Reader
            </button>
          </div>
          <div className="flex items-center gap-2">
            {tab === 'evse' && (
              <AddSidebarButton<EVSE>
                buttonTitle="Add EVSE"
                sidebarTitle="Add New EVSE"
                columns={[
                  { key: 'name', label: 'Name', required: true, type: 'text' },
                  { key: 'location', label: 'Location', required: true, type: 'text' },
                ]}
                onAdd={evse => setEvses(prev => [...prev, evse])}
              />
            )}
            {tab === 'display' && (
              <>
                <ViewToggle view={displayView} onViewChange={setDisplayView} />
                <AddSidebarButton<Display>
                  buttonTitle="Add Display"
                  sidebarTitle="Add New Display"
                  columns={displayColumns.filter(col => col.key !== 'status')}
                  onAdd={display => setDisplays(prev => [...prev, display])}
                />
              </>
            )}
            {tab === 'reader' && (
              <>
                <ViewToggle view={readerView} onViewChange={setReaderView} />
                <AddSidebarButton<Reader>
                  buttonTitle="Add Reader"
                  sidebarTitle="Add New Reader"
                  columns={readerColumns.filter(col => col.key !== 'status')}
                  onAdd={reader => setReaders(prev => [...prev, reader])}
                />
              </>
            )}
          </div>
        </nav>
        {tab === 'evse' && (
          <section aria-labelledby="evse-section">
            {/* Removed ViewToggle and view state for EVSE section */}
            <div className="mb-4"></div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {evses.map((evse) => (
                <EVSECard key={evse.id} evse={evse} onClick={() => {}} />
              ))}
            </div>
          </section>
        )}
        {tab === 'display' && (
          <section aria-labelledby="display-section">
            {/* Removed ViewToggle from here, now in nav */}
            <div className="mb-4"></div>
            <div className="overflow-x-auto">
              {displayView === 'list' ? (
                <GenericTable<Display>
                  columns={displayColumns}
                  data={displays}
                  filterInputClassName="w-1/4"
                  className="min-w-full text-sm"
                  columnVisibilityEnabled={true}
                  columnSearchEnabled={true}
                />
              ) : (
                <GenericCardView<Display>
                  data={displays}
                  columns={displayColumns}
                  className="min-w-full"
                />
              )}
            </div>
          </section>
        )}
        {tab === 'reader' && (
          <section aria-labelledby="reader-section">
            {/* Removed ViewToggle from here, now in nav */}
            <div className="mb-4"></div>
            <div className="overflow-x-auto">
              {readerView === 'list' ? (
                <GenericTable<Reader>
                  columns={readerColumns}
                  data={readers}
                  filterInputClassName="w-1/4"
                  className="min-w-full text-sm"
                  columnVisibilityEnabled={true}
                  columnSearchEnabled={true}
                />
              ) : (
                <GenericCardView<Reader>
                  data={readers}
                  columns={readerColumns}
                  className="min-w-full"
                />
              )}
            </div>
          </section>
        )}
      </main>
    </div>
  );
};

export default Hardware; 