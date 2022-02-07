import { SASLConfig } from "./security/sasl.ts";
import { TLSConfig } from "./security/tls.ts";

export type KafkaDialerConfig = {
  sasl?: SASLConfig;
  tls?: TLSConfig;
};

export interface KafkaConn {
  close: () => Promise<void>;
  apiVersions: () => Promise<string[]>;
}

export interface KafkaDialer {
  dial: (network: "tcp", address: string) => Promise<KafkaConn>;
}
