package stackvm

// Function is a function that can be executed by the virtual machine.
type Function struct {
	proto *FuncProto
}

// FuncProto is a function prototype.
type FuncProto struct {
	nargs     int
	bytecode  []Inst
	constPool map[string]Inst
}

// FuncProtoBuilder is a builder for function prototypes.
type FuncProtoBuilder struct {
	bytecode  []Inst
	constPool map[string]Inst
}

// NewFuncProto creates a new function prototype.
func NewFuncProto(f func(*FuncProtoBuilder)) *FuncProto {
	builder := &FuncProtoBuilder{}
	f(builder)
	return builder.build()
}

// Emit emits an instruction to the bytecode.
func (b *FuncProtoBuilder) Emit(inst Inst) InstPtr {
	b.bytecode = append(b.bytecode, inst)
	return InstPtr(len(b.bytecode) - 1)
}

func (b *FuncProtoBuilder) build() *FuncProto {
	return &FuncProto{
		bytecode:  b.bytecode,
		constPool: b.constPool,
	}
}
