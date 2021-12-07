import { kafkaGoSaur, writerConfig } from "./setup.ts";
import { delay } from "../../deps.ts";
import { Writer } from "../../writer.ts";

export const withWriter = async <T>(
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
