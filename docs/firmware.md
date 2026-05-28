# Firmware

Merchant stand firmware lives in `firmware/merchant-stand`.

Lost tag firmware lives in `firmware/lost-tag`.

Merchant stand modules:

- BLE scanner and advertiser.
- Wi-Fi manager.
- NFC manager.
- Offline sighting buffer.
- Crypto signer.
- OTA manager.
- LED status.
- Power manager.
- Provisioning.

Lost tag modules:

- BLE advertiser.
- Ephemeral ID generator.
- NFC lost-mode payload.
- Buzzer.
- Battery.
- Pairing.
- Anti-stalking.

Hardware-specific BLE, NFC, secure key storage, battery, and buzzer drivers are intentionally isolated.
