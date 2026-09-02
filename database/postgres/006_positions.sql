CREATE TABLE IF NOT EXISTS positions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_address  VARCHAR(64) NOT NULL REFERENCES wallets(address),
    token_address   VARCHAR(64) NOT NULL REFERENCES tokens(address),
    quantity        NUMERIC(38, 18) NOT NULL DEFAULT 0,
    average_entry   NUMERIC(38, 18) NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (wallet_address, token_address)
);

CREATE INDEX IF NOT EXISTS idx_positions_wallet ON positions (wallet_address);