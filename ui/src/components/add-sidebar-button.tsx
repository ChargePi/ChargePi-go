import React, { useState } from 'react';
import { Plus } from 'lucide-react';
import { Button } from './ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';
import SettingsSidebar from './settings-sidebar';
import { v4 as uuidv4 } from 'uuid';
import type { TableColumn } from './generic-table';

interface AddSidebarButtonProps<T> {
  buttonTitle?: string;
  sidebarTitle: string;
  columns: TableColumn<T>[];
  onAdd: (entry: T) => void;
  initial?: Partial<T>;
}

const AddSidebarButton = <T extends { id: string | number }>({
  buttonTitle = 'Add',
  sidebarTitle,
  columns,
  onAdd,
  initial = {},
}: AddSidebarButtonProps<T>) => {
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState<Partial<T>>(initial);
  const [isClosing, setIsClosing] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

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

  const validateForm = () => {
    const newErrors: Record<string, string> = {};
    columns.forEach(col => {
      if (col.editable === false) return;
      const value = form[col.key];
      const error = validateField(col, value);
      if (error) newErrors[col.key as string] = error;
    });
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleOpen = () => {
    setOpen(true);
    setForm(initial);
    setErrors({});
  };
  const handleRequestClose = () => {
    setIsClosing(true);
    setTimeout(() => {
      setIsClosing(false);
      setOpen(false);
    }, 250);
  };
  const handleChange = (key: keyof T, value: any) => {
    setForm(f => ({ ...f, [key]: value }));
    setErrors(e => ({ ...e, [key as string]: '' }));
  };
  const handleSave = () => {
    if (!validateForm()) return;
    const entry = { ...form, id: uuidv4() } as T;
    onAdd(entry);
    handleRequestClose();
  };

  return (
    <>
      <Button
        variant="outline"
        aria-label={buttonTitle}
        title={buttonTitle}
        size="icon"
        className="border-2 border-orange-500 hover:bg-orange-500 hover:border-orange-500 hover:text-white transition-colors"
        onClick={handleOpen}
      >
        <Plus className="h-5 w-5" />
      </Button>
      <SettingsSidebar
        open={open || isClosing}
        onClose={handleRequestClose}
        title={sidebarTitle}
        actions={
          <>
            <Button variant="outline" className="flex-1" onClick={handleRequestClose}>
              Cancel
            </Button>
            <Button variant="accent" className="flex-1" onClick={handleSave}>
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
                <Select value={form[col.key] as string || ''} onValueChange={value => handleChange(col.key, value)}>
                  <SelectTrigger className={errors[col.key as string] ? 'border-red-500' : ''}>
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
                  className={`w-full px-3 py-2 border rounded bg-background ${errors[col.key as string] ? 'border-red-500' : 'border'} focus:outline-none`}
                  value={form[col.key] as string || ''}
                  onChange={e => handleChange(col.key, e.target.value)}
                  placeholder={`Enter ${col.label}`}
                  required={col.required}
                />
              )}
              {errors[col.key as string] && (
                <p className="text-xs text-red-600 mt-1">{errors[col.key as string]}</p>
              )}
            </div>
          ))}
        </form>
      </SettingsSidebar>
    </>
  );
};

export default AddSidebarButton; 