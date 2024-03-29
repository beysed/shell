package execute

import (
	"io"
	"os/exec"
)

// ^C 0x03
// ^D 0x04
// ^Z 0x1a

type Command struct {
	File string
	Args []string
}

func MakeCommand(file string, args ...string) Command {
	return Command{File: file, Args: args}
}

// func MakeExecution(cmd exec.Cmd) Execution {
// }

func forwardRead(pipe io.ReadCloser, ch chan<- []byte) {
	buf := make([]byte, 64)
	for {
		l, err := pipe.Read(buf)
		if l != 0 {
			ch <- buf[:l]
		}

		if err != nil {
			break
		}
	}

}

func forwardWrite(ch <-chan []byte, in io.WriteCloser) {
	for {
		input, ok := <-ch
		to_write := len(input)

		if input != nil {
			for to_write > 0 {
				wrote, err := in.Write(input[len(input)-to_write:])
				if err != nil {
					break
				}
				to_write -= wrote
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

	err = cmd.Start()
	if err != nil {
		return Execution{}, err
	}

	ch_stdin := make(chan []byte)
	ch_stdout := make(chan []byte)
	ch_stderr := make(chan []byte)
	ch_exit := make(chan error)

	execution := Execution{
		Stderr: ch_stderr,
		Stdout: ch_stdout,
		Stdin:  ch_stdin,
		Exit:   ch_exit,
		Kill: func() error {
			if cmd.Process == nil {
				return MakeError("process is nil")
			}

			return cmd.Process.Kill()
		}}

	go forwardRead(stdout, ch_stdout)
	go forwardRead(stderr, ch_stderr)
	go forwardWrite(ch_stdin, stdin)

	go func() {
		err = cmd.Wait()
		ch_exit <- err
	}()

	return execution, nil
}
