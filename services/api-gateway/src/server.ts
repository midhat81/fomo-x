import "dotenv/config";
import Fastify from "fastify";
import cors from "@fastify/cors";
import { Pool } from "pg";

import { loadConfig } from "./config";
import { PortfolioClient } from "./clients/portfolioClient";
import { healthRoutes } from "./routes/health";
import { walletRoutes } from "./routes/wallets";

async function main() {
  const config = loadConfig();

  const pool = new Pool({ connectionString: config.postgresUrl });

  // Fail fast if Postgres isn't reachable, rather than starting a server
  // that will error on every request.
  await pool.query("SELECT 1");
  console.log("Connected to Postgres.");

  const portfolioClient = new PortfolioClient(pool);

  const app = Fastify({ logger: true });

  await app.register(cors, { origin: true });

  await healthRoutes(app);
  await walletRoutes(app, portfolioClient);

  await app.listen({ port: config.port, host: "0.0.0.0" });
  console.log(`FOMO-X API Gateway listening on port ${config.port}`);
}

main().catch((err) => {
  console.error("Failed to start API Gateway:", err);
  process.exit(1);
});