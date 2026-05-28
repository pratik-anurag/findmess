# Architecture

FindMesh uses four cooperating surfaces:

- Mobile app: owner, finder, merchant, provisioning, anti-stalking, and privacy controls.
- Merchant stand firmware: signed zone witness with BLE scan/advertise, NFC setup, Wi-Fi upload, heartbeat, offline buffer, and OTA.
- Lost tag firmware: rotating BLE ephemeral IDs, NFC lost-mode URL, pairing, buzzer, battery, anti-stalking sound behavior.
- Backend: auth, tag ownership, lost mode, sightings, confidence scoring, recovery requests, abuse review, audit, retention, firmware manifests, and admin APIs.

Local development uses an in-memory repository and Docker Compose dependencies. Production should replace the repository with PostgreSQL/PostGIS implementations, Redis-backed rate limiting, NATS event publishing, device-bound auth, and managed observability.

Sightings are processed as append-only raw records, scored, matched against active lost tags, and reduced into coarse last-seen summaries.
