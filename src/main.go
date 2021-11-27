package main

import (
	"github.com/arjun-1/deno-wasm-experiment/jskafka"
	"syscall/js"
)

func main() {
	c := make(chan bool)

	js.Global().Set("newReader", jskafka.JsNewJsReader)
	js.Global().Set("newDialer", jskafka.JsNewJsDialer)

	<-c
}
