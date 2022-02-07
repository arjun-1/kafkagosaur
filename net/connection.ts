/** Represents a network endpoint address, equivalent to a Go net.addr. */
export type Address = {
  /** Name of the network (for example, "tcp", "udp"). */
  network: string;
  /** String form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80"). */
  string: string;
};

/** A generic stream-oriented network connection, equivalent to a Go net.Conn. */
export interface Connection {
  /** Returns the local network address. */
  readonly localAddr: Address;
  /** Returns the remote network address. */
  readonly remoteAddr: Address;

  /**
   * Read reads data from the connection.
   * Read can be made to time out and return an error after a fixed
   * time limit; see SetDeadline and SetReadDeadline.
   */
  read(bytes: Uint8Array): Promise<number | null>;

  /**
   * Write writes data to the connection.
   * Write can be made to time out and return an error after a fixed
   * time limit; see SetDeadline and SetWriteDeadline.
   */
  write(bytes: Uint8Array): Promise<number>;

  close(): Promise<void>;

  setReadDeadline(timeMs: number): void;

  setWriteDeadline(timeMs: number): void;
}

export type Dial = (hostname: string, port: number) => Promise<Connection>;

export const setDialOnGlobal = (dial: Dial) => {
  (globalThis as Record<string, unknown>).dial = dial;
};

/** Specifies the implementation backing a TCP socket connection */
export enum DialBackend {
  Node,
  Deno,
}
