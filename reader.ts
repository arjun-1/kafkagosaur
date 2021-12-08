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
  readMessage: () => Promise<MessageRead>;
  close: () => Promise<void>;
}
