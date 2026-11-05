# API (short)

Base: `http://localhost:<APP_PORT>` (default 8081; `.env` may override). JSON unless noted.

## Health
- `GET /health` -> `200` `{"status":"health is good",...}`

## Sellers
- `POST /sellers` body `{name,email,phone?,address?}` -> `201` seller
- `GET /sellers?limit&offset` -> `200` `[...]`
- `GET /sellers/{id}` -> `200` seller | `404`
- `PATCH /sellers/{id}` body partial -> `200` | `404` `409`
- `DELETE /sellers/{id}` -> `204` | `404`

## Products
- `POST /products` header `X-Seller-ID` (or body `seller_id`) body `{name,description?,price|unit_price,stock}` -> `201`
  - `price` is a plain number in DTO (e.g. `499.50`) converted to paisa at handler; domain/DB store paisa (`int64`).
- `GET /products?seller_id&limit&offset` -> `200` `[...]`
- `GET /products/{id}` -> `200` | `404`
- `PATCH /products/{id}` header/body seller id + partial fields -> `200` | `403` `404`
- `DELETE /products/{id}?seller_id=` or header `X-Seller-ID` -> `204` | `403` `404` `409` (has order items)

## Buyers
- `POST /buyers` `{name,email,phone?,address?}` -> `201`
- `GET /buyers` -> `200` `[...]`
- `GET /buyers/{id}` -> `200` | `404`
- `PATCH /buyers/{id}` -> `200` | `404`
- `DELETE /buyers/{id}` -> `204` | `404` `409` (has orders)

## Orders
- `POST /orders` `{buyer_id, items:[{product_id,quantity}]}` -> `201` order (snapshots name/price/seller)
  Errors: `400` invalid, `404` missing buyer/product, `409` insufficient stock.
- `GET /orders/{id}` -> `200` | `404`
- `PATCH /orders/{id}/status` `{status}` -> `200` | `404` `409` invalid transition
- `POST /orders/{id}/cancel` -> `200` (restores stock in tx) | `404` `409`
- `GET /buyers/{id}/orders` -> `200` `[...]`
- `GET /sellers/{id}/orders` -> `200` `[...]`

## Page to endpoint map
| Page | UI action | Method/Path | Request fields | Response/notes |
|---|---|---|---|---|
| `files/manage-sellers.html` | Save seller form | `POST/PATCH /sellers` | `name,email,phone,address` | seller JSON; page is static (no JS fetch in repo) |
| `files/manage-buyers.html` | Save buyer form | `POST/PATCH /buyers` | `name,email,phone,address` | buyer JSON; static |
| `files/seller.html` | Save product form | `POST/PATCH /products` | `name,description,price,stock` (+ seller id via header/query as used by API) | product JSON; static |
| `files/buyer.html` | Place order / browse | `GET /products`, `POST /orders` | order: `buyer_id, items[{product_id,quantity}]` | order/product JSON; static |

> The HTML pages in `files/` have no embedded `fetch`/AJAX; they define the intended fields and flows. The API and routes are aligned to those field names and simple CRUD/order paths; static files are served from `/` by the Go server.

## Smoke test (curl, against running server)
```bash
BASE=http://localhost:8084
curl -s $BASE/health
curl -s -X POST $BASE/sellers -H 'Content-Type: application/json' -d '{"name":"S1","email":"s1@example.com","phone":"1","address":"a"}'
curl -s -X POST $BASE/products -H 'Content-Type: application/json' -H 'X-Seller-ID: <sid>' -d '{"name":"P1","description":"d","price":499.50,"stock":5}'
curl -s -X POST $BASE/buyers -H 'Content-Type: application/json' -d '{"name":"B1","email":"b1@example.com","phone":"2","address":"b"}'
curl -s -X POST $BASE/orders -H 'Content-Type: application/json' -d '{"buyer_id":<bid>,"items":[{"product_id":<pid>,"quantity":1}]}'
curl -s $BASE/orders/<oid>
curl -s $BASE/buyers/<bid>/orders
curl -s -X PATCH $BASE/orders/<oid>/status -H 'Content-Type: application/json' -d '{"status":"confirmed"}'
curl -s -X POST $BASE/orders/<oid>/cancel
curl -s -o /dev/null -w '%{http_code}' -X DELETE "$BASE/products/<pid>?seller_id=<sid>"
```

### Sample responses (from a real run)
- health: `{"$schema":".../schemas/HealthResponse.json","status":"health is good"}`
- product create (price 499.50 -> paisa): `{"id":4,"seller_id":5,"name":"P1","price":49950,"stock":5,...}`
- order: `{"id":4,"buyer_id":4,"total":49950,"status":"pending","items":[{"product_id":4,"unit_price":49950,"quantity":1,...}]}`
- cancel -> status `cancelled`; delete ordered product -> `409` conflict.
