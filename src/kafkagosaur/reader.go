package kafkagosaur

import (
	"context"
	"log"
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
		message, err := r.underlying.FetchMessage(context.Background())

		if err != nil {
			reject(err)
		}

		resolve(messageToJSObject(message))
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

	if partition := readerConfigJs.Get("partition"); !partition.IsUndefined() {
		kafkaReaderConfig.Partition = partition.Int()
	}

	if topic := readerConfigJs.Get("topic"); !topic.IsUndefined() {
		kafkaReaderConfig.Topic = topic.String()
	}

	if partition := readerConfigJs.Get("partition"); !partition.IsUndefined() {
		kafkaReaderConfig.Partition = partition.Int()
	}

	if queueCapacity := readerConfigJs.Get("queueCapacity"); !queueCapacity.IsUndefined() {
		kafkaReaderConfig.QueueCapacity = queueCapacity.Int()
	}

	if minBytes := readerConfigJs.Get("minBytes"); !minBytes.IsUndefined() {
		kafkaReaderConfig.MinBytes = minBytes.Int()
	}

	if maxBytes := readerConfigJs.Get("maxBytes"); !maxBytes.IsUndefined() {
		kafkaReaderConfig.MaxBytes = maxBytes.Int()
	}

	if maxWait := readerConfigJs.Get("maxWait"); !maxWait.IsUndefined() {
		kafkaReaderConfig.MaxWait = time.Millisecond * time.Duration(maxWait.Int())
	}

	if readLagInterval := readerConfigJs.Get("readLagInterval"); !readLagInterval.IsUndefined() {
		kafkaReaderConfig.ReadLagInterval = time.Millisecond * time.Duration(readLagInterval.Int())
	}

	if heartbeatInterval := readerConfigJs.Get("heartbeatInterval"); !heartbeatInterval.IsUndefined() {
		kafkaReaderConfig.HeartbeatInterval = time.Millisecond * time.Duration(heartbeatInterval.Int())
	}

	if commitInterval := readerConfigJs.Get("commitInterval"); !commitInterval.IsUndefined() {
		kafkaReaderConfig.CommitInterval = time.Millisecond * time.Duration(commitInterval.Int())
	}

	if partitionWatchInterval := readerConfigJs.Get("partitionWatchInterval"); !partitionWatchInterval.IsUndefined() {
		kafkaReaderConfig.PartitionWatchInterval = time.Millisecond * time.Duration(partitionWatchInterval.Int())
	}

	if watchPartitionChanges := readerConfigJs.Get("watchPartitionChanges"); !watchPartitionChanges.IsUndefined() {
		kafkaReaderConfig.WatchPartitionChanges = watchPartitionChanges.Bool()
	}

	if sessionTimeout := readerConfigJs.Get("sessionTimeout"); !sessionTimeout.IsUndefined() {
		kafkaReaderConfig.SessionTimeout = time.Millisecond * time.Duration(sessionTimeout.Int())
	}

	if rebalanceTimeout := readerConfigJs.Get("rebalanceTimeout"); !rebalanceTimeout.IsUndefined() {
		kafkaReaderConfig.RebalanceTimeout = time.Millisecond * time.Duration(rebalanceTimeout.Int())
	}

	if joinGroupBackoff := readerConfigJs.Get("joinGroupBackoff"); !joinGroupBackoff.IsUndefined() {
		kafkaReaderConfig.JoinGroupBackoff = time.Millisecond * time.Duration(joinGroupBackoff.Int())
	}

	if retentionTime := readerConfigJs.Get("retentionTime"); !retentionTime.IsUndefined() {
		kafkaReaderConfig.RetentionTime = time.Millisecond * time.Duration(retentionTime.Int())
	}

	if startOffset := readerConfigJs.Get("startOffset"); !startOffset.IsUndefined() {
		kafkaReaderConfig.StartOffset = int64(startOffset.Int())
	}

	if readBackoffMin := readerConfigJs.Get("readBackoffMin"); !readBackoffMin.IsUndefined() {
		kafkaReaderConfig.ReadBackoffMin = time.Millisecond * time.Duration(readBackoffMin.Int())
	}

	if readBackoffMax := readerConfigJs.Get("readBackoffMax"); !readBackoffMax.IsUndefined() {
		kafkaReaderConfig.ReadBackoffMax = time.Millisecond * time.Duration(readBackoffMax.Int())
	}

	if logger := readerConfigJs.Get("logger"); !logger.IsUndefined() && logger.Bool() {
		kafkaReaderConfig.Logger = log.Default()
	}

	if maxAttempts := readerConfigJs.Get("maxAttempts"); !maxAttempts.IsUndefined() {
		kafkaReaderConfig.MaxAttempts = maxAttempts.Int()
	}

	kafkaReader := kafka.NewReader(kafkaReaderConfig)

	return (&reader{
		underlying: kafkaReader,
	}).toJSObject()

})
