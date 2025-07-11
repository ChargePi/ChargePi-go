import React, { useState } from 'react';
import { Copy, User, Trash2, Pencil } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
  DialogClose
} from '@/components/ui/dialog';

interface CardViewProps<T> {
  data: T[];
  columns: Array<{
    key: keyof T;
    label: string;
    render?: (value: any, row: T) => React.ReactNode;
  }>;
  className?: string;
  onEdit?: (item: T) => void;
  onDelete?: (item: T) => void;
}

function GenericCardView<T extends { id: string | number }>({
  data,
  columns,
  className,
  onEdit,
  onDelete,
}: CardViewProps<T>) {
  const [deleteModalItem, setDeleteModalItem] = useState<T | null>(null);

  const renderValue = (item: T, col: { key: keyof T; label: string; render?: (value: any, row: T) => React.ReactNode }) => {
    if (col.render) {
      return col.render(item[col.key], item);
    }
    return String(item[col.key] || '');
  };

  return (
    <>
      <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 ${className || ''}`}>
        {data.map((item) => (
          <div
            key={item.id}
            className="bg-card border border-border rounded-lg p-4 hover:shadow-md transition-shadow flex items-start gap-4 relative"
          >
            {/* Top right action buttons */}
            {(onEdit || onDelete) && (
              <div className="absolute top-3 right-3 flex gap-2 z-10">
                {onEdit && (
                  <button
                    className="p-2 w-9 h-9 rounded-md hover:bg-muted transition-colors flex items-center justify-center"
                    aria-label="Edit"
                    onClick={() => onEdit(item)}
                  >
                    <Pencil className="h-4 w-4 text-accent-orange" />
                  </button>
                )}
                {onDelete && (
                  <button
                    className="p-2 w-9 h-9 rounded-md hover:bg-red-100 dark:hover:bg-red-900/20 transition-colors flex items-center justify-center"
                    aria-label="Delete"
                    onClick={() => setDeleteModalItem(item)}
                  >
                    <Trash2 className="h-4 w-4 text-red-500" />
                  </button>
                )}
              </div>
            )}
            <div className="flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-lg bg-accent-orange/10 mt-1">
              <User className="h-7 w-7 text-accent-orange" />
            </div>
            <div className="flex-1 space-y-3">
              {columns.map((col) => (
                <div key={String(col.key)} className="flex flex-col items-start">
                  <span className="text-sm font-medium text-accent-orange capitalize min-w-[90px] text-left">
                    {col.label}:
                  </span>
                  <span className="text-sm ml-0 flex-1 text-white text-left">
                    {renderValue(item, col)}
                  </span>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
      {/* Single Delete confirmation modal rendered outside the card map */}
      {onDelete && deleteModalItem && (
        <Dialog open={!!deleteModalItem} onOpenChange={open => !open && setDeleteModalItem(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Delete Item</DialogTitle>
              <DialogDescription>Are you sure you want to delete this item?</DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <DialogClose asChild>
                <button className="px-4 py-2 rounded bg-muted text-foreground hover:bg-muted/80 focus:outline-none focus:ring-2 focus:ring-primary">Cancel</button>
              </DialogClose>
              <button
                className="px-4 py-2 rounded bg-red-600 text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-400"
                onClick={() => {
                  setDeleteModalItem(null);
                  onDelete(deleteModalItem);
                }}
              >
                Delete
              </button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}

export default GenericCardView; 