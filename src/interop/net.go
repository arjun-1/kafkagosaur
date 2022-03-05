package interop

import (
	"context"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"syscall"
	"syscall/js"
	"time"
)

type denoAddr struct {
	jsAddr js.Value
}

func (c *denoAddr) Network() string {
	return c.jsAddr.Get("network").String()
}

func (c *denoAddr) String() string {
	return c.jsAddr.Get("string").String()
}

func mapDeadlineError(reason js.Value) error {
	if reason := reason.Get("name"); !reason.IsUndefined() && reason.String() == "DeadlineError" {
		return os.ErrDeadlineExceeded
	}

	return defaultError(reason)
}

type denoTCPConn struct {
	jsConn    js.Value
	closeOnce sync.Once
}

func (c *denoTCPConn) Read(b []byte) (n int, err error) {

	jsBytes := NewUint8Array(len(b))
	value, err := AwaitWithErrorMapping(c.jsConn.Call("read", jsBytes), mapDeadlineError)

	if err != nil {
		return 0, err
	}
	if value.IsNull() {
		return 0, io.EOF
	}

	js.CopyBytesToGo(b, jsBytes)

	return value.Int(), nil
}

func (c *denoTCPConn) Write(b []byte) (n int, err error) {
	jsBytes := NewUint8Array(len(b))
	js.CopyBytesToJS(jsBytes, b)

	value, err := AwaitWithErrorMapping(c.jsConn.Call("write", jsBytes), mapDeadlineError)

	if err != nil {
		return 0, err
	}

	return value.Int(), nil
}

func (c *denoTCPConn) LocalAddr() net.Addr {
	jsAddr := c.jsConn.Get("localAddr")

	return &denoAddr{
		jsAddr: jsAddr,
	}
}

func (c *denoTCPConn) RemoteAddr() net.Addr {
	jsAddr := c.jsConn.Get("remoteAddr")

	return &denoAddr{
		jsAddr: jsAddr,
	}
}

func (c *denoTCPConn) SetDeadline(t time.Time) error {
	err := c.SetReadDeadline(t)

	if err != nil {
		return c.SetWriteDeadline(t)
	}

	return err
}

func (c *denoTCPConn) SetReadDeadline(t time.Time) error {
	c.jsConn.Call("setReadDeadline", t.UnixMilli())
	return nil
}

func (c *denoTCPConn) SetWriteDeadline(t time.Time) error {
	c.jsConn.Call("setWriteDeadline", t.UnixMilli())
	return nil
}

func (c *denoTCPConn) Close() error {
	c.closeOnce.Do(func() {
		Await(c.jsConn.Call("close"))
	})

	return nil
}

type DialBackend int

const (
	DenoDialBackend DialBackend = iota
	NodeDialBackend DialBackend = iota
)

var dialBackendMap = map[string]DialBackend{
	"deno": DenoDialBackend,
	"node": NodeDialBackend,
}

func StringToDialBackend(str string) DialBackend {
	return dialBackendMap[str]
}

type DialFunc func(ctx context.Context, network string, address string) (net.Conn, error)

const globalDenoInteropKey = "kafkagosaur.deno"

func NewDenoConn(dialBackend DialBackend) func(ctx context.Context, network string, address string) (net.Conn, error) {
	return func(ctx context.Context, network string, address string) (net.Conn, error) {

		if network != "tcp" {
			return nil, syscall.ENOTSUP
		}

		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}

		portInt, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}

		var jsTCPConn js.Value

		switch dialBackend {
		case DenoDialBackend:
			jsTCPConn, err = Await(js.Global().Get(globalDenoInteropKey).Call("dialDeno", host, portInt))
		case NodeDialBackend:
			jsTCPConn, err = Await(js.Global().Get(globalDenoInteropKey).Call("dialNode", host, portInt))
		}

		if err != nil {
			return nil, err
		}

		denoConn := &denoTCPConn{
			jsConn: jsTCPConn,
		}

		return denoConn, nil
	}
}
