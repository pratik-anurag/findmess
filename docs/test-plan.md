# Test Plan

Backend:

- Auth OTP start/verify/logout/delete.
- Tag pairing and ownership conflicts.
- Lost mode open/resolve.
- Sighting ingestion and batch ingestion.
- Stand signature validation.
- Deduplication and replay marking.
- Confidence scoring thresholds.
- Last-seen summary creation.
- Recovery request lifecycle.
- Abuse report lifecycle and disable actions.
- Admin audit logging.
- Retention worker.

Mobile:

- Repository and API client retries.
- Secure token storage.
- Tag list/detail/lost-mode UI logic.
- Nearby finder signal logic.
- Anti-stalking observation heuristic.
- Stand provisioning with BLE/NFC mocks.

Firmware:

- Ephemeral ID test vector.
- RSSI bucketing.
- Sighting buffer ring behavior.
- Signing payload formation.
- State transitions for stand and tag.
- OTA manifest verification path.
