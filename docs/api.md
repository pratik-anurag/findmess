# API

OpenAPI is available at `backend/openapi/openapi.yaml`.

Authentication:

- Users authenticate with OTP and receive a bearer token.
- Local development uses OTP `123456`.
- Admin uses `FINDMESH_ADMIN_TOKEN`.
- User app sightings require authenticated upload.
- Stand sightings require Ed25519 signature validation.

Main groups:

- `/v1/auth/*`
- `/v1/tags/*`
- `/v1/sightings*`
- `/v1/merchants*`
- `/v1/stands*`
- `/v1/recovery*`
- `/v1/found*`
- `/v1/abuse*`
- `/v1/admin*`
- `/v1/firmware*`

Normal users receive coarse last-seen summaries. Admin endpoints are audited.
