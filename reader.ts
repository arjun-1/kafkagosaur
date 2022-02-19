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
  /** The capacity of the internal message queue */
  queueCapacity?: number;
  /** Indicates to the broker the minimum batch size that the consumer will accept */
  minBytes?: number;
  /** Indicates to the broker the maximum batch size that the consumer will accept */
  maxBytes?: number;
  /** Maximum milliseconds to wait for new data to come when fetching batches of messages from kafka. */
  maxWait?: number;
  /** Frequency in milliseconds at which the reader lag is updated */
  readLagInterval?: number;
  /** Frequency in milliseconds at which the reader sends the consumer group heartbeat update. */
  heartBeatInterval?: number;
  /** Interval in milliseconds at which offsets are committed to the broker. */
  commitInterval?: number;
  /** How often in milliseconds a reader checks for partition changes */
  partitionWatchInterval?: number;
  /**
   * Used to inform kafka-go that a consumer group should bepolling the brokers
   * and rebalancing if any partition changes happen to the topic
   */
  watchPartitionChanges?: boolean;
  /**
   * Milliseconds that may pass without a heartbeat before the coordinator
   * considers the consumer dead and initiates a rebalance.
   */
  sessionTimeout?: number;
  /** Milliseconds the coordinator will wait for members to join as part of a rebalance */
  rebalanceTimeout?: number;
  /** Milliseconds to wait between re-joining the consumer group after an error */
  joinGroupBackoff?: number;
  /** Milliseconds the consumer group will be saved by the broker */
  retentionTime?: number;
  /** From whence the consumer group should begin consuming when it finds a partition without a committed offset */
  startOffset?: number;
  /** Smallest amount of milliseconds the reader will wait before polling for new messages */
  readBackoffMin?: number;
  /** Maximum amount of milliseconds the reader will wait before polling for new messages */
  readBackoffMax?: number;
  /** Limit of how many attempts will be made before delivering the error */
  maxAttempts?: number;
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
