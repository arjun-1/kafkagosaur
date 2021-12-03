import KafkaGoSaur, { delay } from "../../mod.ts";
import { Writer } from "../../writer.ts";

const kafkaGoSaur = new KafkaGoSaur();
await delay(50);

const kafkaBroker = "localhost:29092";
const writerConfig = {
  address: kafkaBroker,
  topic: "my-topic",
  idleTimeout: 10,
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

Deno.test(
  "Writer.writeMessages should write messages",
  () =>
    withWriter(async (writer: Writer) => {
      const enc = new TextEncoder();
      const msgs = [{ value: enc.encode("value") }];

      await writer.writeMessages(msgs);
    }),
);
