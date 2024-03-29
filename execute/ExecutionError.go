package execute

type ExecutionError struct {
	ErrorText string
}

func (a ExecutionError) Error() string {
	return a.ErrorText
}

func MakeError(err string) error {
	return &ExecutionError{ErrorText: err}
}
