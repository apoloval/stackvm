package stackvm_test

import (
	"strconv"
	"testing"

	"github.com/apoloval/stackvm"
	"github.com/stretchr/testify/require"
)

func TestVM(t *testing.T) {
	for _, test := range []struct {
		name    string
		code    func(b *stackvm.FuncProtoBuilder)
		samples []funcSample
	}{
		{
			name: "sum(a,b,c:int)->int",
			samples: []funcSample{
				{
					args: []stackvm.Value{
						stackvm.NewInt(1),
						stackvm.NewInt(2),
						stackvm.NewInt(3),
					},
					expected: []stackvm.Value{stackvm.NewInt(6)},
				},
				{
					args: []stackvm.Value{
						stackvm.NewInt(4),
						stackvm.NewInt(5),
						stackvm.NewInt(6),
					},
					expected: []stackvm.Value{stackvm.NewInt(15)},
				},
			},
			code: func(b *stackvm.FuncProtoBuilder) {
				b.Emit(stackvm.ADDI())
				b.Emit(stackvm.ADDI())
				b.Emit(stackvm.RET(1))
			},
		},
		{
			name: "max(a,b:int)->int",
			samples: []funcSample{
				{
					args: []stackvm.Value{
						stackvm.NewInt(123),
						stackvm.NewInt(456),
					},
					expected: []stackvm.Value{stackvm.NewInt(456)},
				},
				{
					args: []stackvm.Value{
						stackvm.NewInt(456),
						stackvm.NewInt(123),
					},
					expected: []stackvm.Value{stackvm.NewInt(456)},
				},
			},
			code: func(b *stackvm.FuncProtoBuilder) {
				b.Emit(stackvm.DUP(0))
				b.Emit(stackvm.DUP(1))
				b.Emit(stackvm.GTI())
				gt := b.NewLabel()
				b.EmitBranch(gt)
				b.Emit(stackvm.DUP(1))
				b.Emit(stackvm.RET(1))
				b.Mark(gt)
				b.Emit(stackvm.DUP(0))
				b.Emit(stackvm.RET(1))
			},
		},
		{
			name: "abs(a:float)->float",
			samples: []funcSample{
				{
					args:     []stackvm.Value{stackvm.NewFloat(-3.1416)},
					expected: []stackvm.Value{stackvm.NewFloat(3.1416)},
				},
				{
					args:     []stackvm.Value{stackvm.NewFloat(3.1416)},
					expected: []stackvm.Value{stackvm.NewFloat(3.1416)},
				},
			},
			code: func(b *stackvm.FuncProtoBuilder) {
				end := b.NewLabel()
				b.Emit(stackvm.DUP(0))
				b.Emit(stackvm.PUSHF(0.0))
				b.Emit(stackvm.GEF())
				b.EmitBranch(end)
				b.Emit(stackvm.NEGF())
				b.Mark(end)
				b.Emit(stackvm.RET(1))
			},
		},
		{
			name: "clamp(x, min, max: float) -> float",
			samples: []funcSample{
				{
					args:     []stackvm.Value{stackvm.NewFloat(1.0), stackvm.NewFloat(0.0), stackvm.NewFloat(2.0)},
					expected: []stackvm.Value{stackvm.NewFloat(1.0)},
				},
				{
					args:     []stackvm.Value{stackvm.NewFloat(3.0), stackvm.NewFloat(0.0), stackvm.NewFloat(2.0)},
					expected: []stackvm.Value{stackvm.NewFloat(2.0)},
				},
				{
					args:     []stackvm.Value{stackvm.NewFloat(-1.0), stackvm.NewFloat(0.0), stackvm.NewFloat(2.0)},
					expected: []stackvm.Value{stackvm.NewFloat(0.0)},
				},
			},
			code: func(b *stackvm.FuncProtoBuilder) {
				checkMax := b.NewLabel()
				end := b.NewLabel()
				b.Emit(stackvm.DUP(0)) // x >= min ?
				b.Emit(stackvm.DUP(1))
				b.Emit(stackvm.GEF())
				b.EmitBranch(checkMax)
				b.Emit(stackvm.DUP(1))
				b.Emit(stackvm.RET(1))
				b.Mark(checkMax)
				b.Emit(stackvm.DUP(0))
				b.Emit(stackvm.DUP(2))
				b.Emit(stackvm.LEF())
				b.EmitBranch(end)
				b.Emit(stackvm.DUP(2))
				b.Emit(stackvm.RET(1))
				b.Mark(end)
				b.Emit(stackvm.DUP(0))
				b.Emit(stackvm.RET(1))
			},
		},
		{
			name: "avg(a, b, c: float) -> float",
			samples: []funcSample{
				{
					args:     []stackvm.Value{stackvm.NewFloat(1.0), stackvm.NewFloat(2.0), stackvm.NewFloat(3.0)},
					expected: []stackvm.Value{stackvm.NewFloat(2.0)},
				},
				{
					args:     []stackvm.Value{stackvm.NewFloat(4.0), stackvm.NewFloat(5.0), stackvm.NewFloat(6.0)},
					expected: []stackvm.Value{stackvm.NewFloat(5.0)},
				},
			},
			code: func(b *stackvm.FuncProtoBuilder) {
				b.Emit(stackvm.DUP(0))
				b.Emit(stackvm.DUP(1))
				b.Emit(stackvm.ADDF())
				b.Emit(stackvm.DUP(2))
				b.Emit(stackvm.ADDF())
				b.Emit(stackvm.PUSHF(3.0))
				b.Emit(stackvm.DIVF())
				b.Emit(stackvm.RET(1))
			},
		},
	} {
		for i, sample := range test.samples {
			t.Run(test.name+"[sample:"+strconv.Itoa(i)+"]", func(t *testing.T) {
				vm := stackvm.New()
				prog, err := stackvm.NewFuncProto(len(sample.args), test.code)
				require.NoError(t, err)
				values, err := vm.Run(prog, sample.args...)
				require.NoError(t, err)
				require.Equal(t, len(sample.expected), len(values))
				for i, value := range values {
					require.Equal(t, sample.expected[i], value)
				}
			})
		}
	}
}

type funcSample struct {
	args     []stackvm.Value
	expected []stackvm.Value
}
