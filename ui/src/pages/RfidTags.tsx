import React, { useState } from 'react';
import Sidebar from '@/components/Sidebar';
import { Table } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Copy, Plus } from 'lucide-react';
import { Dialog, DialogContent, DialogOverlay } from '@/components/ui/dialog';
import AddSidebarButton from '@/components/add-sidebar-button';
import GenericTable, { TableColumn } from '@/components/generic-table';
import GenericCardView from '@/components/generic-card-view';
import ViewToggle from '@/components/view-toggle';
import { formatDate, getUserPrefs } from '@/utils/dateUtils';

import SettingsSidebar from '@/components/settings-sidebar';

interface Tag {
  id: string;
  owner: string;
  status: 'active' | 'inactive' | 'suspended' | 'expired';
  expiryDate?: string; // ISO date string, optional
}

const TAG_STATUS_OPTIONS = [
  { value: 'active', label: 'Active' },
  { value: 'inactive', label: 'Inactive' },
  { value: 'suspended', label: 'Suspended' },
  { value: 'expired', label: 'Expired' },
];

const initialTags: Tag[] = [
  { id: '123456', owner: 'John Doe', status: 'active', expiryDate: '2025-12-31' },
  { id: '654321', owner: 'Jane Smith', status: 'inactive' },
];



const RfidTags: React.FC = () => {
  const [tags, setTags] = useState<Tag[]>(initialTags);
  const [filter, setFilter] = useState('');
  const [selected, setSelected] = useState<Tag | null>(null);
  const [newTag, setNewTag] = useState<Partial<Tag>>({ status: 'active' });
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [view, setView] = useState<'list' | 'cards'>('list');

  // Sidebar edit state
  const [editSidebarOpen, setEditSidebarOpen] = useState(false);
  const [editForm, setEditForm] = useState<Partial<Tag>>({});
  const [editErrors, setEditErrors] = useState<Record<string, string>>({});

  const handleAddTag = () => {
    if (!newTag.id || !newTag.owner) return;
    setTags([...tags, { ...newTag, id: newTag.id.toString() } as Tag]);
    setNewTag({ status: 'active' });
  };

  const validateField = (col: any, value: any) => {
    if (col.required && (!value || value === '')) {
      return 'This field is required.';
    }
    if (col.type === 'email' && value) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
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

  const handleRowClick = (tag: Tag) => {
    setEditForm(tag);
    setEditErrors({});
    setEditSidebarOpen(true);
  };
  const handleEditSidebarClose = () => {
    setEditSidebarOpen(false);
    setEditForm({});
    setEditErrors({});
  };
  const handleEditSidebarChange = (key: keyof Tag, value: any) => {
    setEditForm(f => ({ ...f, [key]: value }));
    setEditErrors(e => ({ ...e, [key as string]: '' }));
  };
  const handleEditSidebarSave = () => {
    if (!validateEditForm()) return;
    setTags(tags => tags.map(t => t.id === editForm.id ? { ...t, ...editForm } as Tag : t));
    handleEditSidebarClose();
  };

  const columns: TableColumn<Tag>[] = [
    {
      key: 'id',
      label: 'Tag ID',
      filterable: true,
      editable: false,
      required: false,
      type: 'text',
      render: (id: string, row: Tag) => (
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
            title="Copy Tag ID"
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
    { key: 'owner', label: 'Owner', filterable: true, required: true, type: 'text' },
    { key: 'status', label: 'Status', filterable: true, required: true, type: 'select', options: TAG_STATUS_OPTIONS, render: (v: string) => v.charAt(0).toUpperCase() + v.slice(1) },
    {
      key: 'expiryDate',
      label: 'Expiry Date',
      filterable: true,
      required: false,
      type: 'date',
      render: (v?: string) => v ? formatDate(v, getUserPrefs().timezone) : '-',
    },
  ];

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">RFID Tags</h1>
          </div>
          <div className="flex items-center space-x-4">
            
          </div>
        </div>
      </header>
      <main className="flex-1 p-6 min-w-0">
        <div className="flex justify-end items-center gap-4 mb-4">
              <ViewToggle view={view} onViewChange={setView} />
              <AddSidebarButton<Tag>
                buttonTitle="Add RFID Tag"
                sidebarTitle="Add New RFID Tag"
                columns={columns}
                onAdd={tag => setTags([...tags, tag])}
              />
          </div>
          {view === 'list' ? (
            <GenericTable<Tag>
              columns={columns}
              data={tags}
              filterInputClassName="w-1/4"
              className="min-w-full text-sm"
              columnVisibilityEnabled={true}
              columnSearchEnabled={true}
            onRowClick={handleRowClick}
            />
          ) : (
            <GenericCardView<Tag>
              data={tags}
              columns={columns}
              className="min-w-full"
            onEdit={handleRowClick}
            onDelete={tag => setTags(tags => tags.filter(t => t.id !== tag.id))}
          />
        )}
        <SettingsSidebar
          open={editSidebarOpen}
          onClose={handleEditSidebarClose}
          title="Edit RFID Tag"
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
                {col.type === 'select' && col.options ? (
                  <Select value={editForm[col.key] as string || ''} onValueChange={value => handleEditSidebarChange(col.key, value)}>
                    <SelectTrigger className={editErrors[col.key as string] ? 'border-red-500' : ''}>
                      <SelectValue placeholder={`Select ${col.label}`} />
                    </SelectTrigger>
                    <SelectContent>
                      {col.options.map(option => (
                        <SelectItem key={option.value} value={option.value}>
                          {option.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                ) : (
                  <input
                    type={col.type || 'text'}
                    className={`w-full px-3 py-2 border rounded bg-background ${editErrors[col.key as string] ? 'border-red-500' : 'border'} focus:outline-none`}
                    value={editForm[col.key] as string || ''}
                    onChange={e => handleEditSidebarChange(col.key, e.target.value)}
                    placeholder={`Enter ${col.label}`}
                    required={col.required}
                  />
                )}
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

export default RfidTags; 