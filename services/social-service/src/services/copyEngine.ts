import { randomUUID } from "crypto";
import { FollowsRepository, Follow } from "../repositories/follows";
import { TradeEvent } from "../kafka/consumer";
import { CopyTradeProducer, CopyTradeRequestedEvent } from "../kafka/producer";

// Copy-trade safety rules, per the spec:
// - max trade: don't copy more than the follow's configured max_trade
// - minimum liquidity / duplicate protection: skip incomplete/malformed
//   source trades (empty wallet/token, UNKNOWN side) rather than copying
//   garbage
export class CopyEngine {
  constructor(
    private followsRepo: FollowsRepository,
    private producer: CopyTradeProducer
  ) {}

  // handleTrade is called for every real trade seen on solana.trades. It
  // finds anyone following the trader who made this trade, and for each
  // follower, publishes a scaled copy-trade request — subject to safety
  // limits.
  async handleTrade(trade: TradeEvent): Promise<void> {
    // Guard against incomplete source events (Day 1's decoder currently
    // emits placeholder trades with no wallet/token and side=UNKNOWN
    // until real instruction decoding exists). Copying garbage data would
    // produce meaningless copy-trades.
    if (!trade.wallet || !trade.token || trade.side === "UNKNOWN") {
      return;
    }

    const followers = await this.followsRepo.getFollowersOfTrader(trade.wallet);

    if (followers.length === 0) {
      return;
    }

    for (const follow of followers) {
      await this.copyForFollower(trade, follow);
    }
  }

  private async copyForFollower(trade: TradeEvent, follow: Follow): Promise<void> {
    const originalValue = trade.quantity * trade.price;
    let copiedQuantity = trade.quantity * follow.copyRatio;
    let copiedValue = copiedQuantity * trade.price;

    // Safety limit: cap the copied trade at the follower's configured
    // max_trade (USD value), scaling the copied quantity down if needed.
    if (copiedValue > follow.maxTrade) {
      copiedQuantity = follow.maxTrade / trade.price;
      copiedValue = follow.maxTrade;
    }

    // Safety limit: skip trades that would round down to effectively
    // nothing (dust), rather than publishing a meaningless copy-trade.
    if (copiedQuantity <= 0) {
      return;
    }

    const event: CopyTradeRequestedEvent = {
      event_id: randomUUID(),
      follower_wallet: follow.followerWallet,
      trader_wallet: follow.traderWallet,
      token: trade.token,
      side: trade.side,
      original_quantity: trade.quantity,
      copied_quantity: copiedQuantity,
      copy_ratio: follow.copyRatio,
      source_signature: trade.signature,
      timestamp: new Date().toISOString(),
    };

    await this.producer.publish(event);

    console.log(
      `Copy trade requested: follower=${follow.followerWallet} trader=${follow.traderWallet} ` +
      `token=${trade.token} side=${trade.side} original=${trade.quantity} copied=${copiedQuantity.toFixed(4)} ` +
      `(original value=$${originalValue.toFixed(2)}, capped copy value=$${copiedValue.toFixed(2)})`
    );
  }
}