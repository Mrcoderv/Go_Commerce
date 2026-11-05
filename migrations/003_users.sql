CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL    PRIMARY KEY,
    email         VARCHAR(255) NOT NULL,
    password_hash TEXT         NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('seller', 'buyer')),
    seller_id     BIGINT       REFERENCES sellers(id) ON DELETE CASCADE,
    buyer_id      BIGINT       REFERENCES buyers(id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_user_role_ref CHECK (
        (role = 'seller' AND seller_id IS NOT NULL AND buyer_id IS NULL) OR
        (role = 'buyer' AND buyer_id IS NOT NULL AND seller_id IS NULL)
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_users_email ON users (LOWER(email));
CREATE INDEX IF NOT EXISTS idx_users_seller_id ON users (seller_id);
CREATE INDEX IF NOT EXISTS idx_users_buyer_id ON users (buyer_id);
