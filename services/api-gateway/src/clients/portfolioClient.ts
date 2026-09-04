import { Pool } from "pg";

export interface Position {
  token: string;
  quantity: number;
  averageEntry: number;
}

export interface Trade {
  side: string;
  quantity: number;
  price: number;
  tradedAt: string;
  signature: string;
}

export interface PnLSummary {
  wallet: string;
  realizedPnL: number;
  unrealizedPnL: number;
  totalPnL: number;
}

export class PortfolioClient {
  constructor(private pool: Pool) {}

  async walletExists(address: string): Promise<boolean> {
    const result = await this.pool.query(
      "SELECT 1 FROM wallets WHERE address = $1",
      [address]
    );
    return (result.rowCount ?? 0) > 0;
  }

  async getPositions(wallet: string): Promise<Position[]> {
    const result = await this.pool.query(
      `SELECT token_address, quantity, average_entry
       FROM positions
       WHERE wallet_address = $1 AND quantity > 0`,
      [wallet]
    );

    return result.rows.map((row) => ({
      token: row.token_address,
      quantity: parseFloat(row.quantity),
      averageEntry: parseFloat(row.average_entry),
    }));
  }

  async getTrades(wallet: string, limit = 50): Promise<Trade[]> {
    const result = await this.pool.query(
      `SELECT side, quantity, price, traded_at, signature
       FROM trades
       WHERE wallet_address = $1
       ORDER BY traded_at DESC
       LIMIT $2`,
      [wallet, limit]
    );

    return result.rows.map((row) => ({
      side: row.side,
      quantity: parseFloat(row.quantity),
      price: parseFloat(row.price),
      tradedAt: row.traded_at,
      signature: row.signature,
    }));
  }

  // getPnL computes a simplified realized P&L: total SELL proceeds minus
  // total BUY cost, summed across all of the wallet's trade history. This
  // is an approximation (not true per-lot FIFO/LIFO accounting) — good
  // enough for Day 2's "paste a wallet, see P&L" deliverable. Unrealized
  // P&L is 0 until a live price feed exists (Day 5+ concern).
  async getPnL(wallet: string): Promise<PnLSummary> {
    const result = await this.pool.query(
      `SELECT
         COALESCE(SUM(CASE WHEN side = 'BUY' THEN quantity * price ELSE 0 END), 0) AS buy_value,
         COALESCE(SUM(CASE WHEN side = 'SELL' THEN quantity * price ELSE 0 END), 0) AS sell_value
       FROM trades
       WHERE wallet_address = $1`,
      [wallet]
    );

    const buyValue = parseFloat(result.rows[0].buy_value);
    const sellValue = parseFloat(result.rows[0].sell_value);
    const realizedPnL = sellValue - buyValue;
    const unrealizedPnL = 0;

    return {
      wallet,
      realizedPnL,
      unrealizedPnL,
      totalPnL: realizedPnL + unrealizedPnL,
    };
  }
}