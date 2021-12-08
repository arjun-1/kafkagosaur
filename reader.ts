import { Header } from "./header.ts";

export type MessageRead = {
  topic: string;
  partition: number;
  offset: number;
  highWaterMark: number;
  key: Uint8Array;
  value: Uint8Array;
  headers: Header[];
  time: number;
};

export type ReaderConfig = {
  brokers: string[];
  topic: string;
  groupId: string;
};

export interface Reader {
  close: () => Promise<void>;
  commitMessages: (msgs: MessageRead[]) => Promise<void>;
  fetchMessage: () => Promise<MessageRead>;
  readMessage: () => Promise<MessageRead>;
  setOffset: (offset: number) => Promise<void>;
  setOffsetAt: (timeMs: number) => Promise<void>
}
