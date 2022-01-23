import { assert } from "../deps.ts";
import {
  broker as brokerNoSASL,
  brokerSASL,
  kafkaGoSaur,
  saslConfig,
} from "./setup.ts";
import { KafkaConn, KafkaDialerConfig } from "../../dialer.ts";

const withKafkaConn = (
  broker: string = brokerNoSASL,
  config: KafkaDialerConfig = {},
) =>
  async <T>(
    resultFn: (conn: KafkaConn) => Promise<T>,
  ): Promise<T> => {
    const dialer = await kafkaGoSaur.dialer(config);
    const conn = await dialer.dial("tcp", broker);

    const result = await resultFn(conn);
    await conn.close();

    return result;
  };

const withKafkaConnSASL = withKafkaConn(brokerSASL, { sasl: saslConfig });

Deno.test(
  "KafkaConn.apiVersions should list api versions",
  () =>
    withKafkaConn()(async (conn: KafkaConn) => {
      const apiVersions = await conn.apiVersions();

      assert(apiVersions.length !== 0, "api versions was empty");
      assert(apiVersions[0].length !== 0, "api versions element was empty");
    }),
);

Deno.test(
  "KafkaConn.apiVersions should list api versions (SASL)",
  () =>
    withKafkaConnSASL(async (conn: KafkaConn) => {
      const apiVersions = await conn.apiVersions();

      assert(apiVersions.length !== 0, "api versions was empty");
      assert(apiVersions[0].length !== 0, "api versions element was empty");
    }),
);
