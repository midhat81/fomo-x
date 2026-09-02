CREATE TABLE IF NOT EXISTS transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    signature       VARCHAR(128) NOT NULL UNIQUE,
    wallet_address  VARCHAR(64) NOT NULL REFERENCES wallets(address),
    success         BOOLEAN NOT NULL DEFAULT true,
    raw_logs        JSONB,
    processed_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_transactions_wallet ON transactions (wallet_address);
CREATE INDEX IF NOT EXISTS idx_transactions_signature ON transactions (signature);