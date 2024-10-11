package mocklistener

import (
	"net"
)

var _ net.Listener = (*Listener)(nil)

type Listener struct {
	// Address is the address of the listener
	Address net.Addr
	// ConnCh must be closed after Close manually
	ConnCh chan net.Conn
}

func (ln *Listener) Accept() (net.Conn, error) {
	return <-ln.ConnCh, nil
}

func (ln *Listener) Close() error {
	return nil
}

func (ln *Listener) Addr() net.Addr {
	return ln.Address
}
