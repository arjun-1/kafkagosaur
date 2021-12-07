package kafkagosaur

import (
	"context"
	"syscall/js"

	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
)

type dialer struct {
	underlying *kafka.Dialer
	ctx        context.Context
}

func (d *dialer) dial(network string, address string) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		kafkaConn, err := d.underlying.Dial(network, address)

		if err != nil {
			reject(err)
		}

		jsConn := conn{
			kafkaConn,
		}

		resolve(jsConn.toJSObject())
	})
}

func (d *dialer) toJSObject() map[string]interface{} {
	return map[string]interface{}{
		"dial": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				// TODO: input validation
				network := args[0].String()
				address := args[1].String()

				return d.dial(network, address)
			},
		),
	}
}

var NewDialerJsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	kafkaDialer := &kafka.Dialer{
		DialFunc: interop.NewDenoConn,
	}

	return (&dialer{
		underlying: kafkaDialer,
		ctx:        context.Background(),
	}).toJSObject()

})
