package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	result := fmt.Sprintf("%d errors occured:\n", len(e.errors))

	for _, err := range e.errors {
		result += fmt.Sprintf("\t* %s", err.Error())
	}
	return result + "\n"
}

func (e *MultiError) Unwrap() []error {
	return e.errors
}

func Append(err error, errs ...error) *MultiError {
	switch err.(type) {
	case *MultiError:
		err.(*MultiError).errors = append(err.(*MultiError).errors, errs...)

		return err.(*MultiError)
	default:
		e := make([]error, 0, len(errs)+1)
		if err != nil {
			e = append(e, err)
		}

		for _, err1 := range errs {
			if err1 != nil {
				e = append(e, err1)
			}
		}

		if len(e) == 0 {
			return nil
		}

		return &MultiError{errors: e}
	}
}

type TestError struct{}

func (e TestError) Error() string {
	return "error 2"
}

func TestMultiError(t *testing.T) {
	var err error
	e1 := errors.New("error 1")
	e2 := TestError{}
	e3 := errors.New("error 3")

	err = Append(err, e1)
	err = Append(err, e2)

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	assert.Equal(t, errors.Is(err, e1), true)
	assert.Equal(t, errors.Is(err, e2), true)
	assert.Equal(t, errors.Is(err, e3), false)

	assert.Equal(t, errors.As(err, &TestError{}), true)
}
