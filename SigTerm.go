package shell

import (
	"os"
	"os/signal"
	"syscall"
)

func CaptureSigTerm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		os.Exit(1)
	}()
}
