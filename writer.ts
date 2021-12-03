export type MessageWrite = {
  topic?: string;
  time?: number;
  key?: Uint8Array;
  value: Uint8Array;
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
