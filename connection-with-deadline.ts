import { deadline } from "./deps.ts";

const withAbsoluteDeadline = <T>(
  p: Promise<T>,
  deadlineMs?: number,
): Promise<T> => {
  if (deadlineMs !== undefined) {
    const delayMs = deadlineMs - Date.now();
    return deadline(p, delayMs);
  }

  return p;
};

class ConnectionWithDeadline implements Deno.Conn {
  #conn: Deno.Conn;
  #readDeadlineMs?: number;
  #writeDeadlineMs?: number;

  readonly localAddr: Deno.Addr;
  readonly remoteAddr: Deno.Addr;
  readonly rid: number;

  constructor(conn: Deno.Conn) {
    this.#conn = conn;
    this.localAddr = conn.localAddr;
    this.remoteAddr = conn.remoteAddr;
    this.rid = conn.rid;
  }

  read(p: Uint8Array): Promise<number | null> {
    return withAbsoluteDeadline(this.#conn.read(p), this.#readDeadlineMs);
  }

  write(p: Uint8Array): Promise<number> {
    return withAbsoluteDeadline(this.#conn.write(p), this.#writeDeadlineMs);
  }

  closeWrite(): Promise<void> {
    return this.#conn.closeWrite();
  }

  close(): void {
    return this.#conn.close();
  }

  setReadDeadline(timeMs: number): void {
    this.#readDeadlineMs = timeMs;
  }

  setWriteDeadline(timeMs: number): void {
    this.#writeDeadlineMs = timeMs;
  }
}

const connect = async (options: Deno.ConnectOptions) => {
  const conn = await Deno.connect(options);
  return new ConnectionWithDeadline(conn);
};

export const setOnGlobal = () =>
  (globalThis as Record<string, unknown>).connectWithDeadline = connect;
