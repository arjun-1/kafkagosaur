package kafkagosaur

import (
	"context"
	"syscall/js"
	"time"

	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
)

type writer struct {
	underlying *kafka.Writer
	transport  *kafka.Transport
}

func (w *writer) writeMessages(msgs []kafka.Message) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := w.underlying.WriteMessages(context.Background(), msgs...)

		if err != nil {
			reject(err)
		}

		resolve(nil)
	})
}

func (w *writer) close() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := w.underlying.Close()
		if err != nil {
			reject(err)
		}
		w.transport.CloseIdleConnections()
		resolve(nil)
	})
}

func (w *writer) toJSObject() map[string]interface{} {
	return map[string]interface{}{
		"writeMessages": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {

				msgsJs := args[0]
				msgs := make([]kafka.Message, msgsJs.Length())

				for i := range msgs {
					msgs[i] = jsObjectToMessage(msgsJs.Index(i))
				}

				return w.writeMessages(msgs)
			},
		),
		"close": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return w.close()
			},
		),
	}
}

var NewWriterJsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	writerConfig := args[0]

	saslMechanism, err := SASLMechanism(writerConfig)
	if err != nil {
		panic(err)
	}

	tls, err := TLSConfig(writerConfig)
	if err != nil {
		panic(err)
	}

	transport := &kafka.Transport{
		Dial: interop.NewDenoConn,
		SASL: saslMechanism,
		TLS:  tls,
	}

	if jsIdleTimeout := writerConfig.Get("idleTimeout"); !jsIdleTimeout.IsUndefined() {
		transport.IdleTimeout = time.Duration(jsIdleTimeout.Int()) * time.Millisecond
	}

	kafkaWriter := kafka.Writer{
		Addr: kafka.TCP(writerConfig.Get("address").String()),
		// Logger:    log.Default(),
		Transport: transport,
	}

	if jsTopic := writerConfig.Get("topic"); !jsTopic.IsUndefined() {
		kafkaWriter.Topic = jsTopic.String()
	}

	return (&writer{
		underlying: &kafkaWriter,
		transport:  transport,
	}).toJSObject()

})
