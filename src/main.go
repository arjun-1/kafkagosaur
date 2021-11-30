package main

import (
	"github.com/arjun-1/kafkagosaur/kafkagosaur"
	"syscall/js"
)

func main() {
	c := make(chan bool)

	js.Global().Set("newDialer", kafkagosaur.NewDialerJsFunc)
	js.Global().Set("newReader", kafkagosaur.NewReaderJsFunc)
	js.Global().Set("newWriter", kafkagosaur.NewWriterJsFunc)

	<-c
}
