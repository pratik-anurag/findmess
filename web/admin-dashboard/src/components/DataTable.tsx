import React from 'react';

type Props = {
  rows: Record<string, unknown>[];
  columns: string[];
};

export function DataTable({ rows, columns }: Props) {
  return (
    <div className="tableWrap">
      <table>
        <thead>
          <tr>
            {columns.map((column) => <th key={column}>{column.replaceAll('_', ' ')}</th>)}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, index) => (
            <tr key={`${row.id ?? index}`}>
              {columns.map((column) => <td key={column}>{format(row[column])}</td>)}
            </tr>
          ))}
          {rows.length === 0 && (
            <tr>
              <td colSpan={columns.length}>No records</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

function format(value: unknown) {
  if (value === undefined || value === null || value === '') return ' ';
  if (typeof value === 'boolean') return value ? 'yes' : 'no';
  if (typeof value === 'object') return JSON.stringify(value);
  return String(value);
}
