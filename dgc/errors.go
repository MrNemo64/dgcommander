package dgc

import "fmt"

type DGCError struct {
	msg  string
	args []any
}

func makeError(msg string) DGCError {
	return DGCError{msg: msg}
}

func (err DGCError) withArgs(arg ...any) DGCError {
	copy := err
	copy.args = arg
	return copy
}

func (err DGCError) Error() string {
	return fmt.Errorf(err.msg, err.args...).Error()
}

func (err DGCError) Is(other error) bool {
	casted, ok := other.(DGCError)
	return ok && casted.msg == err.msg
}

func (err DGCError) Unwrap() []error {
	if err.args == nil {
		return []error{}
	}
	var errors []error
	for _, arg := range err.args {
		if e, ok := arg.(error); ok {
			errors = append(errors, e)
		}
	}
	return errors
}
