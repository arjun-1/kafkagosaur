import { Header } from "./header.ts";
import { SASLConfig } from "./sasl.ts";
export type KafkaWriteMessage = {
  topic?: string;
  offset?: number;
  highWaterMark?: number;
  key?: Uint8Array;
  value?: Uint8Array;
  headers?: Header[];
  time?: number;
};

export type KafkaWriterConfig = {
  topic?: string;
  address: string;
  idleTimeout?: number;
  sasl?: SASLConfig;
};

export interface KafkaWriter {
  writeMessages: (msgs: KafkaWriteMessage[]) => Promise<null>;
  close: () => Promise<null>;
}
