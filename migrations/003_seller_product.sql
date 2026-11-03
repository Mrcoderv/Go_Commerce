CREATE TABLE IF NOT EXISTS sellers (
                        id          BIGSERIAL    PRIMARY KEY,
                        name        VARCHAR(255) NOT NULL,
                        email       VARCHAR(255) NOT NULL,
                        phone       VARCHAR(30)  NOT NULL DEFAULT '',   -- text: keeps +977 and leading zeros
                        address     TEXT         NOT NULL DEFAULT '',
                        created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
                        updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
                        deleted_at  TIMESTAMPTZ                          -- soft delete
);

CREATE UNIQUE INDEX uq_sellers_email_active
    ON sellers (LOWER(email))
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS products (
    id         BIGSERIAL    PRIMARY KEY,
    seller_id BIGINT       NOT NULL REFERENCES sellers (id) ON DELETE RESTRICT,
    product_name VARCHAR(255) NOT NULL,
    unit_price BIGINT       NOT NULL CHECK (unit_price >= 0), -- paisa
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ                          -- soft delete
);

CREATE INDEX idx_products_seller_id
    ON products (seller_id)
    WHERE deleted_at IS NULL;

