/**
 * A Kafka client for Deno built using WebAssembly.
 *
 * ### Example reader
 *
 * ```typescript
 * const kafkaGoSaur = new KafkaGoSaur();
 * const reader = await kafkaGoSaur.createReader({
 *  brokers: ["localhost:9092"],
 *  topic: "test-0",
 * });
 *
 * const readMsg = await reader.readMessage();
 * ```
 *
 * ### Example writer
 *
 * ```typescript
 * const kafkaGoSaur = new KafkaGoSaur();
 * const writer = await kafkaGoSaur.createWriter({
 *  broker: "localhost:9092",
 *  topic: "test-0",
 * });
 *
 * const enc = new TextEncoder();
 * const msgs = [{ value: enc.encode("value") }];
 *
 * await writer.writeMessages(msgs);
 * ```
 *
 * @module
 */

// @deno-types="./global.d.ts"
import "./lib/wasm_exec.js";
import { deadline, DeadlineError, delay } from "./deps.ts";
import { dial as dialNode } from "./net/node-connection.ts";
import { dial as dialDeno } from "./net/deno-connection.ts";
import { KafkaDialer, KafkaDialerConfig } from "./dialer.ts";
import { KafkaReader, KafkaReaderConfig } from "./reader.ts";
import { KafkaWriter, KafkaWriterConfig } from "./writer.ts";

const runGoWasm = async (wasmFilePath: string): Promise<void> => {
  const go = new global.Go();
  const wasmBytes = await Deno.readFile(wasmFilePath);
  const instiatedSource = await WebAssembly.instantiate(
    wasmBytes,
    go.importObject,
  );

  return go.run(instiatedSource.instance);
};

const untilGloballyDefined = (
  key: string,
): Promise<unknown> => {
  const backoffMs = 30;
  const maxDelayMs = 1000;

  const globalGoInteropKey = "kafkagosaur.go";

  const loop = async (): Promise<unknown> => {
    const interopObject =
      (globalThis as Record<string, unknown>)[globalGoInteropKey];

    if (interopObject !== undefined) {
      return Promise.resolve((interopObject as Record<string, unknown>)[key]);
    } else {
      await delay(backoffMs);
      return loop();
    }
  };

  return deadline(loop(), maxDelayMs).catch((e) => {
    if (e instanceof DeadlineError) {
      Promise.reject(`Global key ${key} undefined`);
    } else Promise.reject(e);
  });
};

export const setOnGlobal = (value: unknown) => {
  const globalDenoInteropKey = "kafkagosaur.deno";

  (globalThis as Record<string, unknown>)[globalDenoInteropKey] = value;
};

class KafkaGoSaur {
  constructor() {
    setOnGlobal({ dialNode, dialDeno });
    runGoWasm("./bin/kafkagosaur.wasm");
  }

  async createDialer(config: KafkaDialerConfig): Promise<KafkaDialer> {
    const create = await untilGloballyDefined(
      "createDialer",
    ) as (config: KafkaDialerConfig) => KafkaDialer;
    return create(config);
  }

  async createReader(config: KafkaReaderConfig): Promise<KafkaReader> {
    const create = await untilGloballyDefined(
      "createReader",
    ) as (config: KafkaReaderConfig) => KafkaReader;

    return create(config);
  }

  async createWriter(config: KafkaWriterConfig): Promise<KafkaWriter> {
    const create = await untilGloballyDefined(
      "createWriter",
    ) as (config: KafkaWriterConfig) => KafkaWriter;

    return create(config);
  }
}

export default KafkaGoSaur;
