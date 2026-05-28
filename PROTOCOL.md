# FindMesh Protocol

## Device Types

- `lost_tag`: BLE/NFC lost-item tag with rotating ephemeral IDs.
- `merchant_stand`: BLE scanner/advertiser, NFC setup point, Wi-Fi uploader, signed sighting source.
- `user_app`: mobile finder node and owner/finder interface.

## BLE Advertisements

Lost tag:

```json
{
  "adv_type": "FM_TAG",
  "protocol_version": 1,
  "ephemeral_id": "16 bytes",
  "flags": {
    "lost_hint": false,
    "battery_low": false
  }
}
```

Merchant stand zone:

```json
{
  "adv_type": "FM_ZONE",
  "protocol_version": 1,
  "zone_ephemeral_id": "16 bytes",
  "stand_capabilities": ["scan", "wifi", "nfc"]
}
```

## NFC Payloads

- Lost tag NFC: `findmesh://tag-found?t=<public_lost_token>`
- Merchant stand setup: `findmesh://stand/setup?s=<stand_claim_token>`
- Merchant zone recovery: `findmesh://merchant-zone?z=<zone_public_token>`

## Ephemeral ID Derivation

Each tag has a 32-byte `tag_secret`. The epoch is `floor(unix_seconds / 900)`.

```text
ephemeral_id = Truncate128(HKDF-SHA256(tag_secret, "findmesh-ephemeral-id" || uint64_be(epoch)))
```

Test vector:

- secret: `0123456789abcdef0123456789abcdef`
- epoch: `123456`
- ephemeral_id: `9e40ef0c677adae9870809e1cd952fc2`

## Sighting Payload

```json
{
  "protocol_version": 1,
  "source_type": "merchant_stand",
  "source_id": "internal stand id",
  "tag_ephemeral_id": "00112233445566778899aabbccddeeff",
  "zone_ephemeral_id": "optional 16 bytes",
  "zone_id": "internal zone id",
  "time_bucket": "2026-05-28T10:15:00Z",
  "rssi_bucket": "near",
  "nonce": "random",
  "signature": "ed25519 signature for stands"
}
```

## Signing Rules

Merchant stands sign a canonical string containing protocol version, source type, source ID, tag ephemeral ID, optional zone ephemeral ID, zone ID, time bucket, RSSI bucket, and nonce. Backend verifies the signature against the stand public key before storing the sighting.

## Deduplication

Dedup key:

```text
source_type | source_id | nonce | tag_ephemeral_id | time_bucket
```

Repeated keys are treated as duplicate or replay-suspicious and lower confidence.

## Confidence Scoring

- `+50` merchant stand sighting.
- `+30` user app sighting.
- `+20` two independent sources in same zone/time.
- `+10` RSSI near.
- `+10` high reputation or healthy stand.
- `-20` older than one hour.
- `-30` duplicate/replay suspected.
- `-50` flagged source.

Labels: `low` 0-29, `medium` 30-59, `high` 60-84, `very_high` 85+.

## Lost Mode Matching

The backend derives expected ephemeral IDs for active lost tags around the current time bucket. Matching sightings update a coarse last-seen summary and notify the owner without exposing finder identity.

## Privacy Guarantees

Sightings do not contain phone numbers, owner IDs, raw serials, UPI IDs, POS data, or payment transactions. Owners see coarse areas and confidence labels, not live movement trails.

## Known Limitations

The local repository includes firmware signing and BLE/NFC platform stubs. Production must replace those boundaries with certified hardware drivers, secure storage, and platform-specific background scanning behavior.
