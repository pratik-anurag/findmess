export type Row = Record<string, unknown>;

export type DashboardData = {
  users: Row[];
  tags: Row[];
  merchants: Row[];
  stands: Row[];
  sightings: Row[];
  recovery: Row[];
  abuse: Row[];
  firmware: Row[];
  audit: Row[];
};

class ApiClient {
  baseUrl = import.meta.env.VITE_FINDMESH_API_BASE_URL ?? 'http://localhost:8080';
  token = 'dev-admin-token';

  async get(path: string): Promise<Row[]> {
    const response = await fetch(`${this.baseUrl}${path}`, {
      headers: { authorization: `Bearer ${this.token}` },
    });
    if (!response.ok) throw new Error(await response.text());
    return response.json();
  }
}

export const apiClient = new ApiClient();

export async function loadDashboardData(): Promise<DashboardData> {
  try {
    const [users, tags, merchants, stands, sightings, recovery, abuse, audit] = await Promise.all([
      apiClient.get('/v1/admin/users'),
      apiClient.get('/v1/admin/tags'),
      apiClient.get('/v1/admin/merchants'),
      apiClient.get('/v1/admin/stands'),
      apiClient.get('/v1/admin/sightings'),
      apiClient.get('/v1/recovery/requests'),
      apiClient.get('/v1/admin/abuse/reports'),
      apiClient.get('/v1/admin/audit-events'),
    ]);
    return { users, tags, merchants, stands, sightings, recovery, abuse, firmware: [], audit };
  } catch {
    return mockDashboardData;
  }
}

const now = new Date().toISOString();

const mockDashboardData: DashboardData = {
  users: [{ id: 'demo-user', status: 'active', created_at: now }],
  tags: [{ id: 'demo-tag', status: 'lost', public_label: 'Keys', firmware_version: 'tag-dev', last_seen_at: now }],
  merchants: [{ id: 'demo-merchant', display_name: 'Demo Store', status: 'verified', city: 'Bengaluru', recovery_enabled: true }],
  stands: [{ id: 'demo-stand', status: 'online', firmware_version: 'stand-dev', power_source: 'usb_c', wifi_status: 'connected', last_heartbeat_at: now }],
  sightings: [{ id: 'demo-sighting', source_type: 'merchant_stand', zone_id: 'demo-zone', rssi_bucket: 'near', confidence_score: 90, suspicious: false, created_at: now }],
  recovery: [{ id: 'demo-recovery', status: 'requested', merchant_id: 'demo-merchant', zone_id: 'demo-zone', created_at: now }],
  abuse: [{ id: 'demo-abuse', category: 'unknown_tracker_alert', status: 'open', tag_id: 'demo-tag', stand_id: '', created_at: now }],
  firmware: [{ id: 'demo-fw', device_type: 'merchant_stand', version: '0.1.0', rollout_status: 'staged', created_at: now }],
  audit: [{ id: 'demo-audit', actor_type: 'admin', actor_id: 'admin', action: 'list_users', target_type: 'user', created_at: now }],
};
