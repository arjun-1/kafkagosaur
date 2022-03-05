package main

import (
	"github.com/arjun-1/kafkagosaur/src/kafkagosaur"
	"syscall/js"
)

func main() {
	c := make(chan bool)

	const globalGoInteropKey = "kafkagosaur.go"

	interopObject := map[string]interface{}{
		"createDialer": kafkagosaur.NewDialerJsFunc,
		"createReader": kafkagosaur.NewReaderJsFunc,
		"createWriter": kafkagosaur.NewWriterJsFunc,
	}

	js.Global().Set(globalGoInteropKey, interopObject)

	<-c
}
