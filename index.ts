// @deno-types="./global.d.ts"
import "./wasm_exec.js";

const go = new global.Go();


const wasmCode = await Deno.readFile("./wasm.wasm");
const inst = await WebAssembly.instantiate(wasmCode, go.importObject);

go.run(inst.instance);

// const reader = (global as any).newReader({
//   brokers: ["localhost:29092"],
//   groupId: "my-group-id",
//   topic: "my-topic",
// });


// const message = await reader.readMessage()
// console.log(message)

