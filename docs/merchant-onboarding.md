# Merchant Onboarding

1. Merchant creates profile in app.
2. Field technician or merchant taps stand NFC setup URL.
3. App connects to stand provisioning BLE service.
4. Merchant configures Wi-Fi.
5. Backend binds stand to merchant and coarse zone.
6. Stand starts scanning for FindMesh tag advertisements and advertises zone beacon.
7. Merchant can enable or disable assisted recovery.

Merchant staff never see owner phone by default. Exact merchant details are only used in recovery flow when assistance is enabled.
