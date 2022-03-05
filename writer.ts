import { Header } from "./header.ts";
import { DialBackend } from "./net/connection.ts";
import { SASLConfig } from "./security/sasl.ts";
import { TLSConfig } from "./security/tls.ts";

/** A data structure representing kafka messages written to `KafkaWriter`. */
export type KafkaWriteMessage = {
  /**
   * Can be used to configured the topic if not already specified on
   * `KafkWriter` itself.
   */
  topic?: string;
  offset?: number;
  highWaterMark?: number;
  key?: Uint8Array;
  value?: Uint8Array;
  headers?: Header[];
  /** If not set, will be automatically set when writing the message. */
  time?: number;
};

/** A configuration object used to create new instances of `KafkaWriter`. */
export type KafkaWriterConfig = {
  /**
   * Name of the topic that `KafkaWriter` will produce messages to.
   * Setting this field or not is a mutually exclusive option. If you set `topic`
   * here, you must not set `topic` for any `KafkaWriteMessage`. Otherwise, if you	do
   * not set `topic`, every `KafkaWriteMessage` must have `topic` specified.
   */
  topic?: string;
  /** Address of the kafka cluster that this writer is configured to send messages to. */
  address: string;
  sasl?: SASLConfig;
  tls?: TLSConfig;
  /** Time limit in milliseconds set for establishing connections to the kafka cluster. */
  dialTimeout?: number;
  /** Maximum amount of time that connections will remain open and unused. */
  idleTimeout?: number;
  /** TTL in milliseconds for the metadata cached by this transport. */
  metadataTTL?: number;
  /** Unique identifier that the transport communicates to the brokers when it sends requests. */
  clientId?: string;
  /** Limit on how many attempts will be made to deliver a message. */
  maxAttempts?: number;
  /** Limit on how many messages will be buffered before being sent to a partition */
  batchSize?: number;
  /** Limit the maximum size of a request in bytes before being sent to a partition. */
  batchBytes?: number;
  /** Time limit in milliseconds on how often incomplete message batches will be flushed to kafka. */
  batchTimeout?: number;
  /** Timeout in milliseconds for read operations performed by the Writer. */
  readTimeout?: number;
  /** Timeout in milliseconds for write operations performed by the Writer. */
  writeTimeout?: number;
  /** Setting this flag to true causes the WriteMessages method to never block. */
  async?: boolean;
  /** Setting this to true logs internal changes within the `KafkaReader`. */
  logger?: boolean;
  /** Specifies the implementation backing a TCP socket connection. Defaults to Node */
  dialBackend?: DialBackend;
};

export interface KafkaWriter {
  /** Writes a batch of messages to the configured kafka topic. */
  writeMessages: (msgs: KafkaWriteMessage[]) => Promise<void>;
  /** Close flushes pending writes, and waits for all writes to complete before returning.  */
  close: () => Promise<void>;
  /**
   * Stats returns a snapshot of the writer stats since the last time the method
   * was called, or since the writer was created if it is called for the first time.
   */
  stats: () => string;
}
