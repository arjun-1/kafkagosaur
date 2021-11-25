package main

import (
	"context"
	"fmt"
	"github.com/arjun-1/deno-wasm-experiment/interop"
	"github.com/arjun-1/deno-wasm-experiment/jskafka"
	"github.com/segmentio/kafka-go"
	"syscall/js"
)

func main() {
	c := make(chan bool)

	topic := "my-topic"
	partition := 0

	dialer := &kafka.Dialer{
		DialFunc: interop.NewDenoConn,
	}
	conn, err := dialer.DialLeader(context.Background(), "tcp", "localhost:29092", topic, partition)

	if err != nil {
		fmt.Println("failed to dial leader:", err.Error())
	}

	result, _ := conn.ApiVersions()

	fmt.Println(result)

	js.Global().Set("newReader", jskafka.JsNewJsReader)

	<-c
}
