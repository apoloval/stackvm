package stackvm

type opCode uint8

const (
	opNop opCode = iota
	opAddi
	opPushi
	opRet
)

var opCodeToStr = map[opCode]string{
	opNop:   "NOP",
	opAddi:  "ADDI",
	opPushi: "PUSHI",
	opRet:   "RET",
}

// InstPtr is the pointer to the instruction.
type InstPtr int

// Inst is the instruction of the code segment. It is a 64-bit integer that encodes the
// operation in the first 8 bits and the arguments in the rest.
type Inst uint64

// ADDI encodes an ADDI instruction.
func ADDI() Inst {
	return encodeInst(opAddi)
}

// PUSHI encodes a PUSHI instruction.
func PUSHI(arg int) Inst {
	return encodeInstWithArgInt(opPushi, arg)
}

// RET encodes a RET instruction.
func RET(nres int) Inst {
	return encodeInstWithArgInt(opRet, nres)
}

func encodeInst(op opCode) Inst {
	return Inst(op) << 56
}

func encodeInstWithArgInt(op opCode, arg int) Inst {
	return Inst(op)<<56 | Inst(arg)
}

func (i Inst) opCode() opCode {
	return opCode(i >> 56)
}

func (i Inst) argInt() int {
	return int(i & 0xFFFF_FFFF)
}

func (i Inst) execute(vm *VirtualMachine) error {
	switch i.opCode() {
	case opNop:
		return nil
	case opAddi:
		a, err := vm.stack.popInt()
		if err != nil {
			return err
		}
		b, err := vm.stack.popInt()
		if err != nil {
			return err
		}
		return vm.stack.push(NewInt(a + b))
	case opPushi:
		return vm.stack.push(NewInt(i.argInt()))
	case opRet:
		_, err := vm.stack.unwindFrame(i.argInt())
		return err
	default:
		panic("not implemented")
	}
}
