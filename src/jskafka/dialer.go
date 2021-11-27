package jskafka

import (
	"github.com/arjun-1/deno-wasm-experiment/interop"
	"github.com/segmentio/kafka-go"
	"syscall/js"
)

type jsDialer struct {
	underlying *kafka.Dialer
}

func (d *jsDialer) dialPromise(network string, address string) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		conn, err := d.underlying.Dial(network, address)

		if err != nil {
			reject(err)
		}

		jsConn := jsConn{
			conn,
		}

		resolve(jsConn.toJSObject())
	})
}

func (d *jsDialer) toJSObject() map[string]interface{} {
	return map[string]interface{}{
		"dial": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				// TODO: input validation
				network := args[0].String()
				address := args[1].String()

				return d.dialPromise(network, address)
			},
		),
	}
}

func newJsDialer() *jsDialer {
	dialer := &kafka.Dialer{
		DialFunc: interop.NewDenoConn,
	}

	return &jsDialer{
		underlying: dialer,
	}
}

var JsNewJsDialer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	dialer := newJsDialer()

	return dialer.toJSObject()
})
