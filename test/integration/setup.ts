import { delay } from "../../deps.ts";
import KafkaGoSaur from "../../mod.ts";
import { Writer } from "../../writer.ts";

const kafkaGoSaur = new KafkaGoSaur();
await delay(50);

const kafkaBroker = "localhost:9092";

const writerConfig = {
  address: kafkaBroker,
  topic: "test-0",
  idleTimeout: 10,
};

const readerConfig = {
  brokers: [writerConfig.address],
  topic: writerConfig.topic,
  groupId: "group-id",
};

const withWriter = async <T>(
  resultFn: (writer: Writer) => Promise<T>,
): Promise<T> => {
  const writer = await kafkaGoSaur.writer(writerConfig);

  const result = await resultFn(writer);
  await writer.close();

  await delay(2 * writerConfig.idleTimeout);
  return result;
};

export { kafkaBroker, kafkaGoSaur, readerConfig, writerConfig, withWriter };
