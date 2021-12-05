package kafkagosaur

import (
	"context"
	"log"
	"syscall/js"
	"time"

	"github.com/arjun-1/kafkagosaur/interop"
	"github.com/segmentio/kafka-go"
)

type writer struct {
	underlying *kafka.Writer
	ctx        context.Context
}

func jsObjectToMessage(jsObject js.Value) kafka.Message {
	message := kafka.Message{}

	if jsKey := jsObject.Get("key"); !jsKey.IsUndefined() {
		key := make([]byte, jsKey.Length())
		js.CopyBytesToGo(key, jsKey)
		message.Key = key
	}

	if jsValue := jsObject.Get("value"); !jsValue.IsUndefined() {
		value := make([]byte, jsValue.Length())
		js.CopyBytesToGo(value, jsValue)
		message.Value = value
	}

	if jsTime := jsObject.Get("time"); !jsTime.IsUndefined() {
		time := time.UnixMilli(int64(jsTime.Int()))
		message.Time = time
	}

	if jsTopic := jsObject.Get("topic"); !jsTopic.IsUndefined() {
		message.Topic = jsTopic.String()
	}

	return message
}

func (w *writer) writeMessagesPromise(ctx context.Context, msgs []kafka.Message) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := w.underlying.WriteMessages(context.Background(), msgs...)

		if err != nil {
			reject(err)
		}

		resolve(nil)
	})
}

func (w *writer) closePromise(ctx context.Context) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve(w.underlying.Close())
	})
}

func (w *writer) toJSObject() map[string]interface{} {
	return map[string]interface{}{
		"writeMessages": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				// TODO: input validation

				jsMsgs := args[0]
				msgs := make([]kafka.Message, jsMsgs.Length())

				for i := range msgs {
					msgs[i] = jsObjectToMessage(jsMsgs.Index(i))
				}

				return w.writeMessagesPromise(w.ctx, msgs)
			},
		),
		"close": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return w.closePromise(w.ctx)
			},
		),
	}
}

func newWriter(kafkaWriter kafka.Writer) *writer {

	return &writer{
		underlying: &kafkaWriter,
		ctx:        context.Background(),
	}
}

var NewWriterJsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// TODO: input validation
	writerConfig := args[0]

	transport := &kafka.Transport{
		Dial: interop.NewDenoConn,
	}

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
	// TODO recover GET panic

	writer := newWriter(kafkaWriter)

	return writer.toJSObject()

})
