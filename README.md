# kafkagosaur

[![deno module](https://shield.deno.dev/x/kafkagosaur)](https://deno.land/x/kafkagosaur)
[![CI](https://github.com/arjun-1/kafkagosaur/actions/workflows/CI.yaml/badge.svg)](https://github.com/arjun-1/kafkagosaur/actions/workflows/CI.yaml)
[![codecov](https://codecov.io/gh/arjun-1/kafkagosaur/branch/master/graph/badge.svg)](https://codecov.io/gh/arjun-1/kafkagosaur)

Kafkagosaur is a Kafka client for Deno built using WebAssembly. The project
binds to the [kafka-go](https://github.com/segmentio/kafka-go) library.

Being cross-compiled from Go into a WebAssembly module, it should enjoy
performance characteristics similar to the native code.

## Supported features

- [x] Writer
- [x] Reader
- [x] SASL
- [x] TLS
- [x] TCP
- [ ] UDP
- [ ] Deno streams

## Examples

For comprehensive examples on how to use kafkagosaur, head over to the
[provided examples](#provided-examples).

#### KafkaWriter

To write a message, make use of the `KafkaWriter` object on the `KafkaGoSaur`
instance:

```typescript
const kafkaGoSaur = new KafkaGoSaur();
const writer = await kafkaGoSaur.writer({
  broker: "localhost:9092",
  topic: "test-0",
});

const enc = new TextEncoder();
const msgs = [{ value: enc.encode("value") }];

await writer.writeMessages(msgs);
```

#### KafkaReader

To read a message, make use of the `KafkaReader` object on the `KafkaGoSaur`
instance:

```typescript
const kafkaGoSaur = new KafkaGoSaur();
const reader = await kafkaGoSaur.reader({
  brokers: ["localhost:9092"],
  topic: "test-0",
});

const readMsg = await reader.readMessage();
```

### Provided examples

To run the [provided examples](examples), ensure you have docker up and running.
Then start the kafka broker using

```bash
make docker
```

To run the writer example

```bash
deno run --allow-read --allow-net --unstable examples/writer.ts
```

To run the reader example

```bash
deno run --allow-read --allow-net --unstable examples/reader.ts
```

## Documentation

While the documentation is a work in progress, please see the
[documentation](https://github.com/segmentio/kafka-go/blob/main/README.md) of
the kafka-go project on how to use kafkagosaur. Kafkagosaur tries to follow the
API of kafka-go closely.

## Development

To build the WebAssemnbly module, first run

```bash
make build-wasm
```

To run the tests, ensure first you have docker up and running. Then start the
kafka broker using

```bash
make docker-up
```

Then run

```bash
make test
```

## Performance benchmarks

TODO

## Contributing

Kafkgagosaur is in early stage of development. Nevertheless your contributions
are highly valued and welcomed! Feel free to ask for new features, report bugs,
or submit your code.
