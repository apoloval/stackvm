package stackvm

import "fmt"

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

// FuncProtoLabel is a label in a function prototype.
type FuncProtoLabel int

// FuncProtoBuilder is a builder for function prototypes.
type FuncProtoBuilder struct {
	nargs     int
	bytecode  []Inst
	constPool map[string]Inst
	fixups    []fixup
}

// NewFuncProto creates a new function prototype.
func NewFuncProto(nargs int, f func(*FuncProtoBuilder)) (*FuncProto, error) {
	builder := &FuncProtoBuilder{
		nargs: nargs,
	}
	f(builder)
	return builder.build()
}

// NewLabel creates a new label to be used for jump instructions.
func (b *FuncProtoBuilder) NewLabel() FuncProtoLabel {
	b.fixups = append(b.fixups, fixup{
		refs:  nil,
		value: nullInstPtr,
	})
	return FuncProtoLabel(len(b.fixups) - 1)
}

// NewLabelFixed creates a new label fixed to the current next instruction address.
func (b *FuncProtoBuilder) NewLabelFixed() FuncProtoLabel {
	label := b.NewLabel()
	b.Mark(label)
	return label
}

// Mark marks a label in the bytecode.
func (b *FuncProtoBuilder) Mark(label FuncProtoLabel) InstPtr {
	instPtr := InstPtr(len(b.bytecode))
	b.fixups[label].value = instPtr
	return instPtr
}

// Emit emits an instruction to the bytecode.
func (b *FuncProtoBuilder) Emit(inst Inst) InstPtr {
	b.bytecode = append(b.bytecode, inst)
	return InstPtr(len(b.bytecode) - 1)
}

// EmitBranch emits a branch instruction to the bytecode with a label.
// This configures a fixup for the branch instruction to the given label.
// The label must be marked before the function proto is built.
func (b *FuncProtoBuilder) EmitBranch(to FuncProtoLabel) InstPtr {
	instPtr := b.Emit(BR(0))
	b.fixups[to].refs = append(b.fixups[to].refs, instPtr)
	return instPtr
}

func (b *FuncProtoBuilder) build() (*FuncProto, error) {
	for _, fixup := range b.fixups {
		if fixup.value == nullInstPtr {
			return nil, fmt.Errorf("%w: label not marked", ErrInvalidProgram)
		}
		for _, ref := range fixup.refs {
			b.bytecode[ref] = b.bytecode[ref].withOpInstPtr(fixup.value)
		}
	}
	return &FuncProto{
		nargs:     b.nargs,
		bytecode:  b.bytecode,
		constPool: b.constPool,
	}, nil
}

type fixup struct {
	refs  []InstPtr
	value InstPtr
}
