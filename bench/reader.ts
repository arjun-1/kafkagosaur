import KafkaGoSaur from "../mod.ts";
import { KafkaReaderConfig } from "../reader.ts";
import { SASLMechanism } from "../security/sasl.ts";
import { bench, runBenchmarks } from "./deps.ts";
import { broker, password, topic, username } from "./config.ts";

const nrOfRuns = 10;
const nrOfMessages = 10000;

const readerConfig: KafkaReaderConfig = {
  brokers: [broker],
  topic,
  sasl: {
    mechanism: SASLMechanism.PLAIN,
    username,
    password,
  },
  tls: {
    insecureSkipVerify: true,
  },
};

const kafkaGoSaur = new KafkaGoSaur();
const reader = await kafkaGoSaur.createReader(readerConfig);

bench({
  name: `readMessage#${nrOfMessages}`,
  runs: nrOfRuns,
  async func(b): Promise<void> {
    b.start();

    for (let i = 0; i < nrOfMessages; i++) {
      await reader.readMessage();
    }

    b.stop();
  },
});

const benchmarkRunResults = await runBenchmarks();

console.log(benchmarkRunResults);
console.log(
  `[kafkagosaur] readMessage msgs/s: ${
    nrOfMessages / (benchmarkRunResults.results[0].measuredRunsAvgMs / 1000)
  }`,
);

await reader.close();
