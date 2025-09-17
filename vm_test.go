package stackvm_test

import (
	"testing"

	"github.com/apoloval/stackvm"
	"github.com/stretchr/testify/require"
)

func TestVM(t *testing.T) {
	for _, test := range []struct {
		name     string
		args     []stackvm.Value
		code     func(b *stackvm.FuncProtoBuilder)
		expected []stackvm.Value
	}{
		{
			name: "func add(a, b, c: int) -> int",
			args: []stackvm.Value{
				stackvm.NewInt(1),
				stackvm.NewInt(2),
				stackvm.NewInt(3),
			},
			code: func(b *stackvm.FuncProtoBuilder) {
				b.Emit(stackvm.ADDI())
				b.Emit(stackvm.ADDI())
				b.Emit(stackvm.RET(1))
			},
			expected: []stackvm.Value{stackvm.NewInt(6)},
		},
		{
			name: "func max(a, b: int) -> bool",
			args: []stackvm.Value{
				stackvm.NewInt(123),
				stackvm.NewInt(456),
			},
			code: func(b *stackvm.FuncProtoBuilder) {
				b.Emit(stackvm.DUP(0))
				b.Emit(stackvm.DUP(1))
				b.Emit(stackvm.EVIGT())
				gt := b.NewLabel()
				b.EmitWithFixup(stackvm.BR(0), gt)
				b.Emit(stackvm.DUP(1))
				b.Emit(stackvm.RET(1))
				b.Mark(gt)
				b.Emit(stackvm.DUP(0))
				b.Emit(stackvm.RET(1))
			},
			expected: []stackvm.Value{stackvm.NewInt(456)},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			vm := stackvm.New()
			prog, err := stackvm.NewFuncProto(len(test.args), test.code)
			require.NoError(t, err)
			values, err := vm.Run(prog, test.args...)
			require.NoError(t, err)
			require.Equal(t, len(test.expected), len(values))
			for i, value := range values {
				require.Equal(t, test.expected[i], value)
			}
		})
	}
}
