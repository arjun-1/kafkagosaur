package jskafka

import (
	"context"
	"syscall/js"

	"github.com/arjun-1/deno-wasm-experiment/interop"
	"github.com/segmentio/kafka-go"
)

type jsReader struct {
	underlying *kafka.Reader
}

func (r *jsReader) readMessagePromise() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		message, err := r.underlying.ReadMessage(context.Background())

		if err != nil {
			reject(err)
		}

		// TODO map to object
		resolve(message.Topic)
	})
}

func (r *jsReader) closePromise() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve(r.underlying.Close())
	})
}

func (r *jsReader) toJSObject() map[string]interface{} {
	return map[string]interface{}{
		"readMessage": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return r.readMessagePromise()
			},
		),
		"close": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return r.closePromise()
			},
		),
	}
}

func newJsReader(kafkaReaderConfig kafka.ReaderConfig) *jsReader {
	dialer := kafkaReaderConfig.Dialer
	dialer.DialFunc = interop.NewDenoConn
	kafkaReaderConfig.Dialer = dialer

	reader := kafka.NewReader(kafkaReaderConfig)

	return &jsReader{
		underlying: reader,
	}
}

var JsNewJsReader = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// TODO: input validation
	jsReaderConfig := args[0]

	readerConfig := kafka.ReaderConfig{
		Brokers: interop.MapToString(interop.ToSlice(jsReaderConfig.Get("brokers"))),
		GroupID: jsReaderConfig.Get("groupId").String(),
		Topic:   jsReaderConfig.Get("topic").String(),
	}

	// TODO recover GET panic

	reader := newJsReader(readerConfig)

	return reader.toJSObject()

})
