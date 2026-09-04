const { Kafka } = require("kafkajs");
const { randomUUID } = require("crypto");

async function main() {
  const kafka = new Kafka({ clientId: "whale-test-publisher", brokers: ["127.0.0.1:19092"] });
  const producer = kafka.producer();
  await producer.connect();

  const event = {
    event_id: randomUUID(),
    wallet: "WhaleWa11et1111111111111111111111111111111",
    signature: "whale-test-sig-" + randomUUID(),
    token: "SomeToken1111111111111111111111111111111111",
    side: "BUY",
    quantity: 50,
    price: 100,
    timestamp: new Date().toISOString(),
  };

  await producer.send({
    topic: "solana.trades",
    messages: [{ key: event.event_id, value: JSON.stringify(event) }],
  });

  console.log("Published whale trade:", event);
  await producer.disconnect();
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});