import KafkaGoSaur from "https://raw.githubusercontent.com/arjun-1/kafkagosaur/master/mod.ts";
import { SASLMechanism } from "https://raw.githubusercontent.com/arjun-1/kafkagosaur/master/sasl.ts";

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
  },
};

const kafkaGoSaur = new KafkaGoSaur();
const writer = await kafkaGoSaur.writer(writerConfig);

const enc = new TextEncoder();
const msgs = [{ value: enc.encode("value0") }, { value: enc.encode("value1") }];

await writer.writeMessages(msgs);
await writer.close();
