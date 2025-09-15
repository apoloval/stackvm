package stackvm

import "errors"

var (
	// ErrStackOverflow is returned when the stack is full.
	ErrStackOverflow = errors.New("stack overflow")

	// ErrStackUnderflow is returned when the stack is empty.
	ErrStackUnderflow = errors.New("stack underflow")
)
