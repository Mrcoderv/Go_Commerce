# CODEMAP

## Key files
- `cmd/api/main.go` – Fx app entrypoint and wiring.
- `internal/domain/*` – entities, statuses, sentinel errors.
- `internal/ports/incoming/*` – use-case interfaces and inputs.
- `internal/ports/outgoing/*` – repo and `TxManager` interfaces.
- `internal/core/*/services.go` – business logic; `module.go` Fx modules.
- `internal/adapters/incoming/http/router.go` – Huma/mux routes, static `./files` serving.
- `internal/adapters/incoming/http/handlers/*` – HTTP handlers (DTO -> service -> DTO).
- `internal/adapters/incoming/http/dto/*` – request/response DTOs.
- `internal/adapters/outgoing/postgres/tx.go` – tx context helper and `TxManager`.
- `internal/adapters/outgoing/postgres/repositories/*` – pgx/sql repos.
- `internal/infrastructure/{config,database,logger,server}` – config, pool, logging, HTTP server.
- `migrations/*.sql`, `migrations/embed.go` – embedded SQL migrations run at startup.

## PlaceOrder flow
`POST /orders` -> `handlers.OrderHandler.PlaceOrder` (`internal/adapters/incoming/http/handlers/order.go`)
-> `incoming.OrderService.PlaceOrder` -> `core/order.OrderService.PlaceOrder` (`internal/core/order/services.go`)
-> `TxManager.WithinTx` (`internal/adapters/outgoing/postgres/tx.go`) runs in one tx:
buyer lookup, product reads/snapshots, `ProductRepository.DecreaseStock` (single `UPDATE ... stock>=qty`),
`OrderRepository.Create` inserts order+items. Commit/rollback via tx manager. Failure rolls back stock changes.
