// @deno-types="./global.d.ts"
import "./lib/wasm_exec.js";
import { setOnGlobal as setConnectWithDeadlineOnGlobal } from "./connection-with-deadline.ts";

export interface KafkaConn {
  apiVersions: () => Promise<string[]>;
  close: () => Promise<void>;
}

export interface Dialer {
  dial: (network: "tcp", address: string) => Promise<KafkaConn>;
}

export type MessageRead = {
  topic: string;
  partition: number;
  time: number;
  key: Uint8Array;
  value: Uint8Array;
};

export type ReaderConfig = {
  brokers: string[];
  topic: string;
  groupId: string;
};

export interface Reader {
  readMessage: () => Promise<MessageRead>;
  close: () => Promise<void>;
}

export type MessageWrite = {
  topic?: string;
  time?: number;
  key?: Uint8Array;
  value: Uint8Array;
};

export type WriterConfig = {
  topic?: string;
  address: string;
};

export interface Writer {
  writeMessages: (msgs: MessageWrite[]) => Promise<null>;
  close: () => Promise<void>;
}

const runGoWasm = async (wasmFilePath: string): Promise<unknown> => {
  const go = new global.Go();
  const wasmBytes = await Deno.readFile(wasmFilePath);
  const instiatedSource = await WebAssembly.instantiate(
    wasmBytes,
    go.importObject,
  );

  return go.run(instiatedSource.instance);
};

const delay = (ms: number): Promise<unknown> =>
  new Promise((resolve) => setTimeout(resolve, ms));

const nextBackoffMs = (backoffMs: number): number =>
  initialBackoffMs + backoffMs;

const initialBackoffMs = 30;
const maxDelayMs = 1000;

const untilGloballyDefined = async (
  key: string,
  backoffMs: number = initialBackoffMs,
): Promise<unknown> => {
  if (backoffMs >= maxDelayMs) {
    return Promise.reject(`Global key ${key} undefined`);
  }

  const value = (global as Record<string, unknown>)[key];
  if (value !== undefined) return Promise.resolve(value);
  else {
    await delay(backoffMs);
    return untilGloballyDefined(key, nextBackoffMs(backoffMs));
  }
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
