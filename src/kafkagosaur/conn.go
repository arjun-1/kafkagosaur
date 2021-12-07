package kafkagosaur

import (
	"fmt"
	"github.com/arjun-1/kafkagosaur/src/interop"
	"github.com/segmentio/kafka-go"
	"syscall/js"
)

type conn struct {
	underlying *kafka.Conn
}

func (c *conn) apiVersionsPromise() js.Value {
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

func (c *conn) closePromise() js.Value {
	return interop.NewPromise(func(resolve func(interface{}), reject func(error)) {
		err := c.underlying.Close()

		if err != nil {
			reject(err)
		}

		resolve(nil)
	})
}

func (c *conn) toJSObject() map[string]interface{} {
	return map[string]interface{}{
		"apiVersions": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return c.apiVersionsPromise()
			},
		),
		"close": js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				return c.closePromise()
			},
		),
	}
}
