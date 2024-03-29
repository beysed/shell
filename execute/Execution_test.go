package execute

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestOutput(t *testing.T) {
	var command string
	var args []string

	if runtime.GOOS == "windows" {
		command = "cmd.exe"
		args = []string{"/c", "dir"}
	} else {
		command = "bash"
		args = []string{"-c", "ls -la"}
	}

	exe, _ := Execute(command, args...)

	for run := true; run; {
		select {
		case o := <-exe.Stdout:
			fmt.Print(string(o))

		case <-exe.Exit:
			run = false

		case <-time.After(time.Second * 3):
			t.Error("exit not fired, process hags, timeout")
			run = false
		}
	}
}
