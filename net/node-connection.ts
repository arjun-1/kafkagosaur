import { Address, Connection, Dial } from "./connection.ts";
import { joinHostPort, withAbsoluteDeadline } from "./util.ts";
import { createConnection, deferred, Socket } from "../deps.ts";

export class NodeTCPConnection implements Connection {
  #socket: Socket;
  #readDeadlineMs?: number;
  #writeDeadlineMs?: number;

  readonly localAddr: Address;
  readonly remoteAddr: Address;

  constructor(socket: Socket) {
    if (socket.remoteAddress === undefined || socket.remotePort === undefined) {
      throw Error("Socket remoteAddress or remotePort are undefined.");
    }
    this.#socket = socket;

    this.localAddr = {
      network: "tcp",
      string: joinHostPort(socket.localAddress, socket.localPort),
    };
    this.remoteAddr = {
      network: "tcp",
      string: joinHostPort(socket.remoteAddress, socket.remotePort),
    };
  }

  read(bytes: Uint8Array): Promise<number | null> {
    const p = deferred<number | null>();

    const onReadable = () => {
      const readSize = Math.min(bytes.length, this.#socket.readableLength);
      const chunk = this.#socket.read(readSize);

      if (chunk === null) {
        // wait for next 'readable' event
        return;
      }

      if (chunk !== undefined && typeof chunk !== "string") {
        bytes.set(chunk);
        removeListeners();

        return p.resolve(chunk.length);
      }

      removeListeners();

      p.reject(Error("Invalid chunk read"));
    };
    const onError = (err: Error) => {
      removeListeners();
      p.reject(err);
    };
    const onEnd = () => {
      removeListeners();
      p.resolve(null); // EOF
    };
    const onClose = () => {
      removeListeners();
      p.resolve(0);
    };

    this.#socket.on("readable", onReadable);
    this.#socket.once("error", onError);
    this.#socket.once("end", onEnd);
    this.#socket.once("close", onClose);

    const removeListeners = () => {
      this.#socket.removeListener("readable", onReadable);
      this.#socket.removeListener("error", onError);
      this.#socket.removeListener("end", onEnd);
      this.#socket.removeListener("close", onClose);
    };

    onReadable();

    return withAbsoluteDeadline(p, this.#readDeadlineMs);
  }

  write(bytes: Uint8Array): Promise<number> {
    const p = deferred<number>();

    this.#socket.write(bytes, undefined, (error: Error | null | undefined) => {
      if (error instanceof Error) {
        p.reject(error);
      } else {
        p.resolve(bytes.length);
      }
    });

    return withAbsoluteDeadline(p, this.#writeDeadlineMs);
  }

  close(): Promise<void> {
    const p = deferred<void>();
    const onError = (err: Error) => {
      removeListeners();
      p.reject(err);
    };
    const onClose = () => {
      removeListeners();
      p.resolve();
    };

    this.#socket.once("error", onError);
    this.#socket.once("close", onClose);

    const removeListeners = () => {
      this.#socket.removeListener("error", onError);
      this.#socket.removeListener("close", onClose);
    };

    this.#socket.destroy();

    return p;
  }

  setReadDeadline(timeMs: number): void {
    this.#readDeadlineMs = timeMs;
  }

  setWriteDeadline(timeMs: number): void {
    this.#writeDeadlineMs = timeMs;
  }
}

export const dial: Dial = async (host: string, port: number) => {
  const p = deferred<void>();

  const socket = createConnection({ port, host }, p.resolve);

  await p;

  return new NodeTCPConnection(socket);
};
