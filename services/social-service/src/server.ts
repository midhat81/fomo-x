import "dotenv/config";
import Fastify from "fastify";
import cors from "@fastify/cors";
import { Pool } from "pg";

import { loadConfig } from "./config";
import { FollowsRepository } from "./repositories/follows";
import { TradeConsumer } from "./kafka/consumer";
import { CopyTradeProducer } from "./kafka/producer";
import { CopyEngine } from "./services/copyEngine";
import { healthRoutes } from "./routes/health";
import { followRoutes } from "./routes/follow";

async function main() {
  const config = loadConfig();

  const pool = new Pool({ connectionString: config.postgresUrl });
  await pool.query("SELECT 1");
  console.log("Connected to Postgres.");

  const followsRepo = new FollowsRepository(pool);

  const tradeConsumer = new TradeConsumer(config.kafkaBrokers);
  const copyProducer = new CopyTradeProducer(config.kafkaBrokers);

  await tradeConsumer.connect();
  await copyProducer.connect();
  console.log("Connected to Kafka (consumer + producer).");

  const copyEngine = new CopyEngine(followsRepo, copyProducer);

  // Run the Kafka consume loop in the background — it doesn't block the
  // HTTP server from starting.
  tradeConsumer.run(async (trade) => {
    await copyEngine.handleTrade(trade);
  });

  const app = Fastify({ logger: false });
  await app.register(cors, { origin: true });

  await healthRoutes(app);
  await followRoutes(app, followsRepo);

  await app.listen({ port: config.port, host: "0.0.0.0" });
  console.log(`FOMO-X Social Service listening on port ${config.port}`);
}

main().catch((err) => {
  console.error("Failed to start Social Service:", err);
  process.exit(1);
});