package main

import (
	"context"
	"fmt"

	"github.com/arjun-1/deno-wasm-experiment/jskafka"
	"syscall/js"

	"github.com/arjun-1/deno-wasm-experiment/interop"
	"github.com/segmentio/kafka-go"
)

func main() {
	c := make(chan bool)

	// reader := kafka.NewReader(kafka.ReaderConfig{
	// 	Topic: "my-topic",
	// 	Brokers: []string{"localhost:29092"},
	// 	GroupID: "my-group-id",
	// })

	topic := "my-topic"
	partition := 0

	// dialFunc := func(ctx context.Context, network string, address string) (net.Conn, error) {

	// 	return 1
	// }

	dialer := &kafka.Dialer{
		DialFunc: interop.NewDenoConn,
	}
	conn, err := dialer.DialLeader(context.Background(), "tcp", "localhost:29092", topic, partition)

	// conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:29092", topic, partition)
	if err != nil {
		fmt.Println("failed to dial leader:", err.Error())
	}

	result, _ := conn.ApiVersions()

	fmt.Println(result)

	js.Global().Set("newReader", jskafka.JsNewJsReader)

	<-c
}

// func Foo(this js.Value, p []js.Value) interface{} {
// 	return js.ValueOf(42)
// }

// func main() {
// 	// c := make(chan struct{}, 0)
// 	// fmt.Println("hello deno")
// 	// js.Global().Set("add", js.FuncOf(Foo))
// 	js.Global().Set("export1", "Hello!");
//     <- make(chan bool)

// 	// <-c

// }
