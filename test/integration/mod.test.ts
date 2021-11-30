import { assert } from "../deps.ts";
import KafkaGoSaur, { KafkaConn, Writer } from "../../mod.ts";

const kafkaGoSaur = new KafkaGoSaur();
const kafkaBroker = "localhost:29092";
const writerConfig = {
  address: kafkaBroker,
  topic: "my-topic",
};

const dialer = await kafkaGoSaur.dialer();

const withKafkaConn = async <T>(
  resultFn: (conn: KafkaConn) => Promise<T>,
): Promise<T> => {
  const conn = await dialer.dial("tcp", kafkaBroker);
  const result = await resultFn(conn);
  await conn.close();
  return result;
};

const withWriter = async <T>(
  resultFn: (writer: Writer) => Promise<T>,
): Promise<T> => {
  const writer = await kafkaGoSaur.writer(writerConfig);
  const result = await resultFn(writer);
  await writer.close();
  return result;
};

// Conn

Deno.test(
  "KafkaConn.apiVersions should list api versions",
  () =>
    withKafkaConn(async (conn: KafkaConn) => {
      const apiVersions = await conn.apiVersions();

      assert(apiVersions.length !== 0, "api versions was empty");
      assert(apiVersions[0].length !== 0, "api versions element was empty");
    }),
);

Deno.test(
  "Writer.writeMessages should write messages",
  () =>
    withWriter(async (writer: Writer) => {
      const enc = new TextEncoder();
      const msgs = [{ value: enc.encode("hoihoi") }];

      const result = await writer.writeMessages(msgs);
      await writer.close();
      console.log(result);
    }),
);
