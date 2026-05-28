import React, { useEffect, useMemo, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { Activity, AlertTriangle, Cpu, Radio, Search, Shield, Store, Tags, Users } from 'lucide-react';
import { apiClient, DashboardData, loadDashboardData } from './api/client';
import { DataTable } from './components/DataTable';
import { StatCard } from './components/StatCard';
import { HealthChart } from './charts/HealthChart';
import './styles.css';

const tabs = ['Overview', 'Users', 'Tags', 'Merchants', 'Stands', 'Sightings', 'Recovery', 'Abuse', 'Firmware', 'Audit'];

function App() {
  const [active, setActive] = useState('Overview');
  const [data, setData] = useState<DashboardData | null>(null);
  const [token, setToken] = useState(localStorage.getItem('findmesh_admin_token') ?? 'dev-admin-token');

  useEffect(() => {
    localStorage.setItem('findmesh_admin_token', token);
    apiClient.token = token;
    loadDashboardData().then(setData);
  }, [token]);

  const content = useMemo(() => {
    if (!data) return <div className="empty">Loading dashboard...</div>;
    switch (active) {
      case 'Users':
        return <DataTable rows={data.users} columns={['id', 'status', 'created_at']} />;
      case 'Tags':
        return <DataTable rows={data.tags} columns={['id', 'status', 'public_label', 'firmware_version', 'last_seen_at']} />;
      case 'Merchants':
        return <DataTable rows={data.merchants} columns={['id', 'display_name', 'status', 'city', 'recovery_enabled']} />;
      case 'Stands':
        return <DataTable rows={data.stands} columns={['id', 'status', 'firmware_version', 'power_source', 'wifi_status', 'last_heartbeat_at']} />;
      case 'Sightings':
        return <DataTable rows={data.sightings} columns={['id', 'source_type', 'zone_id', 'rssi_bucket', 'confidence_score', 'suspicious', 'created_at']} />;
      case 'Recovery':
        return <DataTable rows={data.recovery} columns={['id', 'status', 'merchant_id', 'zone_id', 'created_at']} />;
      case 'Abuse':
        return <DataTable rows={data.abuse} columns={['id', 'category', 'status', 'tag_id', 'stand_id', 'created_at']} />;
      case 'Firmware':
        return <DataTable rows={data.firmware} columns={['id', 'device_type', 'version', 'rollout_status', 'created_at']} />;
      case 'Audit':
        return <DataTable rows={data.audit} columns={['id', 'actor_type', 'actor_id', 'action', 'target_type', 'created_at']} />;
      default:
        return <Overview data={data} />;
    }
  }, [active, data]);

  return (
    <div className="app">
      <aside>
        <div className="brand"><Shield size={22} /> FindMesh Admin</div>
        <label className="token">
          Admin token
          <input value={token} onChange={(event) => setToken(event.target.value)} />
        </label>
        <nav>
          {tabs.map((tab) => (
            <button key={tab} className={tab === active ? 'active' : ''} onClick={() => setActive(tab)}>
              {tab}
            </button>
          ))}
        </nav>
      </aside>
      <main>
        <header>
          <div>
            <h1>{active}</h1>
            <p>Audited operational console for private lost-item recovery.</p>
          </div>
          <button onClick={() => loadDashboardData().then(setData)}><Search size={16} /> Refresh</button>
        </header>
        {content}
      </main>
    </div>
  );
}

function Overview({ data }: { data: DashboardData }) {
  return (
    <div className="overview">
      <section className="stats">
        <StatCard icon={<Users />} label="Users" value={data.users.length} />
        <StatCard icon={<Tags />} label="Tags" value={data.tags.length} />
        <StatCard icon={<Store />} label="Merchants" value={data.merchants.length} />
        <StatCard icon={<Radio />} label="Stands" value={data.stands.length} />
        <StatCard icon={<Activity />} label="Sightings" value={data.sightings.length} />
        <StatCard icon={<AlertTriangle />} label="Open abuse" value={data.abuse.filter((r) => r.status === 'open').length} />
        <StatCard icon={<Cpu />} label="Firmware" value={data.firmware.length} />
      </section>
      <HealthChart stands={data.stands} />
      <DataTable rows={data.sightings.slice(0, 8)} columns={['source_type', 'rssi_bucket', 'confidence_score', 'suspicious', 'created_at']} />
    </div>
  );
}

createRoot(document.getElementById('root')!).render(<App />);
