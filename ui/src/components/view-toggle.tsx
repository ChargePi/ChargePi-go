import React from 'react';
import { List, Grid3X3 } from 'lucide-react';
import { Button } from './ui/button';

interface ViewToggleProps {
  view: 'list' | 'cards';
  onViewChange: (view: 'list' | 'cards') => void;
}

const ViewToggle: React.FC<ViewToggleProps> = ({ view, onViewChange }) => {
  return (
    <div className="flex items-center gap-1 p-1 bg-muted rounded-lg h-10">
      <Button
        variant={view === 'list' ? 'default' : 'ghost'}
        size="icon"
        onClick={() => onViewChange('list')}
        className="h-10 w-auto px-4"
      >
        <List className="h-4 w-4 mr-1" />
        List
      </Button>
      <Button
        variant={view === 'cards' ? 'default' : 'ghost'}
        size="icon"
        onClick={() => onViewChange('cards')}
        className="h-10 w-auto px-4"
      >
        <Grid3X3 className="h-4 w-4 mr-1" />
        Cards
      </Button>
    </div>
  );
};

export default ViewToggle; 