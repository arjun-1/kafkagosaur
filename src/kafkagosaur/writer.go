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
}

func jsObjectToMessage(jsObject js.Value) kafka.Message {
	message := kafka.Message{}

	jsKey := jsObject.Get("key")
	if !jsKey.IsUndefined() {
		key := make([]byte, jsKey.Length())
		js.CopyBytesToGo(key, jsKey)
		message.Key = key
	}

	jsValue := jsObject.Get("value")
	if !jsValue.IsUndefined() {
		value := make([]byte, jsValue.Length())
		js.CopyBytesToGo(value, jsValue)
		message.Value = value
	}

	jsTime := jsObject.Get("time")
	if !jsTime.IsUndefined() {
		time := time.UnixMilli(int64(jsTime.Int()))
		message.Time = time
	}

	jsTopic := jsObject.Get("topic")
	if !jsTopic.IsUndefined() {
		message.Topic = jsTopic.String()
	}

	return message
}

func (w *writer) writeMessagesPromise(msgs []kafka.Message) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := w.underlying.WriteMessages(context.Background(), msgs...)

		if err != nil {
			reject(err)
		}

		resolve(nil)
	})
}

func (w *writer) closePromise() js.Value {
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

				return w.writeMessagesPromise(msgs)
			},
		),
		"close": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return w.closePromise()
			},
		),
	}
}

func newWriter(kafkaWriter kafka.Writer) *writer {

	return &writer{
		underlying: &kafkaWriter,
	}
}

var NewWriterJsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// TODO: input validation
	writerConfig := args[0]

	transport := &kafka.Transport{
		Dial: interop.NewDenoConn,
	}

	kafkaWriter := kafka.Writer{
		Addr:      kafka.TCP(writerConfig.Get("address").String()),
		Logger:    log.Default(),
		Transport: transport,
	}

	jsTopic := writerConfig.Get("topic")
	if !jsTopic.IsUndefined() {
		kafkaWriter.Topic = jsTopic.String()
	}
	// TODO recover GET panic

	writer := newWriter(kafkaWriter)

	return writer.toJSObject()

})
