package state

import "fmt"

type ErrorKind string

const (
	ErrRecordNotFound ErrorKind = "ErrRecordNotFound"
	ErrRecordInvalid            = "ErrRecordInvalid"
)

type Error struct {
	Kind        ErrorKind
	Explanation string
}

func NewError(kind ErrorKind, a ...interface{}) *Error {
	exp := ""
	if len(a) == 1 {
		if first, ok := a[0].(string); ok {
			exp = first
		}
	}
	if len(a) > 1 {
		if first, ok := a[0].(string); ok {
			exp = fmt.Sprintf(first, a[1:]...)
		}
	}
	return &Error{kind, exp}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Explanation)
}

func ErrorRecordNotFound(a ...interface{}) *Error {
	return NewError(ErrRecordNotFound, a...)
}

func ErrorRecordInvalid(a ...interface{}) *Error {
	return NewError(ErrRecordInvalid, a...)
}
