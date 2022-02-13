import { delay } from "../../deps.ts";
import KafkaGoSaur from "../../mod.ts";
import { KafkaWriter, KafkaWriterConfig } from "../../writer.ts";
import { KafkaReaderConfig } from "../../reader.ts";
import { SASLMechanism } from "../../security/sasl.ts";

const kafkaGoSaur = new KafkaGoSaur();
// Ensure promise to instantiate wasm is awaited, so no async ops are leaked.
await delay(50);

const broker = "localhost:9092";
const brokerSASL = "localhost:9093";
const topic = "test-0";

const saslConfig = {
  mechanism: SASLMechanism.SCRAMSHA512,
  username: "adminscram",
  password: "admin-secret-512",
};

const writerConfig: KafkaWriterConfig = {
  address: broker,
  topic,
  idleTimeout: 10,
  logger: false,
};

const writerConfigSASL: KafkaWriterConfig = {
  ...writerConfig,
  address: brokerSASL,
  sasl: saslConfig,
  logger: false,
};

const readerConfig: KafkaReaderConfig = {
  brokers: [broker],
  topic,
  groupId: "group-id",
  logger: false,
};

const readerConfigSASL: KafkaReaderConfig = {
  ...readerConfig,
  brokers: [brokerSASL],
  sasl: saslConfig,
  logger: false,
};

const withWriter = (config: KafkaWriterConfig = writerConfig) =>
  async <T>(
    resultFn: (writer: KafkaWriter) => Promise<T>,
  ): Promise<T> => {
    const writer = await kafkaGoSaur.createWriter(config);

    const result = await resultFn(writer);
    await writer.close();

    if (config.idleTimeout) await delay(3 * config.idleTimeout);
    return result;
  };

const withWriterSASL = withWriter(writerConfigSASL);

export {
  broker,
  brokerSASL,
  kafkaGoSaur,
  readerConfig,
  readerConfigSASL,
  saslConfig,
  withWriter,
  withWriterSASL,
  writerConfig,
  writerConfigSASL,
};
