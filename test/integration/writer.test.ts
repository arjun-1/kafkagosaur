import { withWriter } from "./setup.ts";
import { Writer } from "../../writer.ts";

Deno.test(
  "Writer.writeMessages should write messages",
  () =>
    withWriter(async (writer: Writer) => {
      const enc = new TextEncoder();
      const msgs = [{ value: enc.encode("value") }];

      await writer.writeMessages(msgs);
    }),
);
