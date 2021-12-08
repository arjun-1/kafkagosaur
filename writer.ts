import { Header } from "./header.ts";

export type MessageWrite = {
  topic?: string;
  offset?: number;
  highWaterMark?: number;
  key?: Uint8Array;
  value?: Uint8Array;
  headers?: Header[];
  time?: number;
};

export type WriterConfig = {
  topic?: string;
  address: string;
  idleTimeout?: number;
};

export interface Writer {
  writeMessages: (msgs: MessageWrite[]) => Promise<null>;
  close: () => Promise<void>;
}
