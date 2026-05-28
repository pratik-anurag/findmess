# FindMesh Lost Tag Firmware

Default target: BLE-capable MCU using ESP-IDF-compatible components. The same module boundaries can be mapped to nRF Connect SDK.

States:
- `UNPAIRED`
- `PAIRING`
- `PAIRED_NORMAL`
- `LOST_MODE`
- `SEPARATED`
- `RINGING`
- `DISABLED`

The tag never advertises raw serial, owner ID, phone number, or payment data. It advertises a rotating 16-byte ephemeral ID derived from a per-tag secret.

Build:

```sh
idf.py set-target esp32
idf.py build
```

Production replacements:
- Store tag secret in secure storage.
- Replace button/buzzer/battery stubs with board-specific drivers.
- Align unwanted-tracker behavior with regional safety requirements and platform standards.
