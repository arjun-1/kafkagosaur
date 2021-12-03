import { assert } from "../deps.ts";
import KafkaGoSaur from "../../mod.ts";
import { KafkaConn } from "../../dialer.ts";

const kafkaGoSaur = new KafkaGoSaur();
const kafkaBroker = "localhost:29092";

const dialer = await kafkaGoSaur.dialer();

const withKafkaConn = async <T>(
  resultFn: (conn: KafkaConn) => Promise<T>,
): Promise<T> => {
  const conn = await dialer.dial("tcp", kafkaBroker);

  const result = await resultFn(conn);
  await conn.close();
  return result;
};

Deno.test(
  "KafkaConn.apiVersions should list api versions",
  () =>
    withKafkaConn(async (conn: KafkaConn) => {
      const apiVersions = await conn.apiVersions();

      assert(apiVersions.length !== 0, "api versions was empty");
      assert(apiVersions[0].length !== 0, "api versions element was empty");
    }),
);
