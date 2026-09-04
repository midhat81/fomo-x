import { Pool } from "pg";

export interface Follow {
  id: string;
  followerWallet: string;
  traderWallet: string;
  copyRatio: number;
  maxTrade: number;
  dailyLimit: number;
  enabled: boolean;
}

export class FollowsRepository {
  constructor(private pool: Pool) {}

  async follow(
    followerWallet: string,
    traderWallet: string,
    copyRatio = 0.1,
    maxTrade = 500,
    dailyLimit = 2000
  ): Promise<Follow> {
    const result = await this.pool.query(
      `INSERT INTO follows (follower_wallet, trader_wallet, copy_ratio, max_trade, daily_limit, enabled)
       VALUES ($1, $2, $3, $4, $5, true)
       ON CONFLICT (follower_wallet, trader_wallet)
       DO UPDATE SET enabled = true, copy_ratio = $3, max_trade = $4, daily_limit = $5, updated_at = now()
       RETURNING id, follower_wallet, trader_wallet, copy_ratio, max_trade, daily_limit, enabled`,
      [followerWallet, traderWallet, copyRatio, maxTrade, dailyLimit]
    );

    return this.mapRow(result.rows[0]);
  }

  async unfollow(followerWallet: string, traderWallet: string): Promise<void> {
    await this.pool.query(
      `UPDATE follows SET enabled = false, updated_at = now()
       WHERE follower_wallet = $1 AND trader_wallet = $2`,
      [followerWallet, traderWallet]
    );
  }

  async getFollowing(followerWallet: string): Promise<Follow[]> {
    const result = await this.pool.query(
      `SELECT id, follower_wallet, trader_wallet, copy_ratio, max_trade, daily_limit, enabled
       FROM follows
       WHERE follower_wallet = $1 AND enabled = true`,
      [followerWallet]
    );
    return result.rows.map(this.mapRow);
  }

  // getFollowersOfTrader is the critical query for the copy-trading
  // pipeline: given a trader who just made a trade, who is watching them?
  async getFollowersOfTrader(traderWallet: string): Promise<Follow[]> {
    const result = await this.pool.query(
      `SELECT id, follower_wallet, trader_wallet, copy_ratio, max_trade, daily_limit, enabled
       FROM follows
       WHERE trader_wallet = $1 AND enabled = true`,
      [traderWallet]
    );
    return result.rows.map(this.mapRow);
  }

  private mapRow(row: any): Follow {
    return {
      id: row.id,
      followerWallet: row.follower_wallet,
      traderWallet: row.trader_wallet,
      copyRatio: parseFloat(row.copy_ratio),
      maxTrade: parseFloat(row.max_trade),
      dailyLimit: parseFloat(row.daily_limit),
      enabled: row.enabled,
    };
  }
}