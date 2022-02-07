export enum SASLMechanism {
  /** Passes the credentials in clear text. */
  PLAIN = "PLAIN",
  SCRAMSHA512 = "SCRAM-SHA-512",
  SCRAMSHA256 = "SCRAM-SHA-256",
}

export type SASLConfig = {
  mechanism: SASLMechanism;
  username: string;
  password: string;
};
