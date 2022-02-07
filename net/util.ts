import { deadline } from "../deps.ts";

export const withAbsoluteDeadline = <T>(
  p: Promise<T>,
  deadlineMs?: number,
): Promise<T> => {
  if (deadlineMs !== undefined) {
    const delayMs = deadlineMs - Date.now();
    return deadline(p, delayMs);
  }

  return p;
};

export const joinHostPort = (host: string, port: number) =>
  host.indexOf(":") >= 0 ? `[${host}]:${port}` : `${host}:${port}`;
