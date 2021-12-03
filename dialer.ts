export interface KafkaConn {
  apiVersions: () => Promise<string[]>;
  close: () => Promise<void>;
}

export interface Dialer {
  dial: (network: "tcp", address: string) => Promise<KafkaConn>;
}
