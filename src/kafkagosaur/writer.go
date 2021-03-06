package kafkagosaur

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"syscall/js"

	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
)

type writer struct {
	underlying *kafka.Writer
	transport  *kafka.Transport
}

func (w *writer) writeMessages(msgs []kafka.Message) js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {

		switch err := w.underlying.WriteMessages(context.Background(), msgs...).(type) {
		case nil:
		case kafka.WriteErrors:
			writeErrors := []string{err.Error()}

			for i := range msgs {
				if err[i] != nil && len(writeErrors) > 3 {
					writeErrors = append(writeErrors, "...")
					break
				} else if err[i] != nil {
					writeErrors = append(writeErrors, fmt.Sprintf("message %d: %s", i, err[i].Error()))
				}
			}

			reject(errors.New(strings.Join(writeErrors, "\n")))
		default:
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

func (w *writer) stats() js.Value {
	return js.ValueOf(fmt.Sprintf("%+v", w.underlying.Stats()))
}

func (w *writer) toJSObject() map[string]interface{} {
	return map[string]interface{}{
		"writeMessages": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {

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
		"stats": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return w.stats()
			},
		),
	}
}

var NewWriterJsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	writerConfig := args[0]

	saslMechanism, err := SASLMechanism(writerConfig)
	if err != nil {
		panic(err)
	}

	tls, err := TLSConfig(writerConfig)
	if err != nil {
		panic(err)
	}

	var dialBackend = interop.NodeDialBackend

	if dialBackendJs := writerConfig.Get("dialBackend"); !dialBackendJs.IsUndefined() {
		dialBackend = interop.StringToDialBackend(dialBackendJs.String())
	}

	transport := &kafka.Transport{
		Dial: interop.NewDenoConn(dialBackend),
		SASL: saslMechanism,
		TLS:  tls,
	}

	if dialTimeout := writerConfig.Get("dialTimeout"); !dialTimeout.IsUndefined() {
		transport.DialTimeout = JsNumberMillisToDuration(dialTimeout)
	}

	if jsIdleTimeout := writerConfig.Get("idleTimeout"); !jsIdleTimeout.IsUndefined() {
		transport.IdleTimeout = JsNumberMillisToDuration(jsIdleTimeout)
	}

	if metadataTTL := writerConfig.Get("metadataTTL"); !metadataTTL.IsUndefined() {
		transport.MetadataTTL = JsNumberMillisToDuration(metadataTTL)
	}

	if clientId := writerConfig.Get("clientId"); !clientId.IsUndefined() {
		transport.ClientID = clientId.String()
	}

	kafkaWriter := kafka.Writer{
		Addr:      kafka.TCP(writerConfig.Get("address").String()),
		Transport: transport,
	}

	if jsTopic := writerConfig.Get("topic"); !jsTopic.IsUndefined() {
		kafkaWriter.Topic = jsTopic.String()
	}

	if maxAttempts := writerConfig.Get("maxAttempts"); !maxAttempts.IsUndefined() {
		kafkaWriter.MaxAttempts = maxAttempts.Int()
	}

	if batchSize := writerConfig.Get("batchSize"); !batchSize.IsUndefined() {
		kafkaWriter.BatchSize = batchSize.Int()
	}

	if batchBytes := writerConfig.Get("batchBytes"); !batchBytes.IsUndefined() {
		kafkaWriter.BatchBytes = int64(batchBytes.Int())
	}

	if batchTimeout := writerConfig.Get("batchTimeout"); !batchTimeout.IsUndefined() {
		kafkaWriter.BatchTimeout = JsNumberMillisToDuration(batchTimeout)
	}

	if readTimeout := writerConfig.Get("readTimeout"); !readTimeout.IsUndefined() {
		kafkaWriter.ReadTimeout = JsNumberMillisToDuration(readTimeout)
	}

	if writeTimeout := writerConfig.Get("writeTimeout"); !writeTimeout.IsUndefined() {
		kafkaWriter.WriteTimeout = JsNumberMillisToDuration(writeTimeout)
	}

	if async := writerConfig.Get("async"); !async.IsUndefined() {
		kafkaWriter.Async = async.Bool()
	}

	if logger := writerConfig.Get("logger"); !logger.IsUndefined() && logger.Bool() {
		kafkaWriter.Logger = log.Default()
	}

	return (&writer{
		underlying: &kafkaWriter,
		transport:  transport,
	}).toJSObject()

})
