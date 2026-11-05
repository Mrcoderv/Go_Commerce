CREATE TABLE IF NOT EXISTS sellers (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    phone       VARCHAR(30)  NOT NULL DEFAULT '',
    address     TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_sellers_email ON sellers (LOWER(email));

CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL    PRIMARY KEY,
    seller_id   BIGINT       NOT NULL REFERENCES sellers (id) ON DELETE RESTRICT,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    price       BIGINT       NOT NULL CHECK (price > 0),
    stock       INTEGER      NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products (seller_id);
