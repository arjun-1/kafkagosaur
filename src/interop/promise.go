package interop

import (
	"errors"
	"fmt"
	"runtime/debug"
	"syscall/js"
)

func NewPromise(executor func(resolve func(interface{}), reject func(error))) js.Value {

	jsExecutor := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		resolve := func(value interface{}) { args[0].Invoke(value) }
		reject := func(reason error) { args[1].Invoke(reason.Error()) }

		go executor(resolve, reject)
		defer func() {
			if r := recover(); r != nil {
				reject(errors.New(fmt.Sprintln(debug.Stack())))
			}
		}()

		return nil
	})

	defer jsExecutor.Release()

	return js.Global().Get("Promise").New(jsExecutor)
}

func Await(promiseLike js.Value) (js.Value, error) {
	value := make(chan js.Value)
	defer close(value)

	reason := make(chan js.Value)
	defer close(reason)

	onFulfilled := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		value <- args[0]
		return nil
	})
	defer onFulfilled.Release()

	onRejected := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		reason <- args[0]
		return nil
	})
	defer onRejected.Release()

	promiseLike.Call("then", onFulfilled, onRejected)

	select {
	case v := <-value:
		return v, nil
	case r := <-reason:
		return js.Undefined(), errors.New(r.String())
	}

}
