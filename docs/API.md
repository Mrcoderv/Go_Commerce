# API testing guide

- Base URL: `http://localhost:8081` (from `APP_PORT` in `.env`)
- Interactive docs: `http://localhost:8081/docs`
- OpenAPI JSON: `http://localhost:8081/openapi.json`
- JSON schemas: `http://localhost:8081/schemas`

The `/docs` page is the original Huma/Stoplight Elements tester, generated from
the Go code. It includes:
- the local server URL, so **Try It** calls `localhost:8081`
- a **Bearer** security scheme, so the docs' authentication panel sends the JWT
- descriptions for every endpoint

---

## 1. Start the API

```bash
# 1. Start PostgreSQL (skip if you already have one running)
docker compose up -d postgres

# 2. Run the API (migrations run automatically at startup)
go run ./cmd/api
```

Check it is alive:

```bash
curl http://localhost:8081/health
# {"status":"health is good"}
```

Open <http://localhost:8081/docs>.

---

## 2. Test with the `/docs` page (recommended)

### Step 1 - Register a seller

In the docs page open **POST /auth/register** -> **Try It**, send:

```json
{
  "email": "seller1@example.com",
  "password": "secret123",
  "role": "seller",
  "name": "Seller One",
  "phone": "9800000001",
  "address": "Kathmandu"
}
```

The response contains a JWT and the new seller id:

```json
{
  "token": "eyJhbGciOi...",
  "user": { "id": 1, "role": "seller", "seller_id": 1, "buyer_id": null, "email": "seller1@example.com" }
}
```

Copy `token` and `user.seller_id`.

### Step 2 - Set the token

Open the **Auth** / authentication panel on the docs page, paste the token into the
**bearerAuth** field (do **not** type `Bearer ` yourself), and apply it.

### Step 3 - Create a product

Open **POST /products** -> **Try It**:

```json
{
  "name": "Laptop",
  "description": "14 inch",
  "price": 49950,
  "stock": 10
}
```

`price` is **integer paisa**: `49950` = NPR 499.50. There is no `seller_id` field;
the seller is taken from your token.

Copy `id` from the response.

### Step 4 - Register a buyer

Open **POST /auth/register** again with `role: "buyer"`:

```json
{
  "email": "buyer1@example.com",
  "password": "secret123",
  "role": "buyer",
  "name": "Buyer One",
  "phone": "9800000002",
  "address": "Lalitpur"
}
```

Copy the buyer `token` and `user.buyer_id`, then open the **Auth** panel again and
replace the token with the buyer token.

### Step 5 - Place an order

Open **POST /orders** -> **Try It**:

```json
{
  "items": [{ "product_id": 1, "quantity": 2 }]
}
```

No `buyer_id` is needed; it comes from the token. The response snapshots the product
name, seller and unit price, and returns the order `id`.

### Step 6 - Read the order

Open **GET /orders/{id}** and use the order id from step 5.

### Step 7 - Change the status (seller)

Re-authorize with the **seller** token, then call
**PATCH /orders/{id}/status**:

```json
{ "status": "confirmed" }
```

Legal transitions:

```
pending    -> confirmed | cancelled
confirmed  -> shipped   | cancelled
shipped    -> delivered
delivered  -> returned
```

### Step 8 - List orders

- Buyer token: **GET /buyers/{buyer_id}/orders**
- Seller token: **GET /sellers/{seller_id}/orders**

The path id must match the id in the JWT. A seller only sees their own items and subtotal.

### Step 9 - Cancel (buyer)

Re-authorize with the buyer token and call **POST /orders/{id}/cancel**.
Stock is restored in the same transaction.

---

## 3. Test with curl

