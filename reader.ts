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
