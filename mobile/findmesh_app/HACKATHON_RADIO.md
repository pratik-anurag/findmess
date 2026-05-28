# Hackathon BLE/NFC Setup

This app now has two radio modes:

- default: mock BLE/NFC services for simulator and backend demos.
- real radio: Flutter BLE/NFC plugins enabled with `FINDMESH_USE_REAL_RADIO=true`.

## Android Demo

Start the backend:

```sh
cd ../../
make backend
```

Run the app on a physical Android phone:

```sh
cd mobile/findmesh_app
adb reverse tcp:8080 tcp:8080
flutter pub get
flutter run --dart-define=FINDMESH_API_BASE_URL=http://localhost:8080 --dart-define=FINDMESH_USE_REAL_RADIO=true
```

On the login screen, tap `Continue with demo login` to skip manual phone entry. It uses the local dev OTP flow with phone `+15550000000` and OTP `123456`.

If you are not using USB reverse port forwarding and the phone is on the same Wi-Fi network as your laptop, replace `localhost` with the laptop LAN IP:

```sh
flutter run --dart-define=FINDMESH_API_BASE_URL=http://192.168.1.3:8080 --dart-define=FINDMESH_USE_REAL_RADIO=true
```

## BLE Demo

Use two Android phones:

1. Phone A: Debug → `Advertise tag`.
2. Phone B: Debug → `Scan tags`.
3. Phone B should display `FM_TAG`, the ephemeral ID, and RSSI.
4. Phone A can also advertise a merchant zone with `Advertise zone`.
5. Phone B can scan zones with `Scan zones`.

Advertisement format for demo phones:

- Service UUID for external hardware: `f17d0001-2d9a-4c8a-a4f1-f0d641b90f10`
- Android phone-to-phone demo uses manufacturer data only to stay under BLE advertisement size limits.
- iOS phone-to-phone advertising is limited to local-name style payloads by platform behavior.
- Manufacturer ID: `0xF17D`
- Payload: `[type, flags, 16-byte ephemeral_id]`
- `type=0x01`: `FM_TAG`
- `type=0x02`: `FM_ZONE`

## NFC Demo

Use an NTAG213/215/216 sticker or card.

1. Debug → edit NFC payload, for example `findmesh://tag-found?t=demo-lost-token`.
2. Tap `Write NFC` and hold the tag to the phone.
3. Tap `Read NFC` and hold the same tag to the phone.
4. The app should read the payload back.

## Android Permissions

The Android project includes:

- `BLUETOOTH_SCAN`
- `BLUETOOTH_CONNECT`
- `BLUETOOTH_ADVERTISE`
- legacy Bluetooth permissions for Android 11 and lower
- location permission for Android 11 and lower scan compatibility
- `NFC`
- `INTERNET`

Use real Android hardware. Emulators usually do not provide BLE advertising or NFC.
