package stackvm_test

import (
	"testing"

	"github.com/apoloval/stackvm"
	"github.com/stretchr/testify/require"
)

func TestVM(t *testing.T) {
	vm := stackvm.New()

	values, err := vm.Run(stackvm.NewFuncProto(func(b *stackvm.FuncProtoBuilder) {
		b.Emit(stackvm.PUSHI(1))
		b.Emit(stackvm.PUSHI(2))
		b.Emit(stackvm.PUSHI(3))
		b.Emit(stackvm.ADDI())
		b.Emit(stackvm.ADDI())
		b.Emit(stackvm.RET(1))
	}))
	require.NoError(t, err)
	require.Equal(t, 1, len(values))

	val, err := values[0].AsInt()
	require.NoError(t, err)
	require.Equal(t, 6, val)
}
