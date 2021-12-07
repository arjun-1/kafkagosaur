package kafkagosaur

import (
	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
	"syscall/js"
	"time"
)

func jsObjectToMessage(jsObject js.Value) kafka.Message {
	message := kafka.Message{}

	if keyJs := jsObject.Get("key"); !keyJs.IsUndefined() {
		key := make([]byte, keyJs.Length())
		js.CopyBytesToGo(key, keyJs)
		message.Key = key
	}

	if valueJs := jsObject.Get("value"); !valueJs.IsUndefined() {
		value := make([]byte, valueJs.Length())
		js.CopyBytesToGo(value, valueJs)
		message.Value = value
	}

	if timeJs := jsObject.Get("time"); !timeJs.IsUndefined() {
		time := time.UnixMilli(int64(timeJs.Int()))
		message.Time = time
	}

	if topicJs := jsObject.Get("topic"); !topicJs.IsUndefined() {
		message.Topic = topicJs.String()
	}

	if partitionJs := jsObject.Get("partition"); !partitionJs.IsUndefined() {
		message.Partition = partitionJs.Int()
	}

	if offsetJs := jsObject.Get("offset"); !offsetJs.IsUndefined() {
		message.Offset = int64(offsetJs.Int())
	}

	return message
}

func messageToJSObject(m kafka.Message) map[string]interface{} {
	key := interop.NewUint8Array(len(m.Key))
	js.CopyBytesToJS(key, m.Key)

	value := interop.NewUint8Array(len(m.Value))
	js.CopyBytesToJS(value, m.Value)

	return map[string]interface{}{
		"offset":    m.Offset,
		"topic":     m.Topic,
		"partition": m.Partition,
		"time":      m.Time.UnixMilli(),
		"key":       key,
		"value":     value,
	}
}
