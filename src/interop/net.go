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

type DialFunc func(ctx context.Context, network string, address string) (net.Conn, error)

type denoNetAddr struct {
	jsAddr js.Value
}

func (c *denoNetAddr) Network() string {
	return c.jsAddr.Get("transport").String()
}

func (c *denoNetAddr) String() string {
	// TODO: handle unix addr case
	return net.JoinHostPort(c.jsAddr.Get("hostname").String(), c.jsAddr.Get("port").String())
}

func mapDeadlineError(reason js.Value) error {
	if reason.Get("name").String() == "DeadlineError" {
		return os.ErrDeadlineExceeded
	}

	return defaultErrorFn(reason)
}

type denoTCPConn struct {
	jsConn    js.Value
	ctx       context.Context
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

	return &denoNetAddr{
		jsAddr: jsAddr,
	}
}

func (c *denoTCPConn) RemoteAddr() net.Addr {
	jsAddr := c.jsConn.Get("remoteAddr")

	return &denoNetAddr{
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
		c.jsConn.Call("close")
		Await(c.jsConn.Call("closeWrite"))
	})

	return nil
}

func NewDenoConn(ctx context.Context, network string, address string) (net.Conn, error) {

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

	connectOptions := map[string]interface{}{
		"transport": network,
		"hostname":  host,
		"port":      portInt,
	}

	jsTCPConn, err := Await(js.Global().Call("connectWithDeadline", connectOptions))

	if err != nil {
		return nil, err
	}

	denoConn := &denoTCPConn{
		jsConn: jsTCPConn,
		ctx:    ctx,
	}

	return denoConn, nil
}
