# FindMesh

FindMesh is a privacy-preserving lost-item recovery network built around user mobile apps, BLE/NFC lost-item tags, smart merchant counter stands, and a backend that matches anonymous sightings without exposing finder identity, owner identity, phone numbers, payment data, or exact merchant identity by default.

FindMesh is not people tracking, employee tracking, customer surveillance, live tracking, POS integration, or payment processing. Merchant stands may physically display merchant branding, UPI QR, app QR, or recovery QR, but the electronics only participate in lost-item recovery.

## Repository

- `backend/`: Go 1.23 API, worker, admin CLI, migrations, protocol services, tests, OpenAPI.
- `mobile/findmesh_app/`: Flutter app skeleton with user and merchant mode.
- `firmware/merchant-stand/`: ESP-IDF merchant counter stand firmware.
- `firmware/lost-tag/`: ESP-IDF-compatible lost tag firmware.
- `web/admin-dashboard/`: React/TypeScript admin dashboard.
- `docs/`: Architecture, deployment, data model, threat model, runbooks, and test plan.

## Local Setup

For a step-by-step demo timeline, use [HACKATHON_PLAN.md](./HACKATHON_PLAN.md). It explains what to run first, what terminals stay open, and when to use backend, admin, mobile, BLE, NFC, tests, and firmware commands.

Prerequisites:

- Go 1.23+
- Docker Desktop for full-stack local services
- Node.js 22+ for the admin dashboard
- Flutter stable plus Android Studio/platform tools for the mobile app
- A physical Android phone for BLE/NFC hackathon demos
- ESP-IDF for firmware builds

Full local stack:

```sh
cp .env.example .env
make dev
```

Backend only, fastest path for mobile/admin demos:

```sh
make backend
```

The API starts on `http://localhost:8080`. Local development uses OTP `123456`.

Admin dashboard:

```sh
cd web/admin-dashboard
npm install
npm run dev
```

Open `http://localhost:5173`. The local admin token is `dev-admin-token`.

Flutter app:

```sh
cd mobile/findmesh_app
flutter pub get
flutter run --dart-define=FINDMESH_API_BASE_URL=http://localhost:8080
```

To bypass phone entry during local demos, tap `Continue with demo login` on the login screen. It calls the local OTP endpoints with phone `+15550000000` and dev OTP `123456`, so the app still receives a normal bearer token. Disable this button with:

```sh
flutter run --dart-define=FINDMESH_API_BASE_URL=http://localhost:8080 --dart-define=FINDMESH_ENABLE_DEV_LOGIN=false
```

Hackathon BLE/NFC mode on Android:

```sh
adb reverse tcp:8080 tcp:8080 # optional when using a USB-connected Android phone
make mobile-run-hackathon
```

This enables real BLE/NFC plugins with `FINDMESH_USE_REAL_RADIO=true`. Open Debug from the app toolbar to scan FindMesh BLE advertisements, advertise a demo tag or merchant zone from another phone, and read/write `findmesh://...` NFC payloads.

If the phone is not connected over USB, replace `localhost` with your laptop LAN IP:

```sh
cd mobile/findmesh_app
flutter run --dart-define=FINDMESH_API_BASE_URL=http://192.168.1.3:8080 --dart-define=FINDMESH_USE_REAL_RADIO=true
```

Two-phone BLE demo:

1. Install the app on two physical Android phones.
2. Phone A: open Debug and tap `Advertise tag`.
3. Phone B: open Debug and tap `Scan tags`.
4. Phone B should display `FM_TAG`, an ephemeral ID, and RSSI.
5. Use `Advertise zone` and `Scan zones` for the merchant stand beacon demo.

NFC demo:

1. Use an NTAG213/215/216 sticker or card.
2. Open Debug and keep `findmesh://tag-found?t=demo-lost-token` in the NFC payload field.
3. Tap `Write NFC`, then hold the NFC tag to the phone.
4. Tap `Read NFC`, then hold the same tag to confirm the payload.

Firmware:

```sh
cd firmware/merchant-stand && idf.py build
cd firmware/lost-tag && idf.py build
```

## Tests

```sh
make test
make mobile-test
```

Backend tests cover auth, pairing, lost mode, sightings, signature validation, deduplication, confidence scoring, abuse reporting, and API smoke flows. Firmware includes protocol test vectors and module-level seams for Unity tests under ESP-IDF.

## Demo Flow

1. Start the API.
2. POST `/v1/auth/otp/start` and `/v1/auth/otp/verify` with OTP `123456`.
3. Pair a tag through `/v1/tags/pair/complete`.
4. Mark the tag lost through `/v1/tags/{tag_id}/lost-mode`.
5. Upload a signed stand sighting or authenticated user-app sighting.
6. Read `/v1/tags/{tag_id}/last-seen`.
7. Use `/v1/recovery/requests` for masked merchant-assisted recovery.

## Assumptions

- ESP32-class merchant stand with BLE, Wi-Fi, PN532-compatible NFC, LED, optional battery, and secure private key storage.
- Lost tag uses BLE advertising, NFC lost-mode payload, buzzer, battery monitor, pairing button, and a per-tag secret.
- Production PostgreSQL/PostGIS storage is represented by migrations; local API uses an in-memory repository for fast runnable tests.
- NATS is optional and abstracted for production eventing; local development defaults to in-process behavior.
- Real SMS, push, map SDKs, BLE/NFC plugins, secure element signing, and hardware drivers are platform boundaries with mocks or stubs.

## Privacy Defaults

- Phone numbers are hashed and encrypted for storage and never included in sightings.
- Tags broadcast rotating 16-byte ephemeral IDs, not serials or owner IDs.
- Last seen views show coarse areas and confidence levels, not exact finder identity.
- Merchant identity is hidden by default and only revealed through recovery flow with merchant opt-in.
- Raw sightings default to 30-day retention.
