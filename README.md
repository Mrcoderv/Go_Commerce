# Go_Commerce
go 

task for the developer 1 and developer 2 .




 [docs/decisions.md](docs/decisions.md)

## Database setup

Start Postgres with Docker:
docker compose up -d

Run migrations (auto-applied via docker-entrypoint-initdb.d on first start).

Run the app:
go run ./cmd/main.go

Check DB connection:
go run ./cmd/main.go  (should print "Database connection established")

Verify tables:
psql -h localhost -U postgres -d go_commerce -c "\dt" 

ghjkl
