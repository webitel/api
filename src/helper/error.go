package helper

type CodeError struct {
	Code int
	error
}

func NewCodeError(code int, err error) *CodeError {
	return &CodeError{code, err}
}
