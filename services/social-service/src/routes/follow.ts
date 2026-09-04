import { FastifyInstance } from "fastify";
import { FollowsRepository } from "../repositories/follows";

interface FollowBody {
  follower_wallet: string;
  copy_ratio?: number;
  max_trade?: number;
  daily_limit?: number;
}

export async function followRoutes(
  app: FastifyInstance,
  followsRepo: FollowsRepository
) {
  // POST /traders/:wallet/follow
  app.post("/traders/:wallet/follow", async (request, reply) => {
    const { wallet: traderWallet } = request.params as { wallet: string };
    const body = request.body as FollowBody;

    if (!body.follower_wallet) {
      return reply.status(400).send({ error: "follower_wallet is required" });
    }

    const follow = await followsRepo.follow(
      body.follower_wallet,
      traderWallet,
      body.copy_ratio,
      body.max_trade,
      body.daily_limit
    );

    return follow;
  });

  // DELETE /traders/:wallet/follow
  app.delete("/traders/:wallet/follow", async (request, reply) => {
    const { wallet: traderWallet } = request.params as { wallet: string };
    const body = request.body as { follower_wallet: string };

    if (!body?.follower_wallet) {
      return reply.status(400).send({ error: "follower_wallet is required" });
    }

    await followsRepo.unfollow(body.follower_wallet, traderWallet);
    return { unfollowed: true };
  });

  // GET /following?wallet=...
  app.get("/following", async (request, reply) => {
    const { wallet } = request.query as { wallet?: string };

    if (!wallet) {
      return reply.status(400).send({ error: "wallet query param is required" });
    }

    const following = await followsRepo.getFollowing(wallet);
    return { wallet, following };
  });
}