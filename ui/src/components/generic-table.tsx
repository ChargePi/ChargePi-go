import React, { useMemo, useState } from 'react';
import {
  useReactTable,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  flexRender,
  ColumnDef,
  VisibilityState,
  ColumnFiltersState,
} from '@tanstack/react-table';
import { Eye, EyeOff, ChevronDown, Search } from 'lucide-react';

export interface TableColumn<T> {
  key: keyof T;
  label: string;
  filterable?: boolean;
  render?: (value: any, row: T) => React.ReactNode;
  className?: string;
  editable?: boolean; // Add this property
  required?: boolean; // If the field is required for add/edit
  type?: 'text' | 'email' | 'number' | 'date' | 'password' | 'select'; // Input type for validation
  options?: { value: string; label: string }[]; // Options for select type
}

interface GenericTableProps<T extends { id: string | number }> {
  columns: TableColumn<T>[];
  data: T[];
  className?: string;
  filterInputClassName?: string;
  columnVisibilityEnabled?: boolean;
  columnSearchEnabled?: boolean;
  onRowClick?: (row: T) => void;
}

function GenericTable<T extends { id: string | number }>({
  columns,
  data,
  className,
  filterInputClassName,
  columnVisibilityEnabled = false,
  columnSearchEnabled = false,
  onRowClick,
}: GenericTableProps<T>) {
  // Build TanStack column defs
  const columnDefs = useMemo<ColumnDef<T, any>[]>(() =>
    columns.map(col => ({
      accessorKey: col.key as string,
      header: col.label,
      cell: info => col.render ? col.render(info.getValue(), info.row.original) : info.getValue(),
      enableColumnFilter: !!col.filterable,
      meta: { className: col.className },
    })),
    [columns]
  );

  // State for column visibility and filters
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);

  // Table instance
  const table = useReactTable({
    data,
    columns: columnDefs,
    state: {
      columnVisibility,
      columnFilters,
    },
    onColumnVisibilityChange: setColumnVisibility,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    debugTable: false,
  });

  // Dropdown state for columns
  const [showColumnsDropdown, setShowColumnsDropdown] = useState(false);

  return (
    <div className={className ? className + ' w-full' : 'w-full'}>
      <div className="flex items-center mb-4 gap-2">
        {columnVisibilityEnabled && (
          <div className="relative">
            <button
              className={`flex items-center gap-1 px-3 py-2 border border-orange-500 bg-background text-sm rounded-lg shadow-sm hover:bg-orange-50 dark:hover:bg-orange-900/10 transition-colors focus:outline-none focus:ring-2 focus:ring-orange-400/50 group ${showColumnsDropdown ? 'ring-2 ring-orange-400/70' : ''}`}
              onClick={() => setShowColumnsDropdown(v => !v)}
              type="button"
              aria-haspopup="listbox"
              aria-expanded={showColumnsDropdown}
            >
              <ChevronDown className={`h-4 w-4 mr-1 transition-transform ${showColumnsDropdown ? 'rotate-180' : ''}`} />
              <span className="font-medium">Columns</span>
            </button>
            {showColumnsDropdown && (
              <div className="absolute left-0 mt-2 z-50 bg-card border border-orange-200 rounded-lg shadow-xl p-2 min-w-[180px] animate-fade-in">
                <div className="pb-2 mb-2 border-b border-border text-xs font-semibold text-muted-foreground uppercase tracking-wider">Show Columns</div>
                {table.getAllLeafColumns().map(col => (
                  <label key={col.id} className="flex items-center gap-2 py-1 cursor-pointer text-sm hover:bg-muted/40 rounded px-2 transition-colors">
                    <input
                      type="checkbox"
                      checked={col.getIsVisible()}
                      onChange={col.getToggleVisibilityHandler()}
                      className="accent-orange-500"
                    />
                    {col.getIsVisible() ? <Eye className="h-4 w-4 text-green-500" /> : <EyeOff className="h-4 w-4 text-muted-foreground" />}
                    <span>{col.columnDef.header as string}</span>
                  </label>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
      <div className="overflow-x-auto">
        <table className="min-w-full bg-card border border-border rounded-lg">
          <thead>
            {table.getHeaderGroups().map(headerGroup => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map(header => (
                  <th key={header.id} className={`px-4 py-2 text-left font-semibold ${header.column.columnDef.meta && (header.column.columnDef.meta as any).className ? (header.column.columnDef.meta as any).className : ''}`}>
                    {flexRender(header.column.columnDef.header, header.getContext())}
                  </th>
                ))}
              </tr>
            ))}
            {columnSearchEnabled && (
              <tr>
                {table.getAllLeafColumns().map(col => (
                  <th key={col.id} className="px-4 py-1">
                    {col.getCanFilter() ? (
                      <div className="relative flex items-center">
                        <Search className="absolute left-2 h-3.5 w-3.5 text-muted-foreground pointer-events-none" />
                        <input
                          type="text"
                          value={(col.getFilterValue() as string) || ''}
                          onChange={e => col.setFilterValue(e.target.value)}
                          placeholder={`Search ${(col.columnDef.header as string)}`}
                          className="w-full pl-7 pr-2 py-1 border border-border rounded bg-background text-xs focus:ring-2 focus:ring-orange-400/50 focus:outline-none transition-all"
                        />
                      </div>
                    ) : null}
                  </th>
                ))}
              </tr>
            )}
          </thead>
          <tbody>
            {table.getRowModel().rows.map((row, i) => (
              <tr
                key={row.id}
                className={`border-t border-border cursor-pointer ${i % 2 === 1 ? 'bg-muted/40' : ''}`}
                tabIndex={0}
                onClick={() => onRowClick && onRowClick(row.original)}
              >
                {row.getVisibleCells().map(cell => (
                  <td key={cell.id} className={`px-4 py-2 ${cell.column.columnDef.meta && (cell.column.columnDef.meta as any).className ? (cell.column.columnDef.meta as any).className : ''}`}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default GenericTable; 