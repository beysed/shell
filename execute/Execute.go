package execute

import (
	"io"
	"os/exec"
)

// ^C 0x03
// ^D 0x04
// ^Z 0x1a
//

type Command struct {
	File string
	Args []string
}

func MakeCommand(file string, args ...string) Command {
	return Command{File: file, Args: args}
}

func forwardRead(pipe io.ReadCloser, ch chan<- []byte) {
	buf := make([]byte, 10)
	for {
		l, err := pipe.Read(buf)
		if l != 0 {
			c := make([]byte, l)
			copy(c, buf)
			ch <- c
		}

		if err != nil {
			break
		}
	}
}

func forwardWrite(ch <-chan []byte, in io.WriteCloser) {
	for {
		input, ok := <-ch

		if input != nil {
			toWrite := len(input)

			for toWrite > 0 {
				wrote, err := in.Write(input[len(input)-toWrite:])
				if err != nil {
					in.Close()
					break
				}
				toWrite -= wrote
			}
		}

		if !ok {
			in.Close()
			break
		}
	}
}

func Execute(command Command, setup ...func(e *exec.Cmd)) (Execution, error) {
	cmd := exec.Command(command.File, command.Args...)

	for _, v := range setup {
		v(cmd)
	}

	//cmd.CombinedOutput()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return Execution{}, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return Execution{}, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return Execution{}, err
	}

	chStdin := make(chan []byte)
	chStdout := make(chan []byte)
	chStderr := make(chan []byte)
	chExit := make(chan error)

	go forwardRead(stdout, chStdout)
	go forwardRead(stderr, chStderr)
	go forwardWrite(chStdin, stdin)

	execution := Execution{
		Stderr: chStderr,
		Stdout: chStdout,
		Stdin:  chStdin,
		Exit:   chExit,
		Kill: func() error {
			if cmd.Process == nil {
				return MakeError("process is nil")
			}

			return cmd.Process.Kill()
		}}

	err = cmd.Start()
	if err != nil {
		return Execution{}, err
	}

	go func() {
		err = cmd.Wait()
		stdout.Close()
		stderr.Close()

		chExit <- err
	}()

	return execution, nil
}
