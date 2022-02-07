export type X509KeyPair = {
  key: string;
  cert: string;
};

export type TLSConfig = boolean | {
  insecureSkipVerify?: boolean;
  keyPair?: X509KeyPair;
  ca?: string;
};
