# Security

## Cryptographic Defaults

- Ed25519 for device signing.
- X25519-compatible key agreement reserved for encrypted provisioning/session setup.
- HKDF-SHA256 for rotating tag ephemeral IDs.
- AES-GCM for encrypted local/backend stored secrets.
- Secure random nonces for sightings and tokens.
- Protocol version, time bucket, nonce, and replay protection in sighting payloads.

## Key Management

- Tag secret is generated during pairing and stored encrypted server-side.
- Merchant stand private key must be stored in secure element or encrypted NVS in production.
- Backend stores stand public keys and verifies signed sightings.
- Admin tokens in local development must be replaced with SSO and scoped service credentials.

## Abuse Prevention

- Rate limiting is enabled at API edge.
- Stand sightings require signatures.
- Deduplication flags repeated nonce/source/time/tag payloads.
- Confidence scoring lowers suspicious, stale, duplicate, or flagged sources.
- Admin actions are audited.
- Abuse reports can disable tags or stands.

## Threat Model Summary

Primary threats are fake sightings, replay, owner/finder deanonymization, exact merchant disclosure, unwanted tracker misuse, account takeover, and malicious or compromised stands. Mitigations are documented in `docs/threat-model.md`.

## Production Hardening

- Add mTLS or device-bound OAuth for stand uploads and heartbeat.
- Store server keys in KMS/HSM.
- Enforce tenant/role permissions on merchant staff access.
- Add fraud pipelines for impossible travel and high-volume anomalies.
- Add formal privacy reviews for every new data field.
