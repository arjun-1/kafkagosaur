# kafkagosaur

[![deno module](https://img.shields.io/endpoint?url=https%3A%2F%2Fdeno-visualizer.danopia.net%2Fshields%2Flatest-version%2Fx%2Fkafkagosaur%2Fmod.ts)](https://deno.land/x/kafkagosaur)
[![deno doc](https://doc.deno.land/badge.svg)](https://doc.deno.land/https/deno.land/x/kafkagosaur/mod.ts)
[![CI](https://github.com/arjun-1/kafkagosaur/actions/workflows/CI.yaml/badge.svg)](https://github.com/arjun-1/kafkagosaur/actions/workflows/CI.yaml)
[![codecov](https://codecov.io/gh/arjun-1/kafkagosaur/branch/master/graph/badge.svg)](https://codecov.io/gh/arjun-1/kafkagosaur)

Kafkagosaur is a Kafka client for Deno built using WebAssembly. The project
binds to the [kafka-go](https://github.com/segmentio/kafka-go) library.

## Supported features

- [x] Writer
- [x] Reader
- [x] SASL
- [x] TLS
- [ ] Deno streams

## Examples

For comprehensive examples on how to use kafkagosaur, head over to the
[provided examples](#provided-examples).

#### KafkaWriter

To write a message, first create `KafkaWriter` instance using the `createWriter`
function on the `KafkaGoSaur` instance:

```typescript
import KafkaGoSaur from "https://deno.land/x/kafkagosaur/mod.ts";

const kafkaGoSaur = new KafkaGoSaur();
const writer = await kafkaGoSaur.createWriter({
  broker: "localhost:9092",
  topic: "test-0",
});

const enc = new TextEncoder();
const msgs = [{ value: enc.encode("value") }];

await writer.writeMessages(msgs);
```

#### KafkaReader

To read a message, first create `KafkaReader` instance using the `createReader`
function on the `KafkaGoSaur` instance:

```typescript
import KafkaGoSaur from "https://deno.land/x/kafkagosaur/mod.ts";

const kafkaGoSaur = new KafkaGoSaur();
const reader = await kafkaGoSaur.createReader({
  brokers: ["localhost:9092"],
  topic: "test-0",
});

const readMsg = await reader.readMessage();
```

### Provided examples

To run the [provided examples](examples), ensure you have docker up and running.
Then start the kafka broker using

```bash
make docker-up
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

The API documentation is hosted
[here](https://doc.deno.land/https/deno.land/x/kafkagosaur/mod.ts).

## Development

To build the WebAssembly module, first run

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

The Deno benchmarks are located in [bench](bench) and can be run via

```bash
deno run --allow-read --allow-net --allow-env --unstable bench/reader.ts
```

```bash
deno run --allow-read --allow-net --allow-env --unstable bench/writer.ts
```

### Results

<img width="423" alt="bench" src="https://user-images.githubusercontent.com/8102654/152859490-81464138-56ff-43d2-92db-7727458a561b.png">

|                   | kafkagosaur[^3] | kafka-go[^2] |
| ----------------- | --------------- | ------------ |
| writeMessages[^1] | 3977 ± 244      | 4678 ± 207   |
| readMessage       | 1963 ± 190      | 2784 ± 374   |

[^1]: Batching 10.000 messages.

[^2]: Using a single goroutine.

[^3]: Using the `DialBackend.Node`.

#### Environment

- 2,6 GHz 6-Core Intel Core i7
- [Confluent Cloud Basic cluster](https://docs.confluent.io/cloud/current/clusters/cluster-types.html#basic-clusters);
  6 partitions

## Contributing

Kafkgagosaur is in early stage of development. Nevertheless your contributions
are highly valued and welcomed! Feel free to ask for new features, report bugs,
or submit your code.
