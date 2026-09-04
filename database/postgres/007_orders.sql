CREATE TABLE IF NOT EXISTS orders (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key   VARCHAR(128) NOT NULL UNIQUE,
    wallet_address    VARCHAR(64) NOT NULL REFERENCES wallets(address),
    token_address     VARCHAR(64) NOT NULL REFERENCES tokens(address),
    side              VARCHAR(8) NOT NULL CHECK (side IN ('BUY', 'SELL')),
    order_type        VARCHAR(16) NOT NULL DEFAULT 'MARKET',
    quantity          NUMERIC(38, 18) NOT NULL,
    status            VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    reject_reason     TEXT,
    filled_quantity   NUMERIC(38, 18),
    fill_price        NUMERIC(38, 18),
    slippage_bps      NUMERIC(10, 4),
    fee_amount        NUMERIC(38, 18),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_wallet ON orders (wallet_address);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);