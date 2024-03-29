package shell

import (
	"context"
	"testing"
	"time"
)

func Test_Delay(t *testing.T) {
	context, _ := context.WithCancel(context.Background())

	select {
	case <-Delay(time.Second, context):
	case <-time.After(time.Second * 2):
		t.Error("delay hags, timeout")
	}
}

func Test_Delay_Cancel(t *testing.T) {
	context, cancel := context.WithCancel(context.Background())

	ch := Delay(time.Second*3, context)
	cancel()

	select {
	case <-ch:
	case <-time.After(time.Second * 2):
		t.Error("delay hags, timeout")
	}
}