```bash
BASE=http://localhost:8081

# --- Seller ---
SELLER=$(curl -s -X POST "$BASE/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{"email":"seller1@example.com","password":"secret123","role":"seller","name":"Seller One"}')
SELLER_TOKEN=$(echo "$SELLER" | jq -r .token)
SELLER_ID=$(echo "$SELLER" | jq -r .user.seller_id)

# --- Product ---
PRODUCT=$(curl -s -X POST "$BASE/products" \
  -H "Authorization: Bearer $SELLER_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Laptop","description":"14 inch","price":49950,"stock":10}')
PRODUCT_ID=$(echo "$PRODUCT" | jq -r .id)

# --- Buyer ---
BUYER=$(curl -s -X POST "$BASE/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{"email":"buyer1@example.com","password":"secret123","role":"buyer","name":"Buyer One"}')
BUYER_TOKEN=$(echo "$BUYER" | jq -r .token)
BUYER_ID=$(echo "$BUYER" | jq -r .user.buyer_id)

# --- Place order ---
ORDER=$(curl -s -X POST "$BASE/orders" \
  -H "Authorization: Bearer $BUYER_TOKEN" \
  -H 'Content-Type: application/json' \
  -d "{\"items\":[{\"product_id\":$PRODUCT_ID,\"quantity\":2}]}")
ORDER_ID=$(echo "$ORDER" | jq -r .id)

# --- Read order ---
curl -s "$BASE/orders/$ORDER_ID" -H "Authorization: Bearer $BUYER_TOKEN"

# --- Seller confirms ---
curl -s -X PATCH "$BASE/orders/$ORDER_ID/status" \
  -H "Authorization: Bearer $SELLER_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"status":"confirmed"}'

# --- Lists ---
curl -s "$BASE/buyers/$BUYER_ID/orders"   -H "Authorization: Bearer $BUYER_TOKEN"
curl -s "$BASE/sellers/$SELLER_ID/orders" -H "Authorization: Bearer $SELLER_TOKEN"

# --- Cancel (buyer) ---
curl -s -X POST "$BASE/orders/$ORDER_ID/cancel" -H "Authorization: Bearer $BUYER_TOKEN"
```

`jq` is optional; without it, read the `token`, `id` and `user.*_id` values from the JSON manually.

---

## 4. Endpoint reference

### Public

| Method | Path | Purpose |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/auth/register` | Register seller/buyer and get a JWT |
| POST | `/auth/login` | Log in and get a JWT |
| POST | `/sellers` | Legacy seller registration (no JWT) |
| GET | `/sellers` | List sellers |
| POST | `/buyers` | Legacy buyer registration (no JWT) |
| GET | `/buyers` | List buyers |
| GET | `/products` | List products (`?seller_id=&limit=&offset=`) |
| GET | `/products/{id}` | Get a product |

### Requires a JWT

| Method | Path | Who |
|---|---|---|
| GET / PATCH / DELETE | `/sellers/{id}` | That seller |
| GET / PATCH / DELETE | `/buyers/{id}` | That buyer |
| POST | `/products` | Any seller |
| PATCH / DELETE | `/products/{id}` | The owning seller |
| POST | `/orders` | Any buyer |
| GET | `/orders/{id}` | The buyer or a seller in the order |
| PATCH | `/orders/{id}/status` | A seller in the order |
| POST | `/orders/{id}/cancel` | The buyer who placed it |
| GET | `/buyers/{id}/orders` | That buyer |
| GET | `/sellers/{id}/orders` | That seller |

---

## 5. Rules and common errors

- **Money**: always integer paisa. `49950` = NPR 499.50. `499.50` is rejected.
- **Identity**: seller/buyer id always comes from the JWT. The old `X-Seller-ID`
  header, body `seller_id`/`buyer_id`, and seller query parameter no longer exist.
- **Path ids** on `/sellers/{id}`, `/buyers/{id}`, and the order list endpoints must
  match the token, otherwise the API returns `403`.
- **Status codes**:
  - `400` invalid input
  - `401` missing/invalid token
  - `403` token does not own the resource
  - `404` not found
  - `409` duplicate email, insufficient stock, invalid status transition,
    or deleting a buyer/product that is referenced by orders
- **Delete a product** returns `409` if it appears in any order item.
- **Delete a buyer** returns `409` if the buyer has orders.
- Use fresh emails when re-running the flow (emails are unique per role).
