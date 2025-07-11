import React, { useState } from 'react';
import Sidebar from '@/components/Sidebar';
import GenericTable, { TableColumn } from '@/components/generic-table';
import GenericCardView from '@/components/generic-card-view';
import ViewToggle from '@/components/view-toggle';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { MoreVertical, Edit, Trash2, Check, X, Copy, Plus } from 'lucide-react';
import AddSidebarButton from '@/components/add-sidebar-button';

import SettingsSidebar from '@/components/settings-sidebar';

interface User {
  id: string;
  username: string;
  role: string;
  email: string;
  status: 'active' | 'inactive';
}

const USER_ROLES = [
  { value: 'Admin', label: 'Administrator' },
  { value: 'User', label: 'User' },
  { value: 'Operator', label: 'Operator' },
  { value: 'Viewer', label: 'Viewer' },
];

const initialUsers: User[] = [
  { id: '1', username: 'admin', role: 'Admin', email: 'admin@example.com', status: 'active' },
  { id: '2', username: 'jane', role: 'User', email: 'jane@example.com', status: 'active' },
  { id: '3', username: 'bob', role: 'User', email: 'bob@example.com', status: 'inactive' },
];

const Users: React.FC = () => {
  const [users, setUsers] = useState<User[]>(initialUsers);
  const [filter, setFilter] = useState('');
  const [selected, setSelected] = useState<User | null>(null);
  const [editing, setEditing] = useState<User | null>(null);
  const [showDropdown, setShowDropdown] = useState<string | null>(null);
  const [showDelete, setShowDelete] = useState<User | null>(null);
  const [editDraft, setEditDraft] = useState<Partial<User>>({});
  const [newUser, setNewUser] = useState<Partial<User>>({ status: 'active' });
  const [view, setView] = useState<'list' | 'cards'>('list');
  const [copiedId, setCopiedId] = useState<string | null>(null);

  // Sidebar edit state
  const [editSidebarOpen, setEditSidebarOpen] = useState(false);
  const [editForm, setEditForm] = useState<Partial<User>>({});
  const [editErrors, setEditErrors] = useState<Record<string, string>>({});

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

  const handleEdit = (user: User) => {
    setEditing(user);
    setEditDraft(user);
    setShowDropdown(null);
  };

  const handleEditChange = (key: keyof User, value: string) => {
    setEditDraft(draft => ({ ...draft, [key]: value }));
  };

  const handleEditSave = () => {
    if (!editing) return;
    setUsers(users => users.map(u => u.id === editing.id ? { ...u, ...editDraft } as User : u));
    setEditing(null);
    setEditDraft({});
  };

  const handleEditCancel = () => {
    setEditing(null);
    setEditDraft({});
  };

  const handleDelete = () => {
    if (!showDelete) return;
    setUsers(users => users.filter(u => u.id !== showDelete.id));
    setShowDelete(null);
    setSelected(null);
  };

  const handleAddUser = () => {
    if (!newUser.id || !newUser.username || !newUser.role) return;
    setUsers([...users, { ...newUser, id: newUser.id.toString() } as User]);
    setNewUser({ status: 'active' });
  };

  const handleRowClick = (user: User) => {
    setEditForm(user);
    setEditErrors({});
    setEditSidebarOpen(true);
  };
  const handleEditSidebarClose = () => {
    setEditSidebarOpen(false);
    setEditForm({});
    setEditErrors({});
  };
  const handleEditSidebarChange = (key: keyof User, value: any) => {
    setEditForm(f => ({ ...f, [key]: value }));
    setEditErrors(e => ({ ...e, [key as string]: '' }));
  };
  const handleEditSidebarSave = () => {
    if (!validateEditForm()) return;
    setUsers(users => users.map(u => u.id === editForm.id ? { ...u, ...editForm } as User : u));
    handleEditSidebarClose();
  };

  const columns: TableColumn<User>[] = [
    {
      key: 'id',
      label: 'User ID',
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
    { key: 'username', label: 'Name', filterable: true, required: true, type: 'text' },
    { key: 'role', label: 'Role', filterable: true, required: true, type: 'select', options: USER_ROLES },
  ];

  return (
    <div className="min-h-screen bg-background overflow-hidden lg:pl-64">
      <Sidebar isOpen={true} onClose={() => {}} />
      <header className="sticky top-0 z-30 bg-card border-b">
        <div className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center space-x-3">
            <img src="/logo.svg" alt="ChargePi Logo" className="h-8 w-8" />
            <h1 className="text-xl font-semibold">Users</h1>
          </div>
          <div className="flex items-center space-x-4">
            
          </div>
        </div>
      </header>
      <main className="flex-1 p-6 min-w-0">
        <div className="flex justify-end items-center gap-4 mb-4">
          <ViewToggle view={view} onViewChange={setView} />
          <AddSidebarButton<User>
            buttonTitle="Add User"
            sidebarTitle="Add New User"
            columns={columns}
            onAdd={user => setUsers([...users, user])}
          />
        </div>
        <div className="overflow-x-auto">
          {view === 'list' ? (
            <GenericTable<User>
              columns={columns}
              data={users}
              filterInputClassName="w-1/4"
              className="min-w-full text-sm"
              columnVisibilityEnabled={true}
              columnSearchEnabled={true}
              onRowClick={handleRowClick}
            />
          ) : (
            <GenericCardView<User>
              data={users}
              columns={columns}
              className="min-w-full"
              onEdit={handleRowClick}
              onDelete={user => setShowDelete(user)}
            />
          )}
        </div>
        {editing && (
          <div className="flex gap-2 mt-4">
            <Button size="sm" variant="accent" onClick={handleEditSave}><Check className="h-4 w-4 mr-1" />Save</Button>
            <Button size="sm" variant="outline" onClick={handleEditCancel}><X className="h-4 w-4 mr-1" />Cancel</Button>
          </div>
        )}
        <SettingsSidebar
          open={editSidebarOpen}
          onClose={handleEditSidebarClose}
          title="Edit User"
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
        <Dialog open={!!showDelete} onOpenChange={open => !open && setShowDelete(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Confirm Delete</DialogTitle>
            </DialogHeader>
            <div className="mb-4">Are you sure you want to delete this user?</div>
            <div className="flex gap-2 justify-end">
              <Button variant="outline" onClick={() => setShowDelete(null)}>Cancel</Button>
              <Button variant="destructive" onClick={handleDelete}>Delete</Button>
            </div>
          </DialogContent>
        </Dialog>
      </main>
    </div>
  );
};

export default Users; 