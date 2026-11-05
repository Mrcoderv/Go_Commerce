# Go_Commerce

Simple marketplace (Seller, Product, Buyer, Order) in Go. Hexagonal layout, Postgres (pgx), Uber Fx, Huma HTTP, slog.

## Prerequisites
- Go toolchain (see `go.mod`)
- PostgreSQL (or Docker)

## Setup
1. Copy env: `cp .env.example .env` and adjust `DB_*` / `APP_PORT`.
2. Start DB (optional): `docker compose up -d postgres` (maps host `${DB_PORT:-5434}` -> container 5432; match `.env`).
3. Migrations run automatically on app start (`migrations/*.sql`, `IF NOT EXISTS`, safe to re-run).

## Run
```bash
go run ./cmd/api
# health: curl http://localhost:$APP_PORT/health
```

## Test
```bash
go test ./...
# integration (needs DB):
TEST_DATABASE_URL='postgres://postgres:postgres@localhost:5432/go_commerce?sslmode=disable' go test -count=1 ./tests/...
```

## Layout (short)
- `cmd/api` – entrypoint + Fx wiring
- `internal/domain` – entities/errors/status
- `internal/ports/{incoming,outgoing}` – interfaces
- `internal/core/*` – use-case services
- `internal/adapters/incoming/http` – Huma handlers/DTO/router
- `internal/adapters/outgoing/postgres` – repos + `TxManager`
- `internal/infrastructure` – config/db/logger/server
- `migrations` – SQL

Docs: [docs/CODEMAP.md](docs/CODEMAP.md), [docs/API.md](docs/API.md), [docs/decisions.md](docs/decisions.md).
ggg