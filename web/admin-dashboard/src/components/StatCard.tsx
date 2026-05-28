import React from 'react';

type Props = {
  icon: React.ReactNode;
  label: string;
  value: number;
};

export function StatCard({ icon, label, value }: Props) {
  return (
    <div className="stat">
      <div className="statIcon">{icon}</div>
      <div>
        <span>{label}</span>
        <strong>{value}</strong>
      </div>
    </div>
  );
}
