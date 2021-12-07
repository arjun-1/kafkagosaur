import { assertEquals } from "../deps.ts";
import { kafkaGoSaur, readerConfig } from "./setup.ts";
import { Reader } from "../../reader.ts";
import { MessageWrite } from "../../writer.ts";
import { withWriter } from "./writer.test.ts";

const withReader = async <T>(
  resultFn: (reader: Reader) => Promise<T>,
): Promise<T> => {
  const reader = await kafkaGoSaur.reader(readerConfig);
  const result = await resultFn(reader);

  await reader.close();
  return result;
};

const withTestSetup = <T>(msg: MessageWrite) =>
  (
    resultFn: (reader: Reader) => Promise<T>,
  ): Promise<T> =>
    withWriter(async (writer) => {
      await writer.writeMessages([msg]);

      return withReader(resultFn);
    });

Deno.test(
  "Reader.readMessage should read a message",
  async () => {
    const enc = new TextEncoder();
    const dec = new TextDecoder();

    const msgTxt = "value";
    const msg = {
      value: enc.encode(msgTxt),
    };

    await withTestSetup(msg)(async (reader) => {
      const readMsg = await reader.readMessage();
      const readValue = dec.decode(readMsg.value);
      return assertEquals(readValue, msgTxt);
    });
  },
);
