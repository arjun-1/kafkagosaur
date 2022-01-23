# kafkagosaur
Kafkagosaur is a Deno wrapper for the [kafka-go](https://github.com/segmentio/kafka-go) project. 
Kafkagosaur is cross compiled from Go to a WebAssembly module, and runs natively in Deno.

## Supported features
- [x] Writer
- [x] Reader
- [x] SASL
- [ ] TLS
- [ ] Deno streams

## Examples
To run the examples, ensure you have docker up and running
```bash
make docker
```
To run the writer example
```bash
deno run --allow-read --allow-net examples/writer.ts
```

To run the reader example
```bash
deno run --allow-read --allow-net examples/reader.ts
```

## Development
To build the WebAssemnbly module, first run
```bash
make build
```
To run the tests, ensure first you have docker up and running
```bash
make docker
```
Then run 
```bash
make test
```
