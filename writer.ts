import { Header } from "./header.ts";
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
  /** Maximum amount of time that connections will remain open and unused. */
  idleTimeout?: number;
  sasl?: SASLConfig;
  tls?: TLSConfig;
};

export interface KafkaWriter {
  /** Writes a batch of messages to the configured kafka topic. */
  writeMessages: (msgs: KafkaWriteMessage[]) => Promise<void>;
  /** Close flushes pending writes, and waits for all writes to complete before returning.  */
  close: () => Promise<void>;
}
