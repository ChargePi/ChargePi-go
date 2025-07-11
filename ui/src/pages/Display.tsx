import React, { useState } from 'react';
import Sidebar from '@/components/Sidebar';
import GenericTable, { TableColumn } from '@/components/generic-table';
import GenericCardView from '@/components/generic-card-view';
import ViewToggle from '@/components/view-toggle';
import { Copy } from 'lucide-react';

import SettingsSidebar from '@/components/settings-sidebar';
import { Button } from '@/components/ui/button';

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

const Display: React.FC = () => {
  const [displays, setDisplays] = useState<Display[]>(initialDisplays);
  const [view, setView] = useState<'list' | 'cards'>('list');
  const [copiedId, setCopiedId] = useState<string | null>(null);

  // Sidebar edit state
  const [editSidebarOpen, setEditSidebarOpen] = useState(false);
  const [editForm, setEditForm] = useState<Partial<Display>>({});
  const [editErrors, setEditErrors] = useState<Record<string, string>>({});

  const columns: TableColumn<Display>[] = [
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

  const validateField = (col: any, value: any) => {
    if (col.required && (!value || value === '')) {
      return 'This field is required.';
    }
    if (col.type === 'email' && value) {
      const emailRegex = /^[^\s@]+@[^ -\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(value)) return 'Invalid email address.';
    }
    if (col.type === 'number' && value) {
      if (isNaN(Number(value))) return 'Must be a number.';
    }
    if (col.type === 'date' && value) {
      if (isNaN(Date.parse(value))) return 'Invalid date.';
    }
    return '';
  };

  const validateEditForm = () => {
    const newErrors: Record<string, string> = {};
    columns.forEach(col => {
      if (col.editable === false) return;
      const value = editForm[col.key];
      const error = validateField(col, value);
      if (error) newErrors[col.key as string] = error;
    });
    setEditErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleRowClick = (display: Display) => {
    setEditForm(display);
    setEditErrors({});
    setEditSidebarOpen(true);
  };
  const handleEditSidebarClose = () => {
    setEditSidebarOpen(false);
    setEditForm({});
    setEditErrors({});
  };
  const handleEditSidebarChange = (key: keyof Display, value: any) => {
    setEditForm(f => ({ ...f, [key]: value }));
    setEditErrors(e => ({ ...e, [key as string]: '' }));
  };
  const handleEditSidebarSave = () => {
    if (!validateEditForm()) return;
    setDisplays(displays => displays.map(d => d.id === editForm.id ? { ...d, ...editForm } as Display : d));
    handleEditSidebarClose();
  };

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">Display Management</h1>
          </div>
          <div className="flex items-center space-x-4">
            
          </div>
        </div>
      </header>
      <main className="flex-1 p-6 min-w-0">
        <div className="bg-card/80 border-2 border-primary/30 shadow-lg rounded-xl p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">Displays</h2>
            <ViewToggle view={view} onViewChange={setView} />
          </div>
          <div className="overflow-x-auto">
            {view === 'list' ? (
              <GenericTable<Display>
                columns={columns}
                data={displays}
                filterInputClassName="w-1/4"
                className="min-w-full text-sm"
                columnVisibilityEnabled={true}
                columnSearchEnabled={true}
                onRowClick={handleRowClick}
              />
            ) : (
              <GenericCardView<Display>
                data={displays}
                columns={columns}
                className="min-w-full"
                onEdit={handleRowClick}
                onDelete={display => setDisplays(displays => displays.filter(d => d.id !== display.id))}
              />
            )}
          </div>
        </div>
        <SettingsSidebar
          open={editSidebarOpen}
          onClose={handleEditSidebarClose}
          title="Edit Display"
          onSave={handleEditSidebarSave}
          actions={
            <>
              <Button variant="outline" className="flex-1" onClick={handleEditSidebarClose}>
                Cancel
              </Button>
              <Button variant="accent" className="flex-1" onClick={handleEditSidebarSave}>
                Save
              </Button>
            </>
          }
        >
          <form className="space-y-4">
            {columns.filter(col => col.key !== 'id' && col.editable !== false).map(col => (
              <div key={String(col.key)}>
                <label className="block text-sm font-medium mb-1">{col.label}{col.required && <span className="text-red-500 ml-1">*</span>}</label>
                <input
                  type={col.type || 'text'}
                  className={`w-full px-3 py-2 border rounded bg-background ${editErrors[col.key as string] ? 'border-red-500' : 'border'} focus:outline-none`}
                  value={editForm[col.key] as string || ''}
                  onChange={e => handleEditSidebarChange(col.key, e.target.value)}
                  placeholder={`Enter ${col.label}`}
                  required={col.required}
                />
                {editErrors[col.key as string] && (
                  <p className="text-xs text-red-600 mt-1">{editErrors[col.key as string]}</p>
                )}
              </div>
            ))}
          </form>
        </SettingsSidebar>
      </main>
    </div>
  );
};

export default Display; 