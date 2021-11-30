package kafkagosaur

import (
	"context"
	"github.com/arjun-1/kafkagosaur/interop"
	"github.com/segmentio/kafka-go"
	"syscall/js"
)

type reader struct {
	underlying *kafka.Reader
}

func messageToJSObject(m kafka.Message) map[string]interface{} {
	key := interop.NewUint8Array(len(m.Key))
	js.CopyBytesToJS(key, m.Key)

	value := interop.NewUint8Array(len(m.Value))
	js.CopyBytesToJS(value, m.Value)

	return map[string]interface{}{
		"topic":     m.Topic,
		"partition": m.Partition,
		"time":      m.Time.UnixMilli(),
		"key":       key,
		"value":     value,
	}
}

func (r *reader) readMessagePromise() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		message, err := r.underlying.ReadMessage(context.Background())

		if err != nil {
			reject(err)
		}

		resolve(messageToJSObject(message))
	})
}

func (r *reader) closePromise() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		resolve(r.underlying.Close())
	})
}

func (r *reader) toJSObject() map[string]interface{} {
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

func newReader(kafkaReaderConfig kafka.ReaderConfig) *reader {
	kafkaDialer := kafkaReaderConfig.Dialer
	kafkaDialer.DialFunc = interop.NewDenoConn
	kafkaReaderConfig.Dialer = kafkaDialer

	kafkaReader := kafka.NewReader(kafkaReaderConfig)

	return &reader{
		underlying: kafkaReader,
	}
}

var NewReaderJsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// TODO: input validation
	readerConfig := args[0]

	kafkaReaderConfig := kafka.ReaderConfig{
		Brokers: interop.MapToString(interop.ToSlice(readerConfig.Get("brokers"))),
		GroupID: readerConfig.Get("groupId").String(),
		Topic:   readerConfig.Get("topic").String(),
	}

	// TODO recover GET panic

	reader := newReader(kafkaReaderConfig)

	return reader.toJSObject()

})
