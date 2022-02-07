import { Address, Connection, Dial } from "./connection.ts";
import { joinHostPort, withAbsoluteDeadline } from "./util.ts";

export class DenoTCPConnection implements Connection {
  #conn: Deno.Conn;
  #readDeadlineMs?: number;
  #writeDeadlineMs?: number;

  readonly localAddr: Address;
  readonly remoteAddr: Address;

  constructor(conn: Deno.Conn) {
    if (
      conn.localAddr.transport !== "tcp" || conn.remoteAddr.transport !== "tcp"
    ) {
      throw Error(
        `Not supported dial option(s): ${conn.localAddr.transport}, ${conn.remoteAddr.transport}`,
      );
    }

    this.#conn = conn;
    this.localAddr = {
      network: "tcp",
      string: joinHostPort(conn.localAddr.hostname, conn.localAddr.port),
    };
    this.remoteAddr = {
      network: "tcp",
      string: joinHostPort(conn.remoteAddr.hostname, conn.remoteAddr.port),
    };
  }

  read(bytes: Uint8Array): Promise<number | null> {
    return withAbsoluteDeadline(this.#conn.read(bytes), this.#readDeadlineMs);
  }

  write(bytes: Uint8Array): Promise<number> {
    return withAbsoluteDeadline(this.#conn.write(bytes), this.#writeDeadlineMs);
  }

  close(): Promise<void> {
    this.#conn.close();
    return this.#conn.closeWrite();
  }

  setReadDeadline(timeMs: number): void {
    this.#readDeadlineMs = timeMs;
  }

  setWriteDeadline(timeMs: number): void {
    this.#writeDeadlineMs = timeMs;
  }
}

export const dial: Dial = async (hostname: string, port: number) => {
  const conn = await Deno.connect({ hostname, port });
  return new DenoTCPConnection(conn);
};
