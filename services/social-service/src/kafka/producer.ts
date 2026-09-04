import { Kafka, Producer } from "kafkajs";

export interface CopyTradeRequestedEvent {
  event_id: string;
  follower_wallet: string;
  trader_wallet: string;
  token: string;
  side: string;
  original_quantity: number;
  copied_quantity: number;
  copy_ratio: number;
  source_signature: string;
  timestamp: string;
}

export class CopyTradeProducer {
  private kafka: Kafka;
  private producer: Producer;

  constructor(brokers: string[]) {
    this.kafka = new Kafka({ clientId: "social-service", brokers });
    this.producer = this.kafka.producer();
  }

  async connect(): Promise<void> {
    await this.producer.connect();
  }

  async publish(event: CopyTradeRequestedEvent): Promise<void> {
    await this.producer.send({
      topic: "copy.trade.requested",
      messages: [{ key: event.event_id, value: JSON.stringify(event) }],
    });
  }

  async disconnect(): Promise<void> {
    await this.producer.disconnect();
  }
}