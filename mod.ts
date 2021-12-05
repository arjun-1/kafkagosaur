// @deno-types="./global.d.ts"
import "./lib/wasm_exec.js";
import { deadline, DeadlineError, delay } from "./deps.ts";
import { setOnGlobal as setConnectWithDeadlineOnGlobal } from "./connection-with-deadline.ts";
import { Dialer } from "./dialer.ts";
import { Reader, ReaderConfig } from "./reader.ts";
import { Writer, WriterConfig } from "./writer.ts";

const runGoWasm = async (wasmFilePath: string): Promise<unknown> => {
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
  const initialBackoffMs = 30;
  const maxDelayMs = 1000;

  const nextBackoffMs = (backoffMs: number): number =>
    initialBackoffMs + backoffMs;

  const loop = async (
    backoffMs: number = initialBackoffMs,
  ): Promise<unknown> => {
    const value = (global as Record<string, unknown>)[key];
    if (value !== undefined) return Promise.resolve(value);
    else {
      await delay(backoffMs);
      return loop(nextBackoffMs(backoffMs));
    }
  };

  return deadline(loop(), maxDelayMs).catch((e) => {
    if (e instanceof DeadlineError) {
      Promise.reject(`Global key ${key} undefined`);
    } else Promise.reject(e);
  });
};

class KafkaGoSaur {
  constructor() {
    setConnectWithDeadlineOnGlobal();
    runGoWasm("./bin/kafkagosaur.wasm");
  }

  async dialer(): Promise<Dialer> {
    const newDialer = await untilGloballyDefined(
      "newDialer",
    ) as () => Dialer;
    return newDialer();
  }

  async reader(config: ReaderConfig): Promise<Reader> {
    const newReader = await untilGloballyDefined(
      "newReader",
    ) as (config: ReaderConfig) => Reader;

    return newReader(config);
  }

  async writer(config: WriterConfig): Promise<Writer> {
    const newWriter = await untilGloballyDefined(
      "newWriter",
    ) as (config: WriterConfig) => Writer;

    return newWriter(config);
  }
}

export default KafkaGoSaur;
