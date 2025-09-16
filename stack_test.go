package stackvm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	stack := newStack(2)

	err1 := stack.push(NewString("Hello"))
	err2 := stack.push(NewString("World"))
	require.NoError(t, err1)
	require.NoError(t, err2)

	item, err := stack.pop()
	require.NoError(t, err)
	assert.Equal(t, "World", item.v)

	item, err = stack.pop()
	require.NoError(t, err)
	assert.Equal(t, "Hello", item.v)
}

func TestStack_Overflow(t *testing.T) {
	stack := newStack(1)

	err := stack.push(NewString("Hello"))
	require.NoError(t, err)

	err = stack.push(NewString("World"))
	assert.EqualError(t, err, ErrStackOverflow.Error())
}

func TestStack_UnderflowEmpty(t *testing.T) {
	stack := newStack(1)

	_, err := stack.pop()
	assert.EqualError(t, err, ErrStackUnderflow.Error())
}

func TestStack_UnderflowBase(t *testing.T) {
	stack := newStack(1)
	stack.push(NewString("Hello"))
	stack.newFrame(&FuncProto{nargs: 0})

	_, err := stack.pop()
	assert.EqualError(t, err, ErrStackUnderflow.Error())
}
