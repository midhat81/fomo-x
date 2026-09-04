import { Kafka, Consumer } from "kafkajs";

export interface TradeEvent {
  event_id: string;
  wallet: string;
  signature: string;
  token: string;
  side: string;
  quantity: number;
  price: number;
  timestamp: string;
}

export class TradeConsumer {
  private kafka: Kafka;
  private consumer: Consumer;

  constructor(brokers: string[]) {
    this.kafka = new Kafka({ clientId: "social-service", brokers });
    this.consumer = this.kafka.consumer({ groupId: "social-service" });
  }

  async connect(): Promise<void> {
    await this.consumer.connect();
    await this.consumer.subscribe({ topic: "solana.trades", fromBeginning: false });
  }

  async run(onTrade: (event: TradeEvent) => Promise<void>): Promise<void> {
    await this.consumer.run({
      eachMessage: async ({ message }) => {
        if (!message.value) return;

        try {
          const event: TradeEvent = JSON.parse(message.value.toString());
          await onTrade(event);
        } catch (err) {
          console.error("failed to process trade event:", err);
        }
      },
    });
  }

  async disconnect(): Promise<void> {
    await this.consumer.disconnect();
  }
}