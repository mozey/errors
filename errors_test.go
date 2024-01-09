package errors_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/mozey/errors"
)

var packageName = "errors_test"

var msgEmpty = "empty cause"
var msgOne = "first message"
var msgTwo = "second message with args %s %d"
var msgSomething = "something"
var msgThree = fmt.Sprintf("%s bla bla", msgSomething)
var msgFour = "another thing"

// ErrTest uses the package name as the message value
var ErrTest = errors.NewCause(packageName)

// ErrEmpty is a custom error with empty cause
var ErrEmpty = errors.New(msgEmpty)

// ErrOne is a custom error with fixed message
var ErrOne = errors.NewWithCause(ErrTest, msgOne)

// ErrTwo customizes the error message with args
var ErrTwo = func(foo string, bar int) error {
	return errors.NewWithCausef(
		ErrTest, msgTwo, foo, bar)
}

// ErrCheck sets the cause of an error if it matches something,
// or otherwise it returns the original error
var ErrCheck = func(err error) error {
	if strings.Contains(err.Error(), msgSomething) {
		// Replace error
		return errors.NewWithCause(ErrTest, err.Error())
	}
	// Return original error
	return err
}

func TestErrors(t *testing.T) {
	is := is.New(t)

	err1 := ErrOne
	err2 := ErrTwo("a", 1)
	err3 := fmt.Errorf(msgThree)
	err4 := fmt.Errorf(msgFour)

	// Is matches on the error value, i.e the message
	is.True(errors.Is(err1, ErrOne))
	is.True(err1 == ErrOne) // Comparison on value also matches

	// Args don't matter when matching.
	// Matching is done on the code before it's used to format the message
	is.True(errors.Is(err2, ErrTwo("", 0)))
	is.True(err2 != ErrTwo("", 0)) // Comparison on value doesn't match

	// Errors with different message codes don't match
	is.True(!errors.Is(err1, ErrTwo("", 0)))
	is.True(!errors.Is(err2, ErrOne))

	// Error without a cause matches errors.EmptyCause
	is.True(errors.Is(errors.Cause(ErrEmpty), errors.EmptyCause))

	// Errors with cause match ErrTest
	is.True(errors.Is(errors.Cause(err1), ErrTest))
	is.True(errors.Is(errors.Cause(err2), ErrTest))

	// Message for cause is set on the error context (custom type).
	// Error without a cause has an empty message.
	// Errors with cause all have the same message
	is.True(errors.Cause(ErrEmpty).Error() == errors.EmptyCause.Error())
	is.True(errors.Cause(err1).Error() == packageName)
	is.True(errors.Cause(err2).Error() == packageName)

	// Pattern for marking the cause of third-party errors
	is.True(errors.Is(errors.Cause(ErrCheck(err3)), ErrTest))
	// Error and cause message is the same,
	// err3 is not a custom error and does not implement causer
	is.True(!errors.Is(errors.Cause(err3), errors.EmptyCause))
	is.True(errors.Cause(err3).Error() == msgThree)
	// ErrCheck does not set the cause if the substring is not found
	is.True(!errors.Is(errors.Cause(ErrCheck(err4)), ErrTest))
}
