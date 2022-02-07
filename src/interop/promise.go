package interop

import (
	"errors"
	"fmt"
	"runtime/debug"
	"syscall/js"
)

func NewPromise(executor func(_ func(interface{}), _ func(error))) js.Value {

	executorJsFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		resolve := func(value interface{}) { args[0].Invoke(value) }
		reject := func(reason error) {
			err := js.Global().Call("Error", reason.Error())

			args[1].Invoke(err)
		}

		go executor(resolve, reject)
		defer func() {
			if r := recover(); r != nil {
				reject(errors.New(fmt.Sprintln(debug.Stack())))
			}
		}()

		return nil
	})

	defer executorJsFunc.Release()

	return js.Global().Get("Promise").New(executorJsFunc)
}

func defaultError(reason js.Value) error {
	var text string
	if stack := reason.Get("stack"); !stack.IsUndefined() {
		text = stack.String()
	} else {
		text = js.Global().Get("JSON").Call("stringify", reason).String()
	}

	return errors.New(text)
}

func Await(promiseLike js.Value) (js.Value, error) {
	return AwaitWithErrorMapping(promiseLike, defaultError)
}

func AwaitWithErrorMapping(promiseLike js.Value, errorFn func(js.Value) error) (js.Value, error) {
	fullfilledChannel := make(chan js.Value)
	defer close(fullfilledChannel)

	rejectedChannel := make(chan js.Value)
	defer close(rejectedChannel)

	onFulfilled := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fullfilledChannel <- args[0]
		return nil
	})
	defer onFulfilled.Release()

	onRejected := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		rejectedChannel <- args[0]
		return nil
	})
	defer onRejected.Release()

	promiseLike.Call("then", onFulfilled, onRejected)

	select {
	case v := <-fullfilledChannel:
		return v, nil
	case r := <-rejectedChannel:
		return js.Undefined(), errorFn(r)
	}

}
