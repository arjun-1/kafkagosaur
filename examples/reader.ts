import KafkaGoSaur from "https://deno.land/x/kafkagosaur@v0.0.6/mod.ts";
import { SASLMechanism } from "https://deno.land/x/kafkagosaur@v0.0.6/security/sasl.ts";

const broker = "localhost:9093";
const topic = "test-0";

const config = {
  brokers: [broker],
  topic,
  groupId: "group-id",
  sasl: {
    mechanism: SASLMechanism.SCRAMSHA512,
    username: "adminscram",
    password: "admin-secret-512",
  },
};

const kafkaGoSaur = new KafkaGoSaur();
const reader = await kafkaGoSaur.createReader(config);

const dec = new TextDecoder();
const readMsg = await reader.readMessage();
const readValue = dec.decode(readMsg.value);
console.log(`Read message ${readValue} from topic ${topic} at ${broker}`);

await reader.close();
