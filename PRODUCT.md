# Product

FindMesh helps owners recover lost items through anonymous sightings from opted-in mobile devices and participating merchant counter stands.

## Personas

- Tag owner: pairs tags, marks items lost, views approximate last-seen area, uses nearby finder, confirms recovery.
- Finder user: opts into anonymous finder participation and can report found NFC/QR items without seeing owner identity.
- Merchant: claims stands, enables merchant-assisted recovery, sees masked requests.
- Merchant staff: responds to accepted recovery requests without accessing owner phone by default.
- Admin/support operator: manages users, tags, merchants, stands, firmware, health, and data deletion.
- Abuse/safety reviewer: reviews unwanted tracker alerts and suspicious activity.
- Field technician: provisions or replaces merchant stands.

## Product Boundaries

FindMesh does not process payments, read POS systems, access UPI transactions, track people, expose phone numbers in sightings, show live movement trails, or reveal exact merchant/finder/owner identity by default.

## Core Flows

- OTP onboarding.
- BLE/NFC tag pairing.
- Rotating ephemeral ID broadcasts.
- Lost mode with safe message.
- Merchant stand signed sightings and offline buffering.
- User phone anonymous finder network.
- Coarse last-seen summary.
- Nearby finding and ring intent.
- Merchant-assisted masked recovery.
- Found item NFC/QR report.
- Stand onboarding and health.
- OTA manifest and signed firmware update flow.
- Anti-stalking unknown tag alerts and abuse reports.
- Account, tag, finder participation, and retention controls.
