-- 002_buyers_orders.sql
-- Depends on: nothing (no foreign keys to sellers/products on purpose:
-- Order only stores their IDs, see bounded-context rule).
-- Money is stored as BIGINT in paisa (no float / decimal rounding).

-- ---------------------------------------------------------------
-- buyers
-- ---------------------------------------------------------------
CREATE TABLE buyers (
                        id          BIGSERIAL    PRIMARY KEY,
                        name        VARCHAR(255) NOT NULL,
                        email       VARCHAR(255) NOT NULL,
                        phone       VARCHAR(30)  NOT NULL DEFAULT '',   -- text: keeps +977 and leading zeros
                        address     TEXT         NOT NULL DEFAULT '',
                        created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
                        updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
                        deleted_at  TIMESTAMPTZ                          -- soft delete
);

-- Unique email only among active buyers, so a deleted buyer's email can be reused.
CREATE UNIQUE INDEX uq_buyers_email_active
    ON buyers (LOWER(email))
    WHERE deleted_at IS NULL;

-- ---------------------------------------------------------------
-- orders
-- ---------------------------------------------------------------
CREATE TABLE orders (
                        id            BIGSERIAL   PRIMARY KEY,
                        buyer_id      BIGINT      NOT NULL,
                        status        VARCHAR(30) NOT NULL DEFAULT 'pending',
                        total_amount  BIGINT      NOT NULL DEFAULT 0,    -- paisa
                        created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                        updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                        CONSTRAINT fk_orders_buyer
                            FOREIGN KEY (buyer_id) REFERENCES buyers (id)
                                ON DELETE RESTRICT,

                        CONSTRAINT chk_orders_status
                            CHECK (status IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled')),

                        CONSTRAINT chk_orders_total_amount
                            CHECK (total_amount >= 0)
);

-- ---------------------------------------------------------------
-- order_items
-- seller_id, product_name and unit_price are SNAPSHOTS taken at order time.
-- No FK to products/sellers: Order owns its own copy of the data.
-- ---------------------------------------------------------------
CREATE TABLE order_items (
                             id            BIGSERIAL    PRIMARY KEY,
                             order_id      BIGINT       NOT NULL,
                             product_id    BIGINT       NOT NULL,              -- ID only
                             seller_id     BIGINT       NOT NULL,              -- ID only
                             product_name  VARCHAR(255) NOT NULL,              -- snapshot
                             unit_price    BIGINT       NOT NULL,              -- snapshot, paisa
                             quantity      INT          NOT NULL,
                             created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

                             CONSTRAINT fk_order_items_order
                                 FOREIGN KEY (order_id) REFERENCES orders (id)
                                     ON DELETE CASCADE,

                             CONSTRAINT chk_order_items_quantity
                                 CHECK (quantity > 0),

                             CONSTRAINT chk_order_items_unit_price
                                 CHECK (unit_price >= 0)
);
CREATE INDEX idx_orders_buyer_created  ON orders (buyer_id, created_at DESC);  -- ListByBuyer
CREATE INDEX idx_order_items_order_id  ON order_items (order_id);              -- load items
CREATE INDEX idx_order_items_seller_id ON order_items (seller_id);             -- ListBySeller
CREATE INDEX idx_order_items_product_id ON order_items (product_id);