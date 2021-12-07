package kafkagosaur

import (
	"context"
	"github.com/arjun-1/kafkagosaur/interop"
	"github.com/segmentio/kafka-go"
	"syscall/js"
)

type reader struct {
	underlying *kafka.Reader
	ctx        context.Context
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
		message, err := r.underlying.ReadMessage(r.ctx)

		if err != nil {
			reject(err)
		}

		resolve(messageToJSObject(message))
	})
}

func (r *reader) closePromise() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := r.underlying.Close()

		if err != nil {
			reject(err)
		}

		resolve(nil)
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

var NewReaderJsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// TODO: input validation
	readerConfig := args[0]

	// TODO recover GET panic

	kafkaReaderConfig := kafka.ReaderConfig{
		Brokers: interop.MapToString(interop.ToSlice(readerConfig.Get("brokers"))),
		GroupID: readerConfig.Get("groupId").String(),
		Topic:   readerConfig.Get("topic").String(),
	}

	kafkaReaderConfig.Dialer = &kafka.Dialer{
		DialFunc: interop.NewDenoConn,
	}

	kafkaReader := kafka.NewReader(kafkaReaderConfig)

	return (&reader{
		underlying: kafkaReader,
		ctx:        context.Background(),
	}).toJSObject()

})
