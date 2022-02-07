/**
 * X509KeyPair holds a public/private key of PEM encoded data.
 */
export type X509KeyPair = {
  key: string;
  cert: string;
};

export type TLSConfig = boolean | {
  /**
   * Controls whether a client verifies the server's certificate chain and
   * host name. If set true, crypto/tls accepts any certificate presented
   * by the server and any host name in that certificate. In this mode, TLS
   * is susceptible to machine-in-the-middle attacks unless custom verification
   * is used. This should be used only for testing.
   */
  insecureSkipVerify?: boolean;
  keyPair?: X509KeyPair;
  /**
   * Defines the set of root certificate authority that clients use when
   * verifying server certificates. If not set, use the host's root CA set.
   */
  ca?: string;
};
