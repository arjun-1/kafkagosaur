// import KafkaGoSaur from "../mod.ts";

// const kafkaGoSaur = new KafkaGoSaur();

// // const dialer = await kafkaGoSaur.dialer();
// // const kafkaConn = await dialer.dial("tcp", "localhost:29092");
// // const apiVersions = await kafkaConn.apiVersions();
// // console.log(apiVersions)

// const config = {
//   topic: "my-topic",
//   address: "localhost:29092",
// };

// const enc = new TextEncoder();

// const writer = await kafkaGoSaur.writer(config);
// const msgs = [{ value: enc.encode("hoihoi"), key: enc.encode("Key-A") }];
// await writer.writeMessages(msgs);

// console.log("finished?", new Date());
// JSON.
const conn = await Deno.connect({
  hostname: "localhost",
  port: 29092,
  transport: "tcp",
});
await conn.closeWrite();
await conn.close();
// const reader = (global as any).newReader({
//   brokers: ["localhost:29092"],
//   groupId: "my-group-id",
//   topic: "my-topic",
// });

// const message = await reader.readMessage()
// console.log(message)
