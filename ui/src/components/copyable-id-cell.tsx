import React, { useState } from 'react';
import { Copy } from 'lucide-react';

interface CopyableIdCellProps {
  id: string;
  tooltip?: string;
}

const CopyableIdCell: React.FC<CopyableIdCellProps> = ({ id, tooltip = 'Copy ID' }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = (e: React.MouseEvent) => {
    e.stopPropagation();
    navigator.clipboard.writeText(id);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <span className="flex items-center gap-2 relative">
      <span>{id}</span>
      <button
        className="p-1 rounded hover:bg-muted relative"
        onClick={handleCopy}
        title={tooltip}
      >
        <Copy className="h-4 w-4" />
        {copied && (
          <span className="absolute left-1/2 -translate-x-1/2 top-8 bg-black text-white text-xs rounded px-2 py-1 shadow z-10 animate-fade-in">
            Copied!
          </span>
        )}
      </button>
    </span>
  );
};

export default CopyableIdCell; 