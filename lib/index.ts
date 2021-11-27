// @deno-types="../global.d.ts"
import "./wasm_exec.js";

const go = new global.Go();

const wasmCode = await Deno.readFile("./bin/wasm.wasm");
const inst = await WebAssembly.instantiate(wasmCode, go.importObject);

go.run(inst.instance);

const newDialer = (global as any).newDialer;

const dialer = newDialer();
const kafkaConn = await dialer.dial("tcp", "localhost:29092");
const apiVersions = await kafkaConn.apiVersions();
console.log(apiVersions);

// const conn = await Deno.connect({})
// conn.close
// const reader = (global as any).newReader({
//   brokers: ["localhost:29092"],
//   groupId: "my-group-id",
//   topic: "my-topic",
// });

// const message = await reader.readMessage()
// console.log(message)
