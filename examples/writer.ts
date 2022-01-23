import KafkaGoSaur from "https://github.com/arjun-1/kafkagosaur/mod.ts";
import { SASLMechanism } from "https://github.com/arjun-1/kafkagosaur/mod.ts/sasl.ts";

const broker = "localhost:9093";
const topic = "test-0";

const writerConfig = {
  address: broker,
  topic,
  idleTimeout: 10,
  sasl: {
    mechanism: SASLMechanism.SCRAMSHA512,
    username: "adminscram",
    password: "admin-secret-512",
  }
};

const kafkaGoSaur = new KafkaGoSaur();
const writer = await kafkaGoSaur.writer(writerConfig);

const enc = new TextEncoder();
const msgs = [{ value: enc.encode("value") }];

await writer.writeMessages(msgs);
