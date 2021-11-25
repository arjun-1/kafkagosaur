package interop

import (
	"context"
	"io"
	"net"
	"strconv"
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

type denoConn struct {
	jsConn js.Value
}

func (c *denoConn) Close() error {
	c.jsConn.Call("close")
	return nil
}

func (c *denoConn) Read(b []byte) (n int, err error) {

	jsBytes := NewUint8Array(len(b))
	value, err := Await(c.jsConn.Call("read", jsBytes))
	js.CopyBytesToGo(b, jsBytes)

	if err != nil {
		return 0, err
	}
	if value.IsNull() {
		return 0, io.EOF
	}

	return value.Int(), nil
}

func (c *denoConn) Write(b []byte) (n int, err error) {
	jsBytes := NewUint8Array(len(b))
	js.CopyBytesToJS(jsBytes, b)

	value, err := Await(c.jsConn.Call("write", jsBytes))

	if err != nil {
		return 0, err
	}

	return value.Int(), nil
}

func (c *denoConn) LocalAddr() net.Addr {
	jsAddr := c.jsConn.Get("localAddr")

	return &denoNetAddr{
		jsAddr: jsAddr,
	}
}

func (c *denoConn) RemoteAddr() net.Addr {
	jsAddr := c.jsConn.Get("remoteAddr")

	return &denoNetAddr{
		jsAddr: jsAddr,
	}
}

func (c *denoConn) SetDeadline(t time.Time) error {
	return syscall.ENOTSUP
}

func (c *denoConn) SetReadDeadline(t time.Time) error {
	return syscall.ENOTSUP
}

func (c *denoConn) SetWriteDeadline(t time.Time) error {
	return syscall.ENOTSUP
}

func NewDenoConn(ctx context.Context, network string, address string) (net.Conn, error) {

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	// TODO: only allow transport == 'tcp' ?
	connectOptions := map[string]interface{}{
		"transport": network,
		"hostname":  host,
		"port":      portInt,
	}

	jsConn, err := Await(js.Global().Get("Deno").Call("connect", connectOptions))

	if err != nil {
		return nil, err
	}

	denoConn := &denoConn{
		jsConn: jsConn,
	}

	return denoConn, nil
}
