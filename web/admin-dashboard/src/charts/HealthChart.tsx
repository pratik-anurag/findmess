import React from 'react';

export function HealthChart({ stands }: { stands: Record<string, unknown>[] }) {
  const online = stands.filter((stand) => stand.status === 'online').length;
  const offline = Math.max(stands.length - online, 0);
  const total = Math.max(stands.length, 1);
  return (
    <section className="health">
      <div>
        <h2>Device health</h2>
        <p>Stand status is based on heartbeat and reported errors.</p>
      </div>
      <div className="bars">
        <div style={{ width: `${(online / total) * 100}%` }}>online {online}</div>
        <div style={{ width: `${(offline / total) * 100}%` }}>attention {offline}</div>
      </div>
    </section>
  );
}
