package jskafka

import (
	"fmt"
	"github.com/arjun-1/deno-wasm-experiment/interop"
	"github.com/segmentio/kafka-go"
	"syscall/js"
)

type jsConn struct {
	underlying *kafka.Conn
}

func (c *jsConn) apiVersionsPromise() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		apiVersions, err := c.underlying.ApiVersions()

		if err != nil {
			reject(err)
		}

		apiVersionsString := make([]interface{}, len(apiVersions))

		for i, v := range apiVersions {
			apiVersionsString[i] = fmt.Sprint(v)
		}

		resolve(apiVersionsString)

	})
}

func (c *jsConn) toJSObject() map[string]interface{} {
	return map[string]interface{}{
		"apiVersions": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return c.apiVersionsPromise()
			},
		),
	}
}
