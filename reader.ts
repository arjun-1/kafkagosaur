import { Header } from "./header.ts";
import { SASLConfig } from "./security/sasl.ts";
import { TLSConfig } from "./security/tls.ts";

/** A data structure representing kafka messages read from `KafkaReader`. */
export type KafkaReadMessage = {
  /** Topic indicates which topic this message was consumed from via `KafkaReader`. */
  topic: string;
  partition: number;
  offset: number;
  highWaterMark: number;
  key: Uint8Array;
  value: Uint8Array;
  headers: Header[];
  time: number;
};

/** A configuration object used to create new instances of `KafkaReader`. */
export type KafkaReaderConfig = {
  /** The list of broker addresses used to connect to the kafka cluster. */
  brokers: string[];
  /** The topic to read messages from. */
  topic: string;
  /**
   * Holds the optional consumer group id. If `groupId` is specified, then
   * `partition` should NOT be specified.
   */
  groupId?: string;
  /**
   * Partition to read messages from. Either `partition` or `groupId` may
   * be assigned, but not both.
   */
  partition?: number;
  sasl?: SASLConfig;
  tls?: TLSConfig;
  /** Setting this to true logs internal changes within the `KafkaReader`. */
  logger?: boolean;
};

export interface KafkaReader {
  /** Closes the stream, preventing the program from reading any more messages from it. */
  close: () => Promise<void>;
  /** Commits the list of messages passed as argument. */
  commitMessages: (msgs: KafkaReadMessage[]) => Promise<void>;
  /** Reads and return the next message. Does not commit offsets automatically
   *  when using consumer groups. Use `commitMessages` to commit the offset.
   */
  fetchMessage: () => Promise<KafkaReadMessage>;
  /**
   * Reads and return the next message. If consumer groups are used, `readMessage`
   * will automatically commit the offset when called. Note that this could result
   * in an offset being committed before the message is fully processed.
   *
   * If more fine grained control of when offsets are  committed is required, it
   * is recommended to use `fetchMessage` with `commitMessages` instead.
   */
  readMessage: () => Promise<KafkaReadMessage>;
  /**
   * SetOffset changes the offset from which the next batch of messages will be
   * read.
   */
  setOffset: (offset: number) => Promise<void>;
  /**
   * SetOffset changes the offset from which the next batch of messages will be
   * read given the timestamp.
   */
  setOffsetAt: (timeMs: number) => Promise<void>;
}
