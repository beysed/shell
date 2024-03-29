package shell

import (
	"context"
	"time"
)

func Delay(t time.Duration, c context.Context) <-chan error {
	ch := make(chan error)

	go func() {
		select {
		case <-time.After(t):

			ch <- nil
		case <-c.Done():
			ch <- nil
		}
	}()

	return ch
}
