# Operations Runbook

## Daily Checks

- API health: `/healthz`.
- Prometheus metrics: `/metrics`.
- Stand heartbeat freshness.
- Abuse queue.
- Firmware rollout status.
- Raw sighting retention jobs.

## Incident: Suspicious Stand

1. Search stand in admin dashboard.
2. Review heartbeat, sighting volume, and audit records.
3. Disable stand if needed through abuse action.
4. Dispatch technician for key rotation or replacement.

## Incident: Unwanted Tracker Report

1. Review abuse report.
2. Check tag status and ownership history.
3. Disable tag/account when policy threshold is met.
4. Preserve audit record.

## Incident: Firmware Rollout Failure

1. Pause rollout.
2. Check device heartbeats and last error.
3. Revert manifest to previous staged version.
4. Use field technician flow for devices stuck in OTA state.
