package execute

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	var command Command

	setup := func(c *exec.Cmd) {
		c.Env = []string{"MY=MY"}
	}

	if runtime.GOOS == "windows" {
		command = MakeCommand("cmd.exe", "/c", "echo %MY%")
	} else {
		command = MakeCommand("bash", "-c", "printenv MY")
	}

	exe, _ := Execute(command, setup)

	hasOutput := false
	for run := true; run; {
		select {
		case o := <-exe.Stdout:
			assert.True(t, strings.HasPrefix(string(o), "MY"))
			hasOutput = true
		case o := <-exe.Stderr:
			fmt.Println(string(o))
		case <-exe.Exit:
			run = false
		case <-time.After(time.Second * 3):
			t.Error("exit not fired, process hags, timeout")
			run = false
		}
	}
	assert.True(t, hasOutput)
}

func TestConversationSed(t *testing.T) {
	exe, _ := Execute(MakeCommand("sed", "-e", "s/s/S/"))

	exe.Stdin <- []byte("sss")
	close(exe.Stdin)

	var output string
	for run := true; run; {
		select {
		case o := <-exe.Stdout:
			output += string(o)
		case o := <-exe.Stderr:
			fmt.Println(string(o))
		case <-exe.Exit:
			run = false

		case <-time.After(time.Second * 3):
			t.Error("exit not fired, process hags, timeout")
			run = false
		}
	}

	assert.Equal(t, "Sss", string(output))
}

func TestConversationCat(t *testing.T) {
	exe, _ := Execute(MakeCommand("cat"))

	exe.Stdin <- []byte{97, 97, 97}
	close(exe.Stdin)

	var output string
	for run := true; run; {
		select {
		case o := <-exe.Stdout:
			output += string(o)
		case o := <-exe.Stderr:
			fmt.Println(string(o))
		case <-exe.Exit:
			run = false

		case <-time.After(time.Second * 3):
			t.Error("exit not fired, process hags, timeout")
			run = false
		}
	}

	assert.Equal(t, "aaa", string(output))
}

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

	exe, _ := Execute(MakeCommand(command, args...))

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
