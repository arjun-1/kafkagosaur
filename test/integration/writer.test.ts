import { withWriter, withWriterSASL } from "./setup.ts";
import { KafkaWriter } from "../../writer.ts";

Deno.test(
  "Writer.writeMessages should write messages",
  () =>
    withWriter()(async (writer: KafkaWriter) => {
      const enc = new TextEncoder();
      const msgs = [{ value: enc.encode("value") }];

      await writer.writeMessages(msgs);
    }),
);

Deno.test(
  "Writer.writeMessages should write messages (SASL)",
  () =>
    withWriterSASL(async (writer: KafkaWriter) => {
      const enc = new TextEncoder();
      const msgs = [{ value: enc.encode("value") }];

      await writer.writeMessages(msgs);
    }),
);
