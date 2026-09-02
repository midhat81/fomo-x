CREATE TABLE IF NOT EXISTS trades (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL UNIQUE,
    transaction_id  UUID NOT NULL REFERENCES transactions(id),
    wallet_address  VARCHAR(64) NOT NULL REFERENCES wallets(address),
    token_address   VARCHAR(64) NOT NULL REFERENCES tokens(address),
    side            VARCHAR(8) NOT NULL CHECK (side IN ('BUY', 'SELL', 'UNKNOWN')),
    quantity        NUMERIC(38, 18) NOT NULL DEFAULT 0,
    price           NUMERIC(38, 18) NOT NULL DEFAULT 0,
    signature       VARCHAR(128) NOT NULL,
    traded_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_trades_wallet ON trades (wallet_address);
CREATE INDEX IF NOT EXISTS idx_trades_token ON trades (token_address);
CREATE INDEX IF NOT EXISTS idx_trades_traded_at ON trades (traded_at);