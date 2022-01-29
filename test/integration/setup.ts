import { delay } from "../../deps.ts";
import KafkaGoSaur from "../../mod.ts";
import { KafkaWriter, KafkaWriterConfig } from "../../writer.ts";
import { SASLMechanism } from "../../sasl.ts";

const kafkaGoSaur = new KafkaGoSaur();
await delay(50);

const broker = "localhost:9092";
const brokerSASL = "localhost:9093";
const topic = "test-0";

const saslConfig = {
  mechanism: SASLMechanism.SCRAMSHA512,
  username: "adminscram",
  password: "admin-secret-512",
};

const writerConfig = {
  address: broker,
  topic,
  idleTimeout: 10,
};

const writerConfigSASL = {
  ...writerConfig,
  address: brokerSASL,
  sasl: saslConfig,
};

const readerConfig = {
  brokers: [broker],
  topic,
  groupId: "group-id",
};

const readerConfigSASL = {
  ...readerConfig,
  brokers: [brokerSASL],
  sasl: saslConfig,
};

const withWriter = (config: KafkaWriterConfig = writerConfig) =>
  async <T>(
    resultFn: (writer: KafkaWriter) => Promise<T>,
  ): Promise<T> => {
    const writer = await kafkaGoSaur.writer(config);

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
