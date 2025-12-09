package pkg

import (
	"net"
	"time"
)

func WaitFor(addr string, timeout time.Duration) bool {
	portAvailable := make(chan struct{})
	timeoutCh := time.After(timeout)

	go func() {
		for {
			select {
			case <-timeoutCh:
				return
			default:
			}

			_, err := net.Dial("tcp", addr)
			if err == nil {
				close(portAvailable)
				return
			}
			time.Sleep(200 * time.Millisecond)
		}

	}()

	select {
	case <-portAvailable:
		return true
	case <-timeoutCh:
		return false
	}
}
