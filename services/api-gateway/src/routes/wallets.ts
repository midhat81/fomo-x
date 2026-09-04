import { FastifyInstance } from "fastify";
import { PortfolioClient } from "../clients/portfolioClient";

export async function walletRoutes(
  app: FastifyInstance,
  portfolioClient: PortfolioClient
) {
  app.get("/wallet/:address", async (request, reply) => {
    const { address } = request.params as { address: string };

    const exists = await portfolioClient.walletExists(address);
    if (!exists) {
      return reply.status(404).send({ error: "wallet not found" });
    }

    return { wallet: address };
  });

  app.get("/wallet/:address/portfolio", async (request, reply) => {
    const { address } = request.params as { address: string };

    const exists = await portfolioClient.walletExists(address);
    if (!exists) {
      return reply.status(404).send({ error: "wallet not found" });
    }

    const positions = await portfolioClient.getPositions(address);
    return { wallet: address, positions };
  });

  app.get("/wallet/:address/trades", async (request, reply) => {
    const { address } = request.params as { address: string };

    const exists = await portfolioClient.walletExists(address);
    if (!exists) {
      return reply.status(404).send({ error: "wallet not found" });
    }

    const trades = await portfolioClient.getTrades(address);
    return { wallet: address, trades };
  });

  app.get("/wallet/:address/pnl", async (request, reply) => {
    const { address } = request.params as { address: string };

    const exists = await portfolioClient.walletExists(address);
    if (!exists) {
      return reply.status(404).send({ error: "wallet not found" });
    }

    const pnl = await portfolioClient.getPnL(address);
    return pnl;
  });
}