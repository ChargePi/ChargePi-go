import React, { useState, useEffect, useRef } from 'react';
import {
  Command,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem
} from 'cmdk';
import { useNavigate } from 'react-router-dom';

const actions = [
  { label: 'Go to Users', path: '/users' },
  { label: 'Go to RFID Tags', path: '/rfid-tags' },
  { label: 'Open Settings', path: '/settings' },
  { label: 'Go to Dashboard', path: '/dashboard' },
];

const CommandPalette: React.FC = () => {
  const [open, setOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault();
        setOpen((o) => !o);
      }
      if (e.key === 'Escape') setOpen(false);
    };
    window.addEventListener('keydown', down);
    return () => window.removeEventListener('keydown', down);
  }, []);

  useEffect(() => {
    if (open && inputRef.current) inputRef.current.focus();
  }, [open]);

  return (
    open && (
      <div className="fixed top-1/4 left-1/2 -translate-x-1/2 z-[100] w-full max-w-lg rounded-xl bg-card border border-border shadow-2xl p-0 text-white dark:bg-card">
        <Command label="Command Palette">
          <CommandInput
            ref={inputRef}
            placeholder="Type a command or search..."
            className="w-full px-4 py-3 bg-background text-lg rounded-t-xl border-b border-border outline-none"
          />
          <CommandList className="max-h-72 overflow-y-auto">
            <CommandEmpty className="p-4 text-muted-foreground">No results found.</CommandEmpty>
            <CommandGroup heading="Navigation" className="">
              {actions.map((action) => (
                <CommandItem
                  key={action.label}
                  onSelect={() => {
                    setOpen(false);
                    navigate(action.path);
                  }}
                  className="px-4 py-3 cursor-pointer hover:bg-accent-orange/10 transition-colors flex items-center gap-2"
                >
                  {action.label}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </div>
    )
  );
};

export default CommandPalette; 