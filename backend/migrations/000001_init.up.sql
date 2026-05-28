CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  phone_hash text UNIQUE NOT NULL,
  phone_encrypted text,
  status text NOT NULL DEFAULT 'active',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE user_devices (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  platform text NOT NULL,
  push_token text,
  app_version text NOT NULL DEFAULT '',
  finder_participation_enabled boolean NOT NULL DEFAULT false,
  last_seen_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE tags (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  serial_hash text UNIQUE NOT NULL,
  owner_user_id uuid REFERENCES users(id),
  status text NOT NULL CHECK (status IN ('unpaired','active','lost','disabled','recovered')),
  public_label text,
  tag_secret_encrypted text,
  battery_level int,
  firmware_version text,
  last_seen_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE tag_ownership_events (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tag_id uuid NOT NULL REFERENCES tags(id),
  user_id uuid NOT NULL REFERENCES users(id),
  event_type text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE merchants (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  display_name text NOT NULL,
  status text NOT NULL DEFAULT 'pending_verification',
  city text,
  category text,
  recovery_enabled boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE merchant_zones (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL REFERENCES merchants(id),
  coarse_geohash text NOT NULL,
  display_area text NOT NULL,
  latitude numeric,
  longitude numeric,
  location_precision_meters int NOT NULL DEFAULT 500,
  public_visibility text NOT NULL DEFAULT 'coarse_only',
  geom geography(Point, 4326),
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE stands (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid REFERENCES merchants(id),
  zone_id uuid REFERENCES merchant_zones(id),
  serial_hash text UNIQUE NOT NULL,
  public_key text NOT NULL DEFAULT '',
  status text NOT NULL DEFAULT 'unclaimed',
  firmware_version text,
  battery_level int,
  power_source text,
  wifi_status text,
  last_heartbeat_at timestamptz,
  last_error text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE stand_claim_tokens (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  stand_id uuid NOT NULL REFERENCES stands(id),
  token_hash text NOT NULL,
  expires_at timestamptz NOT NULL,
  claimed_at timestamptz
);

CREATE TABLE sightings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tag_ephemeral_id bytea NOT NULL,
  source_type text NOT NULL CHECK (source_type IN ('merchant_stand','user_app')),
  source_id uuid,
  zone_id uuid REFERENCES merchant_zones(id),
  time_bucket timestamptz NOT NULL,
  rssi_bucket text NOT NULL CHECK (rssi_bucket IN ('near','medium','far')),
  confidence_score int NOT NULL,
  nonce text NOT NULL,
  signature text,
  raw_payload_hash text NOT NULL,
  suspicious boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX sightings_dedup_idx ON sightings(source_type, source_id, nonce, tag_ephemeral_id, time_bucket);
CREATE INDEX sightings_tag_time_idx ON sightings(tag_ephemeral_id, time_bucket DESC);

CREATE TABLE lost_mode_sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tag_id uuid NOT NULL REFERENCES tags(id),
  owner_user_id uuid NOT NULL REFERENCES users(id),
  status text NOT NULL,
  safe_message text,
  public_lost_token text UNIQUE,
  created_at timestamptz NOT NULL DEFAULT now(),
  resolved_at timestamptz
);

CREATE TABLE last_seen_summaries (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tag_id uuid NOT NULL REFERENCES tags(id),
  lost_mode_session_id uuid REFERENCES lost_mode_sessions(id),
  zone_id uuid REFERENCES merchant_zones(id),
  display_area text NOT NULL,
  confidence_level text NOT NULL,
  confidence_score int NOT NULL,
  last_seen_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE recovery_requests (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  lost_mode_session_id uuid NOT NULL REFERENCES lost_mode_sessions(id),
  merchant_id uuid REFERENCES merchants(id),
  zone_id uuid REFERENCES merchant_zones(id),
  status text NOT NULL,
  masked_thread_id text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE abuse_reports (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  reporter_user_id uuid REFERENCES users(id),
  tag_id uuid REFERENCES tags(id),
  stand_id uuid REFERENCES stands(id),
  merchant_id uuid REFERENCES merchants(id),
  category text NOT NULL,
  description text NOT NULL DEFAULT '',
  status text NOT NULL DEFAULT 'open',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE audit_events (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  actor_type text NOT NULL,
  actor_id uuid,
  action text NOT NULL,
  target_type text NOT NULL,
  target_id uuid,
  metadata jsonb NOT NULL DEFAULT '{}',
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE firmware_releases (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  device_type text NOT NULL,
  version text NOT NULL,
  manifest_url text NOT NULL,
  binary_url text NOT NULL,
  signature text NOT NULL,
  rollout_status text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE device_heartbeats (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  stand_id uuid NOT NULL REFERENCES stands(id),
  firmware_version text NOT NULL,
  battery_level int,
  power_source text NOT NULL,
  wifi_rssi int NOT NULL,
  buffer_count int NOT NULL,
  uptime_seconds bigint NOT NULL,
  last_error text,
  created_at timestamptz NOT NULL DEFAULT now()
);
