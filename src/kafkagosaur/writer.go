package kafkagosaur

import (
	"context"
	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
	"log"
	"syscall/js"
	"time"
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
				// TODO: input validation

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
	// TODO: input validation
	writerConfig := args[0]

	transport := &kafka.Transport{
		Dial: interop.NewDenoConn,
	}

	// TODO recover GET panic

	if jsIdleTimeout := writerConfig.Get("idleTimeout"); !jsIdleTimeout.IsUndefined() {
		transport.IdleTimeout = time.Duration(jsIdleTimeout.Int()) * time.Millisecond
	}

	kafkaWriter := kafka.Writer{
		Addr:      kafka.TCP(writerConfig.Get("address").String()),
		Logger:    log.Default(),
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
