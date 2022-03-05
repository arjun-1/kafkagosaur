package kafkagosaur

import (
	"syscall/js"
	"time"
)

func JsNumberMillisToDuration(value js.Value) time.Duration {
	return time.Duration(value.Float() * float64(time.Millisecond))
}
