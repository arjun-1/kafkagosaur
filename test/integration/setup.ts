import { delay } from "../../deps.ts";
import KafkaGoSaur from "../../mod.ts";

const kafkaGoSaur = new KafkaGoSaur();
await delay(50);

const kafkaBroker = "localhost:29092";
const writerConfig = {
  address: kafkaBroker,
  topic: "my-topic",
  idleTimeout: 10,
};

const readerConfig = {
  brokers: [writerConfig.address],
  topic: writerConfig.topic,
  groupId: "group-id",
};

export { kafkaBroker, kafkaGoSaur, readerConfig, writerConfig };
