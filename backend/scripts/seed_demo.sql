INSERT INTO users (id, phone_hash, phone_encrypted, status)
VALUES ('00000000-0000-4000-8000-000000000101', 'demo-phone-hash', 'encrypted-demo-phone', 'active')
ON CONFLICT (id) DO NOTHING;

INSERT INTO merchants (id, name, display_name, status, city, category, recovery_enabled)
VALUES ('00000000-0000-4000-8000-000000000201', 'Demo Corner Store', 'Demo Store', 'verified', 'Bengaluru', 'retail', true)
ON CONFLICT (id) DO NOTHING;

INSERT INTO merchant_zones (id, merchant_id, coarse_geohash, display_area, location_precision_meters, public_visibility)
VALUES ('00000000-0000-4000-8000-000000000202', '00000000-0000-4000-8000-000000000201', 'tdr1w', 'near a participating merchant zone in Indiranagar', 500, 'coarse_only')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tags (id, serial_hash, owner_user_id, status, public_label, firmware_version)
VALUES ('00000000-0000-4000-8000-000000000301', 'demo-tag-serial-hash', '00000000-0000-4000-8000-000000000101', 'lost', 'Keys', 'tag-dev')
ON CONFLICT (id) DO NOTHING;

INSERT INTO stands (id, merchant_id, zone_id, serial_hash, public_key, status, firmware_version, power_source, wifi_status, last_heartbeat_at)
VALUES ('00000000-0000-4000-8000-000000000401', '00000000-0000-4000-8000-000000000201', '00000000-0000-4000-8000-000000000202', 'demo-stand-serial-hash', '', 'online', 'stand-dev', 'usb_c', 'connected', now())
ON CONFLICT (id) DO NOTHING;

INSERT INTO lost_mode_sessions (id, tag_id, owner_user_id, status, safe_message, public_lost_token)
VALUES ('00000000-0000-4000-8000-000000000501', '00000000-0000-4000-8000-000000000301', '00000000-0000-4000-8000-000000000101', 'active', 'If found, contact me via FindMesh.', 'demo-lost-token')
ON CONFLICT (id) DO NOTHING;

INSERT INTO firmware_releases (id, device_type, version, manifest_url, binary_url, signature, rollout_status)
VALUES ('00000000-0000-4000-8000-000000000601', 'merchant_stand', '0.1.0', 'https://example.invalid/findmesh/merchant-stand/0.1.0.json', 'https://example.invalid/findmesh/merchant-stand/0.1.0.bin', 'dev-signature', 'staged')
ON CONFLICT (id) DO NOTHING;
