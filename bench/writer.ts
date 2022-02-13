import KafkaGoSaur from "../mod.ts";
import { KafkaWriterConfig } from "../writer.ts";
import { SASLMechanism } from "../security/sasl.ts";
import { bench, runBenchmarks } from "./deps.ts";
import { broker, password, topic, username } from "./config.ts";

const nrOfRuns = 10;
const nrOfMessages = 100000;
const msgBatchSize = 10000;
const msgSize = 1024;

const writerConfig: KafkaWriterConfig = {
  address: broker,
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
const writer = await kafkaGoSaur.createWriter(writerConfig);

const value = new Uint8Array(msgSize);
crypto.getRandomValues(value);

bench({
  name: `writeMessages#${nrOfMessages}`,
  runs: nrOfRuns,
  async func(b): Promise<void> {
    b.start();

    for (let i = 0; i < nrOfMessages / msgBatchSize; i++) {
      const msgs = [];

      for (let j = 0; j < msgBatchSize; j++) {
        msgs.push({ value });
      }
      await writer.writeMessages(msgs);
    }

    b.stop();
  },
});

const benchmarkRunResults = await runBenchmarks();

console.log(benchmarkRunResults);
console.log(
  `[kafkagosaur] writeMessages msgs/s: ${
    nrOfMessages / (benchmarkRunResults.results[0].measuredRunsAvgMs / 1000)
  }`,
);

await writer.close();
