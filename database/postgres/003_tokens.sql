CREATE TABLE IF NOT EXISTS tokens (
    address         VARCHAR(64) PRIMARY KEY,
    symbol          VARCHAR(32),
    name            VARCHAR(128),
    decimals        SMALLINT NOT NULL DEFAULT 9,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tokens_symbol ON tokens (symbol);