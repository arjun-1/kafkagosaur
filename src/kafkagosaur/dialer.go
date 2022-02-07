package kafkagosaur

import (
	"syscall/js"

	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
)

type dialer struct {
	underlying *kafka.Dialer
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

func NewKafkaDialer(dialConfigJs js.Value) *kafka.Dialer {

	saslMechanism, err := SASLMechanism(dialConfigJs)
	if err != nil {
		panic(err)
	}

	tls, err := TLSConfig(dialConfigJs)
	if err != nil {
		panic(err)
	}

	kafkaDialer := &kafka.Dialer{
		DialFunc:      interop.NewDenoConn,
		SASLMechanism: saslMechanism,
		TLS:           tls,
	}

	return kafkaDialer
}

var NewDialerJsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	dialConfigJs := args[0]

	kafkaDialer := NewKafkaDialer(dialConfigJs)

	return (&dialer{
		underlying: kafkaDialer,
	}).toJSObject()

})
