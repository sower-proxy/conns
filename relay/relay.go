package relay

import (
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

func RelayTo(conn net.Conn, addr string) (dur time.Duration, err error) {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = net.JoinHostPort(addr, "80")
	}

	start := time.Now()
	rc, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return time.Since(start), err
	}
	defer rc.Close()

	err = Relay(conn, rc)
	return time.Since(start), err
}

func Relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()

	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()

	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
