package interop

import (
	"syscall/js"
)

func MapToString(vs []js.Value) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = v.String()
	}
	return vsm
}

func ToSlice(v js.Value) []js.Value {
	length := v.Length()
	slice := make([]js.Value, length)

	for i := range slice {
		slice[i] = v.Index(i)
	}

	return slice
}

func NewUint8Array(length int) js.Value {
	return js.Global().Get("Uint8Array").New(length)
}
