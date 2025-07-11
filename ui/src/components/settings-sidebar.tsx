import React, { useState } from 'react';

interface SettingsSidebarProps {
  open: boolean;
  onClose: () => void;
  onSave?: () => void;
  title: string;
  children: React.ReactNode;
  actions?: React.ReactNode;
}

const SettingsSidebar: React.FC<SettingsSidebarProps> = ({ open, onClose, onSave, title, children, actions }) => {
  const [isClosing, setIsClosing] = useState(false);

  const handleRequestClose = () => {
    setIsClosing(true);
    setTimeout(() => {
      setIsClosing(false);
      onClose();
    }, 250);
  };

  React.useEffect(() => {
    if (!open) return;
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        e.preventDefault();
        handleRequestClose();
      } else if (e.key === 'Enter' && onSave) {
        // Only trigger save if not inside a textarea or button
        if (document.activeElement && ['TEXTAREA', 'BUTTON'].indexOf(document.activeElement.tagName) !== -1) return;
        e.preventDefault();
        onSave();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [open, onSave]);

  if (!open && !isClosing) return null;
  return (
    <>
      <div
        className="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm"
        onClick={handleRequestClose}
      />
      <div className={`fixed right-0 top-0 h-full w-96 z-50 bg-card shadow-2xl border-l flex flex-col justify-between p-6 animate-in slide-in-from-right-1/2 transition-transform duration-200 ${isClosing ? 'animate-slide-out-right' : ''}`}>
        <div>
          <h3 className="text-lg font-semibold mb-4">{title}</h3>
          {children}
        </div>
        {actions && (
          <div className="flex gap-2 mt-6">{actions}</div>
        )}
      </div>
    </>
  );
};

export default SettingsSidebar; 