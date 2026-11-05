CREATE TABLE IF NOT EXISTS buyers (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    phone       VARCHAR(30)  NOT NULL DEFAULT '',
    address     TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_buyers_email ON buyers (LOWER(email));

CREATE TABLE IF NOT EXISTS orders (
    id            BIGSERIAL   PRIMARY KEY,
    buyer_id      BIGINT      NOT NULL REFERENCES buyers (id) ON DELETE RESTRICT,
    status        VARCHAR(30) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','returned', 'confirmed', 'shipped', 'delivered', 'cancelled')),
    total_amount  BIGINT      NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id            BIGSERIAL    PRIMARY KEY,
    order_id      BIGINT       NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id    BIGINT       NOT NULL,
    seller_id     BIGINT       NOT NULL,
    product_name  VARCHAR(255) NOT NULL,
    unit_price    BIGINT       NOT NULL CHECK (unit_price >= 0),
    quantity      INT          NOT NULL CHECK (quantity > 0),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_buyer_id ON orders (buyer_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items (order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_seller_id ON order_items (seller_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items (product_id);
