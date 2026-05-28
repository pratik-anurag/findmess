# Mobile

The Flutter app is in `mobile/findmesh_app`.

Architecture:

- `core`: API client, app state, models.
- `platform`: BLE, NFC, notifications, location, secure storage abstractions.
- `features`: auth, tags, lost mode, nearby find, finder network, merchant, stand provisioning, recovery, privacy, anti-stalking, settings, debug.

Default state management: Riverpod.

Android is the first BLE background scanning target. iOS support is intentionally abstracted because background BLE scanning/advertising behavior depends on platform policies and app entitlements.

Map UI is a coarse placeholder until provider keys are supplied.

## Hackathon BLE/NFC Mode

Run:

```sh
flutter pub get
flutter run --dart-define=FINDMESH_API_BASE_URL=http://localhost:8080 --dart-define=FINDMESH_USE_REAL_RADIO=true
```

The Debug screen exposes:

- BLE permission request through the real platform plugins.
- Foreground scanning for `FM_TAG` and `FM_ZONE` advertisements.
- Android BLE advertising for a demo lost tag or merchant zone.
- NFC NDEF read/write for `findmesh://tag-found?t=<token>` payloads.

Two-phone demo:

1. Install the app on two Android phones with Bluetooth and NFC.
2. On phone A, open Debug and tap `Advertise tag`.
3. On phone B, open Debug and tap `Scan tags`.
4. Phone B should show the `FM_TAG` ephemeral ID and RSSI.
5. Write an NFC payload to an NTAG-compatible sticker, then read it back with `Read NFC`.

iOS can read NFC with entitlements and can scan foreground BLE, but BLE advertising and background behavior are more constrained. Use Android for the hackathon radio demo.
