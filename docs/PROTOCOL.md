# Protocol

The canonical protocol document lives at `../PROTOCOL.md`. Keep implementation test vectors in sync across:

- `backend/internal/crypto/crypto_test.go`
- `firmware/lost-tag/main/test_vectors.c`
- mobile BLE/NFC parsing in `mobile/findmesh_app/lib/platform`
