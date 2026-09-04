export interface Config {
  port: number;
  postgresUrl: string;
  kafkaBrokers: string[];
}

export function loadConfig(): Config {
  return {
    port: parseInt(process.env.PORT ?? "3003", 10),
    postgresUrl:
      process.env.POSTGRES_URL ??
      "postgres://fomox:fomox@127.0.0.1:5433/fomox?sslmode=disable",
    kafkaBrokers: (process.env.KAFKA_BROKERS ?? "127.0.0.1:19092").split(","),
  };
}