package kafkagosaur

import (
	"context"
	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
	"syscall/js"
)

type reader struct {
	underlying *kafka.Reader
	ctx        context.Context
}

func (r *reader) close() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := r.underlying.Close()

		if err != nil {
			reject(err)
		}

		resolve(nil)
	})
}

func (r *reader) commitMessages(msgs []kafka.Message) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := r.underlying.CommitMessages(r.ctx, msgs...)

		if err != nil {
			reject(err)
		}

		resolve(nil)
	})
}

// func (r *reader) config(msgs []kafka.Message) js.Value {}

func (r *reader) fetchMessage() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		msg, err := r.underlying.FetchMessage(r.ctx)

		if err != nil {
			reject(err)
		}

		resolve(msg)
	})
}

func (r *reader) readMessage() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		message, err := r.underlying.ReadMessage(r.ctx)

		if err != nil {
			reject(err)
		}

		resolve(messageToJSObject(message))
	})
}

func (r *reader) toJSObject() map[string]interface{} {
	return map[string]interface{}{
		"close": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return r.close()
			},
		),
		"commitMessages": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				// TODO: input validation

				msgsJs := args[0]
				msgs := make([]kafka.Message, msgsJs.Length())

				for i := range msgs {
					msgs[i] = jsObjectToMessage(msgsJs.Index(i))
				}

				return r.commitMessages(msgs)
			},
		),
		"fetchMessage": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return r.fetchMessage()
			},
		),
		"readMessage": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return r.readMessage()
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
