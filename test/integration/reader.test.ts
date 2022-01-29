import { assertEquals } from "../deps.ts";
import {
  kafkaGoSaur,
  readerConfig as readerConfigNoSASL,
  readerConfigSASL,
  withWriter,
  writerConfig as writerConfigNoSASL,
  writerConfigSASL,
} from "./setup.ts";
import { KafkaReader, KafkaReaderConfig } from "../../reader.ts";
import { KafkaWriteMessage, KafkaWriterConfig } from "../../writer.ts";

const withReader = (config: KafkaReaderConfig) =>
  async <T>(
    resultFn: (reader: KafkaReader) => Promise<T>,
  ): Promise<T> => {
    const reader = await kafkaGoSaur.reader(config);
    const result = await resultFn(reader);

    await reader.close();
    return result;
  };

const withTestSetup = (
  writerConfig: KafkaWriterConfig = writerConfigNoSASL,
  readerConfig: KafkaReaderConfig = readerConfigNoSASL,
) =>
  <T>(msg: KafkaWriteMessage) =>
    (
      resultFn: (reader: KafkaReader) => Promise<T>,
    ): Promise<T> =>
      withWriter(writerConfig)(async (writer) => {
        await writer.writeMessages([msg]);

        return withReader(readerConfig)(resultFn);
      });

const withTestSetupSASL = withTestSetup(writerConfigSASL, readerConfigSASL);

Deno.test(
  "Reader.readMessage should read a message",
  async () => {
    const enc = new TextEncoder();
    const dec = new TextDecoder();

    const msgTxt = "value";
    const msg = {
      value: enc.encode(msgTxt),
    };

    await withTestSetup()(msg)(async (reader) => {
      const readMsg = await reader.readMessage();
      const readValue = dec.decode(readMsg.value);
      return assertEquals(readValue, msgTxt);
    });
  },
);

Deno.test(
  "Reader.readMessage should read a message (SASL)",
  async () => {
    const enc = new TextEncoder();
    const dec = new TextDecoder();

    const msgTxt = "value";
    const msg = {
      value: enc.encode(msgTxt),
    };

    await withTestSetupSASL(msg)(async (reader) => {
      const readMsg = await reader.readMessage();
      const readValue = dec.decode(readMsg.value);
      return assertEquals(readValue, msgTxt);
    });
  },
);
