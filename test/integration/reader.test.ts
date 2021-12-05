import KafkaGoSaur from "../../mod.ts";
import { assertEquals } from "../deps.ts";
import { delay } from "../../deps.ts";
import { Reader } from "../../reader.ts";
import { withWriter } from "./writer.test.ts";

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
const withProducedMessage = async <T>(
  resultFn: (reader: Reader) => Promise<T>,
): Promise<T> => {
  withWriter;
  const writer = await kafkaGoSaur.writer(writerConfig);
  const reader = await kafkaGoSaur.reader(readerConfig);

  const enc = new TextEncoder();
  const msgs = [{ value: enc.encode("value") }];

  await writer.writeMessages(msgs);

  await writer.close();
  const result = await resultFn(reader);

  await delay(2 * writerConfig.idleTimeout);
  return result;
};

Deno.test(
  "Reader.readMessage should read a message",
  () =>
    withProducedMessage(async (reader: Reader) => {
      const message = await reader.readMessage();

      assertEquals(message, 1);
    }),
);
