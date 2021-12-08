export interface KafkaConn {
  close: () => Promise<void>;
  apiVersions: () => Promise<string[]>;
}

export interface Dialer {
  dial: (network: "tcp", address: string) => Promise<KafkaConn>;
}
