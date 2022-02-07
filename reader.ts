import { Header } from "./header.ts";
import { SASLConfig } from "./security/sasl.ts";
import { TLSConfig } from "./security/tls.ts";

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
  groupId?: string;
  sasl?: SASLConfig;
  tls?: TLSConfig;
};

export interface KafkaReader {
  close: () => Promise<void>;
  commitMessages: (msgs: KafkaReadMessage[]) => Promise<void>;
  fetchMessage: () => Promise<KafkaReadMessage>;
  readMessage: () => Promise<KafkaReadMessage>;
  setOffset: (offset: number) => Promise<void>;
  setOffsetAt: (timeMs: number) => Promise<void>;
}
