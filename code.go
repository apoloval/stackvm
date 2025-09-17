package stackvm

type opCode uint8

const (
	opNop opCode = iota
	opAddi
	opBr
	opDup
	opEviGt
	opJmp
	opPushi
	opRet
)

var opCodeToStr = map[opCode]string{
	opNop:   "NOP",
	opAddi:  "ADDI",
	opBr:    "BR",
	opDup:   "DUP",
	opEviGt: "EVIGT",
	opJmp:   "JUMP",
	opPushi: "PUSHI",
	opRet:   "RET",
}

// InstPtr is the pointer to the instruction.
type InstPtr uint32

var nullInstPtr InstPtr = ^InstPtr(0)

// Inst is the instruction of the code segment. It is a 64-bit integer that encodes the
// operation in the first 8 bits and the arguments in the rest.
type Inst uint64

// ADDI encodes an ADDI instruction.
func ADDI() Inst {
	return makeInst(opAddi)
}

// BR encodes a BR instruction.
func BR(arg InstPtr) Inst { return makeInst(opBr).withOpInstPtr(arg) }

// DUP encodes a DUP instruction.
func DUP(arg int) Inst { return makeInst(opDup).withOpInt(int32(arg)) }

// EVIGT encodes a EVIGT instruction.
func EVIGT() Inst { return makeInst(opEviGt) }

// JMP encodes a JMP instruction.
func JMP(arg InstPtr) Inst { return makeInst(opJmp).withOpInstPtr(arg) }

// PUSHI encodes a PUSHI instruction.
func PUSHI(arg int32) Inst { return makeInst(opPushi).withOpInt(arg) }

// RET encodes a RET instruction.
func RET(nres uint) Inst { return makeInst(opRet).withOpInt(int32(nres)) }

func makeInst(op opCode) Inst {
	return Inst(op) << 56
}

func (i Inst) withOpInt(arg int32) Inst {
	return i | Inst(arg)
}

func (i Inst) withOpInstPtr(to InstPtr) Inst {
	return i | Inst(to)
}

func (i Inst) opCode() opCode {
	return opCode(i >> 56)
}

func (i Inst) argInt() int32 {
	return int32(i & 0xFFFF_FFFF)
}

func (i Inst) execute(vm *VirtualMachine) error {
	switch i.opCode() {
	case opNop:
		return nil
	case opAddi:
		b, err := vm.stack.popInt()
		if err != nil {
			return err
		}
		a, err := vm.stack.popInt()
		if err != nil {
			return err
		}
		return vm.stack.push(NewInt(a + b))
	case opBr:
		cond, err := vm.stack.popBool()
		if err != nil {
			return err
		}
		if cond {
			vm.stack.currentFrame().ip = InstPtr(i.argInt())
		}
		return nil
	case opDup:
		return vm.stack.dup(int(i.argInt()))
	case opEviGt:
		b, err := vm.stack.popInt()
		if err != nil {
			return err
		}
		a, err := vm.stack.popInt()
		if err != nil {
			return err
		}
		return vm.stack.push(NewBool(a > b))
	case opJmp:
		vm.stack.currentFrame().ip = InstPtr(i.argInt())
		return nil
	case opPushi:
		return vm.stack.push(NewInt(i.argInt()))
	case opRet:
		_, err := vm.stack.unwindFrame(int(i.argInt()))
		return err
	default:
		panic("not implemented")
	}
}
