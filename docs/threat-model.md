# Threat Model

## Assets

- Tag secrets.
- Stand private keys.
- Phone hashes/encrypted phones.
- Sighting data.
- Lost-mode safe messages.
- Admin audit logs.

## Threats

- Fake sightings.
- Replay attacks.
- Compromised merchant stand.
- Finder or owner deanonymization.
- Exact merchant identity disclosure.
- Unwanted tracker misuse.
- Admin overreach.
- Account takeover.
- Firmware tampering.

## Mitigations

- Signed stand sightings and nonce deduplication.
- Authenticated user app sightings.
- Coarse time/location buckets.
- Hidden exact merchant identity by default.
- Audit logging for admin actions.
- Abuse review and device disable actions.
- Firmware manifests with signatures.
- Retention worker for raw sightings.

## Residual Risk

Rotating IDs make unwanted tracker detection harder. Mobile anti-stalking is a local heuristic and should be integrated with platform standards and regulatory guidance before broad launch.
