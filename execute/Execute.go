package execute

import (
	"io"
	"os/exec"
)

// ^C 0x03
// ^D 0x04
// ^Z 0x1a

func Execute(file string, args ...string) (Execution, error) {
	cmd := exec.Command(file, args...)

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

	forward := func(pipe io.ReadCloser, ch chan<- []byte) {
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
	go forward(stdout, ch_stdout)
	go forward(stderr, ch_stderr)

	go func() {
		for {
			input, ok := <-ch_stdin
			to_write := len(input)

			if input != nil {
				for to_write > 0 {
					wrote, err := stdin.Write(input[len(input)-to_write:])
					if err != nil {
						break
					}
					to_write -= wrote
				}
			}

			if !ok {
				stdin.Close()
				break
			}
		}
	}()

	go func() {
		err = cmd.Wait()
		ch_exit <- err
	}()

	return execution, nil
}
