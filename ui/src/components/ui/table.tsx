import * as React from 'react';

interface Column<T> {
  key: keyof T;
  label: string;
  filterable?: boolean;
  render?: (value: any, row: T) => React.ReactNode;
  className?: string;
}

interface TableProps<T> {
  columns: Column<T>[];
  data: T[];
  filter?: string;
  onFilterChange?: (value: string) => void;
  selectedRow?: T | null;
  onRowSelect?: (row: T | null) => void;
  rowActions?: (row: T) => React.ReactNode;
}

export function Table<T extends { id: string | number }>({
  columns,
  data,
  filter = '',
  onFilterChange,
  selectedRow,
  onRowSelect,
  rowActions,
}: TableProps<T>) {
  const [internalFilter, setInternalFilter] = React.useState(filter);

  React.useEffect(() => {
    setInternalFilter(filter);
  }, [filter]);

  const filteredData = React.useMemo(() => {
    if (!internalFilter) return data;
    return data.filter(row =>
      columns.some(col =>
        col.filterable && String(row[col.key]).toLowerCase().includes(internalFilter.toLowerCase())
      )
    );
  }, [internalFilter, data, columns]);

  return (
    <div className="w-full">
      {onFilterChange && (
        <div className="mb-4">
          <input
            type="text"
            value={internalFilter}
            onChange={e => {
              setInternalFilter(e.target.value);
              onFilterChange(e.target.value);
            }}
            placeholder="Filter..."
            className="w-full px-3 py-2 border border-border rounded-md bg-background text-sm"
          />
        </div>
      )}
      <div className="overflow-x-auto">
        <table className="min-w-full bg-card border border-border rounded-lg">
          <thead>
            <tr>
              {columns.map(col => (
                <th key={String(col.key)} className={`px-4 py-2 text-left font-semibold ${col.className || ''}`}>{col.label}</th>
              ))}
              {rowActions && <th className="px-4 py-2 text-left">Actions</th>}
            </tr>
          </thead>
          <tbody>
            {filteredData.map(row => (
              <tr
                key={row.id}
                className={`border-t border-border cursor-pointer transition-colors ${selectedRow && selectedRow.id === row.id ? 'bg-muted' : 'hover:bg-muted/50'}`}
                onClick={() => onRowSelect && onRowSelect(row)}
              >
                {columns.map(col => (
                  <td key={String(col.key)} className={`px-4 py-2 ${col.className || ''}`}>{col.render ? col.render(row[col.key], row) : String(row[col.key] || '')}</td>
                ))}
                {rowActions && <td className="px-4 py-2">{rowActions(row)}</td>}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
} 