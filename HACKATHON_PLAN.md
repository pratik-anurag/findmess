# FindMesh Hackathon Run Plan

This is the practical runbook for demo day: what to start, when to start it, what stays running, and what to click in the app.

## 0. What You Need

Install these before the demo:

- Go 1.23+
- Node.js 22+
- Flutter stable
- Android Studio or Android platform tools
- A physical Android phone for BLE/NFC
- Optional: a second Android phone for the phone-to-phone BLE demo
- Optional: NTAG213/215/216 NFC sticker or card
- Optional: Docker Desktop if you want the full Compose stack

You do not need ESP-IDF unless you are building firmware during the demo.

## 1. Before the Demo

Run these once from the repo root:

```sh
cd /Users/pratikanurag/hack-mvp/findmesh
cd web/admin-dashboard && npm install && cd ../..
cd mobile/findmesh_app && flutter pub get && cd ../..
```

If you are using a USB-connected Android phone:

```sh
adb devices
adb reverse tcp:8080 tcp:8080
```

Use `adb reverse` so the phone can reach the laptop backend at `http://localhost:8080`.

## 2. Terminal 1: Start the Backend

Run this first and keep it running:

```sh
cd /Users/pratikanurag/hack-mvp/findmesh
make backend
```

Expected result:

```text
findmesh api starting addr=:8080
```

Health check:

```sh
curl http://localhost:8080/healthz
```

Expected response:

```json
{"status":"ok"}
```

The local OTP is always:

```text
123456
```

## 3. Terminal 2: Start the Admin Dashboard

Run this after the backend:

```sh
cd /Users/pratikanurag/hack-mvp/findmesh/web/admin-dashboard
npm run dev
```

Open:

```text
http://localhost:5173
```

Use admin token:

```text
dev-admin-token
```

Show this during the demo when you want to explain:

- users
- tags
- merchants
- stands
- sightings
- recovery requests
- abuse reports
- audit events

## 4. Terminal 3: Run the Mobile App Without Real BLE/NFC

Use this for simulator or backend-only app demos:

```sh
cd /Users/pratikanurag/hack-mvp/findmesh/mobile/findmesh_app
flutter run --dart-define=FINDMESH_API_BASE_URL=http://localhost:8080
```

On the login screen, tap:

```text
Continue with demo login
```

This bypasses manual phone entry by calling the local OTP flow with:

```text
phone: +15550000000
otp: 123456
```

Use this mode when you want to demo screens quickly without phone radio issues.

## 5. Terminal 3 Alternative: Run the Mobile App With Real BLE/NFC

Use this for the hackathon phone demo:

```sh
cd /Users/pratikanurag/hack-mvp/findmesh
adb reverse tcp:8080 tcp:8080
make mobile-run-hackathon
```

This runs Flutter with:

```text
FINDMESH_USE_REAL_RADIO=true
```

On the login screen, tap:

```text
Continue with demo login
```

If the phone is not connected over USB, use your laptop LAN IP instead:

```sh
cd /Users/pratikanurag/hack-mvp/findmesh/mobile/findmesh_app
flutter run \
  --dart-define=FINDMESH_API_BASE_URL=http://YOUR_LAPTOP_IP:8080 \
  --dart-define=FINDMESH_USE_REAL_RADIO=true
```

## 6. Main Demo Path

Run this sequence in the mobile app:

1. Tap `Continue with demo login`.
2. Go to `Tags`.
3. Tap `Pair tag`.
4. Use demo serial `FM-TAG-DEV-1`.
5. Label it `Keys`.
6. Open the tag detail.
7. Tap `Mark lost`.
8. Keep the safe message as `If found, contact me via FindMesh.`
9. Tap `Enable lost mode`.
10. Open `Privacy` from the top bar and explain what is hidden.
11. Open `Safety` and show unknown tracker reporting.
12. Open the admin dashboard and show records/audit views.

Use this message while presenting:

```text
FindMesh is lost-item recovery using anonymous sightings. It is not live tracking and does not reveal finder identity, owner phone number, or exact merchant identity by default.
```

## 7. BLE Demo Path

Use two Android phones.

Phone A:

1. Open the app.
2. Tap `Continue with demo login`.
3. Tap the bug/debug icon.
4. Tap `Advertise tag`.

Phone B:

1. Open the app.
2. Tap `Continue with demo login`.
3. Tap the bug/debug icon.
4. Tap `Scan tags`.
5. Confirm it shows `FM_TAG`, an ephemeral ID, and RSSI.

Merchant-zone version:

1. Phone A: tap `Advertise zone`.
2. Phone B: tap `Scan zones`.
3. Confirm it shows `FM_ZONE`.

## 8. NFC Demo Path

Use a real NFC sticker or card.

1. Open Debug in the mobile app.
2. Keep payload:

```text
findmesh://tag-found?t=demo-lost-token
```

3. Tap `Write NFC`.
4. Hold the NFC tag to the phone.
5. Tap `Read NFC`.
6. Hold the same tag to the phone.
7. Confirm the payload is read back.

Use this to explain the found-item flow:

```text
Anyone who finds an item can tap NFC and report it through FindMesh without seeing the owner's phone number.
```

## 9. When to Run Tests

Before pushing code:

```sh
cd /Users/pratikanurag/hack-mvp/findmesh
make test
cd web/admin-dashboard && npm run build
```

If Flutter is installed:

```sh
cd /Users/pratikanurag/hack-mvp/findmesh/mobile/findmesh_app
flutter test
```

Do not run firmware builds unless ESP-IDF is installed:

```sh
cd firmware/merchant-stand && idf.py build
cd firmware/lost-tag && idf.py build
```

## 10. Troubleshooting

Backend is not reachable from phone:

```sh
adb reverse tcp:8080 tcp:8080
curl http://localhost:8080/healthz
```

If not using USB, replace `localhost` with your laptop LAN IP in the Flutter command.

Login fails:

- Check backend is running.
- Use `Continue with demo login`.
- If entering manually, use OTP `123456`.

BLE scan finds nothing:

- Use physical Android phones, not emulator.
- Turn Bluetooth on.
- Grant Bluetooth and location permissions.
- Start advertising on Phone A before scanning on Phone B.
- Keep both phones close during the demo.

NFC does not read:

- Use a real NTAG-compatible sticker/card.
- Make sure NFC is enabled in Android settings.
- Hold the tag still against the phone NFC antenna area.

Admin dashboard is empty:

- Check the token is `dev-admin-token`.
- Check backend is running on `http://localhost:8080`.
- Refresh the dashboard.

## 11. Recommended Demo Order

1. Open admin dashboard.
2. Start mobile app and demo-login.
3. Pair a tag.
4. Mark it lost.
5. Show privacy screen.
6. Show anti-stalking screen.
7. Show BLE phone-to-phone scan.
8. Show NFC write/read.
9. Return to admin dashboard and show audit/abuse/stand areas.

