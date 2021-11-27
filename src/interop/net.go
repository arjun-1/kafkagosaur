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

type denoTCPConn struct {
	jsConn js.Value
}

func (c *denoTCPConn) Close() error {
	c.jsConn.Call("close")
	return nil
}

func (c *denoTCPConn) Read(b []byte) (n int, err error) {

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

func (c *denoTCPConn) Write(b []byte) (n int, err error) {
	jsBytes := NewUint8Array(len(b))
	js.CopyBytesToJS(jsBytes, b)

	value, err := Await(c.jsConn.Call("write", jsBytes))

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
	return syscall.ENOTSUP
}

func (c *denoTCPConn) SetReadDeadline(t time.Time) error {
	return syscall.ENOTSUP
}

func (c *denoTCPConn) SetWriteDeadline(t time.Time) error {
	return syscall.ENOTSUP
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

	jsTCPConn, err := Await(js.Global().Get("Deno").Call("connect", connectOptions))

	if err != nil {
		return nil, err
	}

	denoConn := &denoTCPConn{
		jsConn: jsTCPConn,
	}

	return denoConn, nil
}
