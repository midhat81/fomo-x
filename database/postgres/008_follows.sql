CREATE TABLE IF NOT EXISTS follows (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    follower_wallet   VARCHAR(64) NOT NULL,
    trader_wallet     VARCHAR(64) NOT NULL,
    copy_ratio        NUMERIC(6, 4) NOT NULL DEFAULT 0.1,
    max_trade         NUMERIC(38, 18) NOT NULL DEFAULT 500,
    daily_limit       NUMERIC(38, 18) NOT NULL DEFAULT 2000,
    enabled           BOOLEAN NOT NULL DEFAULT true,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (follower_wallet, trader_wallet)
);

CREATE INDEX IF NOT EXISTS idx_follows_trader ON follows (trader_wallet) WHERE enabled = true;
CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows (follower_wallet);