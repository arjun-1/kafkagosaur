package kafkagosaur

import (
	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
	"syscall/js"
	"time"
)

func jsObjectToMessage(jsObject js.Value) kafka.Message {
	message := kafka.Message{}

	if topicJs := jsObject.Get("topic"); !topicJs.IsUndefined() {
		message.Topic = topicJs.String()
	}

	if partitionJs := jsObject.Get("partition"); !partitionJs.IsUndefined() {
		message.Partition = partitionJs.Int()
	}

	if offsetJs := jsObject.Get("offset"); !offsetJs.IsUndefined() {
		message.Offset = int64(offsetJs.Int())
	}

	if highWaterMark := jsObject.Get("highWaterMark"); !highWaterMark.IsUndefined() {
		message.HighWaterMark = int64(highWaterMark.Int())
	}

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

	if headersJs := jsObject.Get("headers"); !headersJs.IsUndefined() {
		headers := make([]kafka.Header, headersJs.Length())

		for i := range headers {
			headerJs := headersJs.Index(i)
			valueJs := headerJs.Get("value")
			value := make([]byte, valueJs.Length())
			js.CopyBytesToGo(value, valueJs)

			headers[i] = kafka.Header{
				Key:   headerJs.Get("Key").String(),
				Value: value,
			}
		}

		message.Headers = headers
	}

	if timeJs := jsObject.Get("time"); !timeJs.IsUndefined() {
		time := time.UnixMilli(int64(timeJs.Int()))
		message.Time = time
	}

	return message
}

func messageToJSObject(m kafka.Message) map[string]interface{} {
	key := interop.NewUint8Array(len(m.Key))
	js.CopyBytesToJS(key, m.Key)

	value := interop.NewUint8Array(len(m.Value))
	js.CopyBytesToJS(value, m.Value)

	headersJs := make([]interface{}, len(m.Headers))

	for i, header := range m.Headers {
		valueJs := interop.NewUint8Array(len(header.Value))
		js.CopyBytesToJS(valueJs, header.Value)

		headerJs := map[string]interface{}{
			"key":   header.Key,
			"value": valueJs,
		}

		headersJs[i] = headerJs
	}

	return map[string]interface{}{
		"topic":         m.Topic,
		"partition":     m.Partition,
		"offset":        m.Offset,
		"highWaterMark": m.HighWaterMark,
		"time":          m.Time.UnixMilli(),
		"key":           key,
		"headers":       headersJs,
		"value":         value,
	}
}
