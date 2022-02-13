package main

import (
	"github.com/arjun-1/kafkagosaur/src/kafkagosaur"
	"syscall/js"
)

func main() {
	c := make(chan bool)

	js.Global().Set("createDialer", kafkagosaur.NewDialerJsFunc)
	js.Global().Set("createReader", kafkagosaur.NewReaderJsFunc)
	js.Global().Set("createWriter", kafkagosaur.NewWriterJsFunc)

	<-c
}
