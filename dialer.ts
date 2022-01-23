import { SASLConfig } from "./sasl.ts";

export type KafkaDialerConfig = {
  sasl?: SASLConfig;
};

export interface KafkaConn {
  close: () => Promise<void>;
  apiVersions: () => Promise<string[]>;
}

export interface KafkaDialer {
  dial: (network: "tcp", address: string) => Promise<KafkaConn>;
}
