# Data Model

The PostgreSQL schema is in `backend/migrations/000001_init.up.sql`.

Important entities:

- `users`: phone hash, encrypted phone, lifecycle status.
- `user_devices`: app installations, push token, finder participation.
- `tags`: serial hash, encrypted tag secret, owner, lifecycle status.
- `merchants`: merchant profile and recovery opt-in.
- `merchant_zones`: coarse display area and optional PostGIS point with precision.
- `stands`: stand public key, zone binding, firmware, health status.
- `sightings`: ephemeral ID, source, time/RSSI bucket, nonce, signature, confidence, suspicious flag.
- `lost_mode_sessions`: owner-controlled active recovery state and public found token.
- `last_seen_summaries`: privacy-preserving display record for owner.
- `recovery_requests`: masked owner/merchant/finder flow.
- `abuse_reports`: safety review workflow.
- `audit_events`: sensitive admin/support action log.
- `firmware_releases` and `device_heartbeats`: fleet operations.

Raw phone numbers, payment details, POS data, UPI transaction data, and finder identities visible to owners are not modeled.
