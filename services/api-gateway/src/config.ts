export interface Config {
    port: number;
    postgresUrl: string;
  }
  
  export function loadConfig(): Config {
    return {
      port: parseInt(process.env.PORT ?? "3000", 10),
      postgresUrl:
        process.env.POSTGRES_URL ??
        "postgres://fomox:fomox@127.0.0.1:5433/fomox?sslmode=disable",
    };
  }