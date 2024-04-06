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

func execute(t *testing.T, exe Execution) (bool, string) {
	hasOutput := false
	output := strings.Builder{}

	for run := true; run; {
		select {
		case o := <-exe.Stdout:
			hasOutput = true
			output.Write(o)
		case o := <-exe.Stderr:
			fmt.Println(string(o))
		case <-exe.Exit:
			run = false
		case <-time.After(time.Second * 3):
			t.Error("exit not fired, process hags, timeout")
			exe.Kill()

			run = false
		}
	}

	return hasOutput, output.String()
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

	exe, _ := Execute(command, setup)
	hasOutput, output := execute(t, exe)

	assert.True(t, hasOutput)
	assert.True(t, strings.HasPrefix(output, "MY"))
}

func TestConversationSed(t *testing.T) {
	exe, err := Execute(MakeCommand("sed", "-e", "s/s/S/"))

	if err != nil {
		t.Error(err)
		return
	}

	exe.Stdin <- []byte("sss")
	close(exe.Stdin)

	hasOutput, output := execute(t, exe)

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
	hasOutput, _ := execute(t, exe)
	assert.True(t, hasOutput)
}

func TestCheckOutput(t *testing.T) {
	var command string
	var args []string

	command = "bash"
	testString := "00000000001111111111222222222233333333334444444444555555555566666666667777777777"

	args = []string{"-c", fmt.Sprintf("echo %s", testString)}

	exe, _ := Execute(MakeCommand(command, args...))
	hasOutput, output := execute(t, exe)
	assert.Equal(t, testString+"\n", output)
	assert.True(t, hasOutput)
}

func TestExample(t *testing.T) {
	command := MakeCommand("sed", "-e", "s/a/A/g")
	execution, err := Execute(command)
	if err != nil {
		t.Error(err)
		return
	}

	execution.Stdin <- []byte("aaa")
	close(execution.Stdin)

	for run := true; run; {
		select {
		case out := <-execution.Stdout:
			fmt.Print(string(out))
		case err := <-execution.Stderr:
			fmt.Println(string(err))
		case <-execution.Exit:
			run = false
		case <-time.After(time.Second * 3):
			t.Error("process killed by timeout")
			execution.Kill()
			run = false
		}
	}
}
