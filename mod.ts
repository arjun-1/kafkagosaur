// @deno-types="./global.d.ts"
import "./lib/wasm_exec.js";
import { deadline, DeadlineError, delay } from "./deps.ts";
import { Dial, setDialOnGlobal } from "./net/connection.ts";
import { dial as nodeDial } from "./net/node-connection.ts";
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

  const loop = async (): Promise<unknown> => {
    const value = (globalThis as Record<string, unknown>)[key];
    if (value !== undefined) return Promise.resolve(value);
    else {
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

class KafkaGoSaur {
  constructor(dial: Dial = nodeDial) {
    setDialOnGlobal(dial);
    runGoWasm("./bin/kafkagosaur.wasm");
  }

  async dialer(config: KafkaDialerConfig): Promise<KafkaDialer> {
    const newDialer = await untilGloballyDefined(
      "newDialer",
    ) as (config: KafkaDialerConfig) => KafkaDialer;
    return newDialer(config);
  }

  async reader(config: KafkaReaderConfig): Promise<KafkaReader> {
    const newReader = await untilGloballyDefined(
      "newReader",
    ) as (config: KafkaReaderConfig) => KafkaReader;

    return newReader(config);
  }

  async writer(config: KafkaWriterConfig): Promise<KafkaWriter> {
    const newWriter = await untilGloballyDefined(
      "newWriter",
    ) as (config: KafkaWriterConfig) => KafkaWriter;

    return newWriter(config);
  }
}

export default KafkaGoSaur;
