package kafkagosaur

import (
	"context"
	"syscall/js"
	"time"

	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
)

type reader struct {
	underlying *kafka.Reader
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
		err := r.underlying.CommitMessages(context.Background(), msgs...)

		if err != nil {
			reject(err)
		}

		resolve(nil)
	})
}

func (r *reader) fetchMessage() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		msg, err := r.underlying.FetchMessage(context.Background())

		if err != nil {
			reject(err)
		}

		resolve(msg)
	})
}

func (r *reader) readMessage() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		message, err := r.underlying.ReadMessage(context.Background())

		if err != nil {
			reject(err)
		}

		resolve(messageToJSObject(message))
	})
}

func (r *reader) setOffset(offset int64) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := r.underlying.SetOffset(offset)

		if err != nil {
			reject(err)
		}

		resolve(nil)
	})
}

func (r *reader) setOffsetAt(t time.Time) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := r.underlying.SetOffsetAt(context.Background(), t)

		if err != nil {
			reject(err)
		}

		resolve(nil)
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
		"setOffset": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				offset := int64(args[0].Int())

				return r.setOffset(offset)
			},
		),
		"setOffsetAt": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				time := time.UnixMilli(int64(args[0].Int()))

				return r.setOffsetAt(time)
			},
		),
	}
}

var NewReaderJsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

	readerConfigJs := args[0]

	kafkaDialer := NewKafkaDialer(readerConfigJs)

	kafkaReaderConfig := kafka.ReaderConfig{
		Dialer: kafkaDialer,
	}

	if brokers := readerConfigJs.Get("brokers"); !brokers.IsUndefined() {
		kafkaReaderConfig.Brokers = interop.MapToString(interop.ToSlice(brokers))
	}

	if groupId := readerConfigJs.Get("groupId"); !groupId.IsUndefined() {
		kafkaReaderConfig.GroupID = groupId.String()
	}

	if topic := readerConfigJs.Get("topic"); !topic.IsUndefined() {
		kafkaReaderConfig.Topic = topic.String()
	}

	if partition := readerConfigJs.Get("partition"); !partition.IsUndefined() {
		kafkaReaderConfig.Partition = partition.Int()
	}

	if minBytes := readerConfigJs.Get("minBytes"); !minBytes.IsUndefined() {
		kafkaReaderConfig.MinBytes = minBytes.Int()
	}

	if maxBytes := readerConfigJs.Get("maxBytes"); !maxBytes.IsUndefined() {
		kafkaReaderConfig.MaxBytes = maxBytes.Int()
	}

	kafkaReader := kafka.NewReader(kafkaReaderConfig)

	return (&reader{
		underlying: kafkaReader,
	}).toJSObject()

})
