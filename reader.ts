import { Header } from "./header.ts";
import { SASLConfig } from "./sasl.ts";

export type KafkaReadMessage = {
  topic: string;
  partition: number;
  offset: number;
  highWaterMark: number;
  key: Uint8Array;
  value: Uint8Array;
  headers: Header[];
  time: number;
};

export type KafkaReaderConfig = {
  brokers: string[];
  topic: string;
  groupId: string;
  sasl?: SASLConfig;
};

export interface KafkaReader {
  close: () => Promise<null>;
  // TODO: is actually null!
  commitMessages: (msgs: KafkaReadMessage[]) => Promise<void>;
  fetchMessage: () => Promise<KafkaReadMessage>;
  readMessage: () => Promise<KafkaReadMessage>;
  setOffset: (offset: number) => Promise<null>;
  setOffsetAt: (timeMs: number) => Promise<null>;
}
