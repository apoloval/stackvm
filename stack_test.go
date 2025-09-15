package stackvm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	stack := newStack(2)

	err1 := stack.push("Hello")
	err2 := stack.push("World")
	require.NoError(t, err1)
	require.NoError(t, err2)

	item, err := stack.pop()
	require.NoError(t, err)
	require.Equal(t, "World", item)

	item, err = stack.pop()
	require.NoError(t, err)
	require.Equal(t, "Hello", item)
}

func TestStackOverflow(t *testing.T) {
	stack := newStack(1)

	err := stack.push("Hello")
	require.NoError(t, err)

	err = stack.push("World")
	assert.EqualError(t, err, ErrStackOverflow.Error())
}

func TestStackUnderflow(t *testing.T) {
	stack := newStack(1)

	_, err := stack.pop()
	assert.EqualError(t, err, ErrStackUnderflow.Error())
}
