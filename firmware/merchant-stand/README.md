# FindMesh Merchant Stand Firmware

Default target: ESP32-class device with BLE, Wi-Fi, NFC through PN532-compatible module, LED, USB-C power, optional battery, and secure key storage.

States:
- `UNCLAIMED`
- `PROVISIONING`
- `ONLINE`
- `OFFLINE_BUFFERING`
- `RECOVERY_MODE`
- `ERROR`
- `OTA_UPDATING`

The stand is a zone witness. It does not process payments, POS data, UPI IDs, phone numbers, finder identity, or owner identity.

Build:

```sh
idf.py set-target esp32
idf.py build
```

Production replacements:
- Store the Ed25519 private key in secure element/NVS encryption.
- Replace `crypto_signer` placeholder with hardware-backed Ed25519.
- Replace PN532 stub with board-specific I2C/SPI driver.
- Pin TLS backend upload to the FindMesh API certificate policy.
