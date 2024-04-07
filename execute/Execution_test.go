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

func execute(t *testing.T, exe Execution) (string, bool) {
	hasOutput := false
	output := strings.Builder{}

	run := true
	for run {
		select {
		case o := <-exe.Stdout:
			hasOutput = true
			output.Write(o)
		case o := <-exe.Stderr:
			hasOutput = true
			output.Write(o)
		case <-exe.Exit:
			run = false
		case <-time.After(time.Second * 3):
			exe.Kill()
			t.Error("exit was not fired, process hags, timeout")
			t.Fail()

			run = false
		}
	}

	return output.String(), hasOutput
}

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

	exe, err := Execute(command, setup)
	assert.Nil(t, err)
	output, hasOutput := execute(t, exe)

	assert.True(t, hasOutput)
	assert.True(t, strings.HasPrefix(output, "MY"))
}

func TestConversationSed(t *testing.T) {
	exe, err := Execute(MakeCommand("sed", "-e", "s/s/S/"))

	assert.Nil(t, err)

	exe.Stdin <- []byte("sss")
	close(exe.Stdin)

	output, hasOutput := execute(t, exe)

	assert.True(t, hasOutput)
	assert.Equal(t, "Sss", string(output))
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
	_, hasOutput := execute(t, exe)
	assert.True(t, hasOutput)
}

func TestCheckOutput(t *testing.T) {
	var command Command
	testString := "00000000001111111111222222222233333333334444444444555555555566666666667777777777"
	echoCmd := fmt.Sprintf("echo %s", testString)
	if runtime.GOOS == "windows" {
		command = MakeCommand("cmd.exe", "/c", echoCmd)
	} else {
		command = MakeCommand("bash", "-c", echoCmd)
	}

	exe, err := Execute(command)
	assert.Nil(t, err)
	output, hasOutput := execute(t, exe)
	assert.True(t, hasOutput)
	assert.True(t, strings.HasPrefix(output, testString))
}
