package stackvm

import "errors"

var (
	// ErrIllegalState is returned when the requested action is not allowed in the current state.
	ErrIllegalState = errors.New("illegal state")

	// ErrInvalidProgram is returned when the program is invalid.
	ErrInvalidProgram = errors.New("invalid program")

	// ErrStackOverflow is returned when the stack is full.
	ErrStackOverflow = errors.New("stack overflow")

	// ErrStackUnderflow is returned when the stack is empty.
	ErrStackUnderflow = errors.New("stack underflow")

	// ErrTypeMismatch is returned when the type of the value is not expected.
	ErrTypeMismatch = errors.New("type mismatch")
)
