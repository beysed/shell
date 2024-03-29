package execute

type Execution struct {
	Stdout <-chan []byte
	Stderr <-chan []byte
	Stdin  chan<- []byte
	Exit   <-chan error
	Kill   func() error
}
