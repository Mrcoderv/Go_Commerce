.PHONY: run build vet test test-integration db-up migrate

run:
	go run ./cmd/api

build:
	go build -v ./...

vet:
	go vet ./...

test:
	go test -v ./...

test-integration:
	TEST_DATABASE_URL=$${TEST_DATABASE_URL:-"postgres://postgres:postgres@localhost:5432/go_commerce?sslmode=disable"} go test -count=1 -v ./tests/...

db-up:
	docker compose up -d postgres

migrate:
	go run ./cmd/api --migrate-only || true
