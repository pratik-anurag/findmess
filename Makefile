SHELL := /bin/sh

.PHONY: dev test migrate backend mobile-test firmware-stand-build firmware-tag-build lint seed admin-dev

dev:
	docker compose up --build

backend:
	cd backend && go run ./cmd/api

test:
	cd backend && go test ./...

migrate:
	migrate -path backend/migrations -database "$${DATABASE_URL}" up

seed:
	psql "$${DATABASE_URL}" -f backend/scripts/seed_demo.sql

mobile-test:
	cd mobile/findmesh_app && flutter test

mobile-run-hackathon:
	cd mobile/findmesh_app && flutter pub get && flutter run --dart-define=FINDMESH_API_BASE_URL=http://localhost:8080 --dart-define=FINDMESH_USE_REAL_RADIO=true

firmware-stand-build:
	cd firmware/merchant-stand && idf.py build

firmware-tag-build:
	cd firmware/lost-tag && idf.py build

admin-dev:
	cd web/admin-dashboard && npm install && npm run dev

lint:
	cd backend && go test ./...
	cd mobile/findmesh_app && flutter analyze
	cd web/admin-dashboard && npm run build
