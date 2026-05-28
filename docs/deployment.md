# Deployment

Local:

```sh
cp .env.example .env
docker compose up --build
```

Production baseline:

- API behind TLS load balancer.
- PostgreSQL with PostGIS and encrypted backups.
- Redis for rate limiting and job coordination.
- NATS for event delivery.
- Object storage/CDN for signed firmware binaries.
- KMS/HSM for backend secrets.
- SSO for admin dashboard.
- Central logs, traces, and Prometheus metrics.

Run migrations with `golang-migrate` before deploying a new backend revision.
