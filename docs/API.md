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
