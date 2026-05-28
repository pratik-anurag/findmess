# Privacy

FindMesh uses data minimization as a product requirement.

## Never Included In Sightings

- Phone number.
- Raw tag serial.
- Owner user ID.
- Finder identity visible to owner.
- Owner identity visible to finder.
- UPI ID or payment transaction data.
- POS data.

## Stored Data

- Phone hash and encrypted phone value for authentication.
- Tag ownership and encrypted tag secret.
- Merchant and stand operational records.
- Coarse merchant zone metadata.
- Raw sightings with ephemeral ID, source type, internal source ID, time bucket, RSSI bucket, nonce, and confidence.
- Aggregated last-seen summaries.
- Audit and abuse workflow records.

## Retention

- Raw sightings: 30 days by default.
- Lost-mode relevant sightings: active session plus 30 days.
- Aggregated last-seen summaries: retained longer for recovery history.
- Audit and abuse logs: configurable per policy.

## User Controls

- Disable anonymous finder participation.
- Delete account.
- Delete tag.
- Disable lost mode.
- Export basic account data.
- Report unwanted tracker or safety concerns.

## Display Rules

Owners see coarse areas and confidence labels. Exact merchant identity requires recovery flow and merchant opt-in. No live trail UI is provided.
