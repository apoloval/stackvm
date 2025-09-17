package stackvm

import "math"

type opCode uint16

const (
	typNone   opCode = 0x0 // none type
	typInt    opCode = 0x1 // integer type
	typFloat  opCode = 0x2 // float type
	typBool   opCode = 0x3 // boolean type
	typString opCode = 0x4 // string type

	// Control flow instructions
	opNop opCode = 0x0000 // NOP: no operation
	opBr  opCode = 0x0010 // BR: branch (conditional)
	opJmp opCode = 0x0020 // JMP: jump
	opRet opCode = 0x0030 // RET: return from function

	// Stack handling instructions
	opDup   opCode = 0x0100            // DUP: duplicate value
	opPushi opCode = 0x0110 | typInt   // PUSHI: push integer value
	opPushf opCode = 0x0110 | typFloat // PUSHF: push float value
	opPop   opCode = 0x0120            // POP: pop value

	// Arithmetic-logical instructions
	opAddi opCode = 0x0200 | typInt   // ADDI: add integer values
	opAddf opCode = 0x0200 | typFloat // ADDF: add float values
	opSubi opCode = 0x0210 | typInt   // SUBI: subtract integer values
	opSubf opCode = 0x0210 | typFloat // SUBF: subtract float values
	opMuli opCode = 0x0220 | typInt   // MULI: multiply integer values
	opMulf opCode = 0x0220 | typFloat // MULF: multiply float values
	opDivi opCode = 0x0230 | typInt   // DIVI: divide integer values
	opDivf opCode = 0x0230 | typFloat // DIVF: divide float values
	opModi opCode = 0x0240 | typInt   // MODI: modulo integer values
	opNegi opCode = 0x0250 | typInt   // NEGI: negate integer values
	opNegf opCode = 0x0250 | typFloat // NEGF: negate float values

	// Evaluation instructions
	opEqi opCode = 0x0300 | typInt    // EQI: evaluate integers equal to
	opEqf opCode = 0x0300 | typFloat  // EQF: evaluate floats equal to
	opEqb opCode = 0x0300 | typBool   // EQB: evaluate booleans equal to
	opEqs opCode = 0x0300 | typString // EQS: evaluate strings equal to
	opNei opCode = 0x0310 | typInt    // NEI: evaluate integers not equal to
	opNef opCode = 0x0310 | typFloat  // NEF: evaluate floats not equal to
	opNeb opCode = 0x0310 | typBool   // NEB: evaluate booleans not equal to
	opNes opCode = 0x0310 | typString // NES: evaluate strings not equal to

	opGti opCode = 0x0320 | typInt   // GTI: evaluate integers greater than
	opGtf opCode = 0x0320 | typFloat // GTF: evaluate floats greater than
	opGei opCode = 0x0330 | typInt   // GEI: evaluate integers greater than or equal to
	opGef opCode = 0x0330 | typFloat // GEF: evaluate floats greater than or equal to
	opLti opCode = 0x0340 | typInt   // LTI: evaluate integers less than
	opLtf opCode = 0x0340 | typFloat // LTF: evaluate floats less than
	opLei opCode = 0x0350 | typInt   // LEI: evaluate integers less than or equal to
	opLef opCode = 0x0350 | typFloat // LEF: evaluate floats less than or equal to
)

// InstPtr is the pointer to the instruction.
type InstPtr uint32

var nullInstPtr InstPtr = ^InstPtr(0)

// Inst is the instruction of the code segment. It is a 64-bit integer that encodes the
// operation in the first 8 bits and the arguments in the rest.
type Inst uint64

// NOP encodes a NOP instruction.
func NOP() Inst { return makeInst(opNop) }

// BR encodes a BR instruction.
func BR(arg InstPtr) Inst { return makeInst(opBr).withOpInstPtr(arg) }

// JMP encodes a JMP instruction.
func JMP(arg InstPtr) Inst { return makeInst(opJmp).withOpInstPtr(arg) }

// RET encodes a RET instruction.
func RET(arg uint) Inst { return makeInst(opRet).withOpInt(int32(arg)) }

// DUP encodes a DUP instruction.
func DUP(arg int) Inst { return makeInst(opDup).withOpInt(int32(arg)) }

// PUSHI encodes a PUSHI instruction.
func PUSHI(arg int32) Inst { return makeInst(opPushi).withOpInt(arg) }

// PUSHF encodes a PUSHF instruction.
func PUSHF(arg float32) Inst { return makeInst(opPushf).withOpFloat(arg) }

// POP encodes a POP instruction.
func POP(arg int) Inst { return makeInst(opPop).withOpInt(int32(arg)) }

// ADDI encodes an ADDI instruction.
func ADDI() Inst { return makeInst(opAddi) }

// ADDF encodes an ADDF instruction.
func ADDF() Inst { return makeInst(opAddf) }

// SUBI encodes a SUBI instruction.
func SUBI() Inst { return makeInst(opSubi) }

// SUBF encodes a SUBF instruction.
func SUBF() Inst { return makeInst(opSubf) }

// MULI encodes a MULI instruction.
func MULI() Inst { return makeInst(opMuli) }

// DIVI encodes a DIVI instruction.
func DIVI() Inst { return makeInst(opDivi) }

// DIVF encodes a DIVF instruction.
func DIVF() Inst { return makeInst(opDivf) }

// MODI encodes a MODI instruction.
func MODI() Inst { return makeInst(opModi) }

// NEGI encodes a NEGI instruction.
func NEGI() Inst { return makeInst(opNegi) }

// NEGF encodes a NEGF instruction.
func NEGF() Inst { return makeInst(opNegf) }

// EQI encodes a EQI instruction.
func EQI() Inst { return makeInst(opEqi) }

// EQF encodes a EQF instruction.
func EQF() Inst { return makeInst(opEqf) }

// EQB encodes a EQB instruction.
func EQB() Inst { return makeInst(opEqb) }

// EQS encodes a EQS instruction.
func EQS() Inst { return makeInst(opEqs) }

// NEI encodes a NEI instruction.
func NEI() Inst { return makeInst(opNei) }

// NEF encodes a NEF instruction.
func NEF() Inst { return makeInst(opNef) }

// NEB encodes a NEB instruction.
func NEB() Inst { return makeInst(opNeb) }

// NES encodes a NES instruction.
func NES() Inst { return makeInst(opNes) }

// GTI encodes a GTI instruction.
func GTI() Inst { return makeInst(opGti) }

// GTF encodes a GTF instruction.
func GTF() Inst { return makeInst(opGtf) }

// GEI encodes a GEI instruction.
func GEI() Inst { return makeInst(opGei) }

// GEF encodes a GEF instruction.
func GEF() Inst { return makeInst(opGef) }

// LTI encodes a LTI instruction.
func LTI() Inst { return makeInst(opLti) }

// LTF encodes a LTF instruction.
func LTF() Inst { return makeInst(opLtf) }

// LEI encodes a LEI instruction.
func LEI() Inst { return makeInst(opLei) }

// LEF encodes a LEF instruction.
func LEF() Inst { return makeInst(opLef) }

func makeInst(op opCode) Inst {
	return Inst(op) << 48
}

func (i Inst) withOpInt(arg int32) Inst {
	return i | Inst(arg)
}

func (i Inst) withOpFloat(arg float32) Inst {
	return i | Inst(math.Float32bits(arg))
}

func (i Inst) withOpInstPtr(to InstPtr) Inst {
	return i | Inst(to)
}

func (i Inst) opCode() opCode {
	return opCode(i >> 48)
}

func (i Inst) argInt() int32 {
	return int32(i & 0xFFFF_FFFF)
}

func (i Inst) argFloat() float32 {
	return math.Float32frombits(uint32(i & 0xFFFF_FFFF))
}

func (i Inst) execute(vm *VirtualMachine) error {
	switch i.opCode() {
	case opNop:
		return nil
	case opBr:
		return withBoolSingle(vm, func(cond bool) error {
			if cond {
				vm.stack.currentFrame().ip = InstPtr(i.argInt())
			}
			return nil
		})
	case opJmp:
		vm.stack.currentFrame().ip = InstPtr(i.argInt())
		return nil
	case opRet:
		_, err := vm.stack.unwindFrame(int(i.argInt()))
		return err
	case opDup:
		v, err := vm.stack.peek(int(i.argInt()))
		if err != nil {
			return err
		}
		return vm.stack.push(v)
	case opPushi:
		return vm.stack.push(NewInt(i.argInt()))
	case opPushf:
		return vm.stack.push(NewFloat(i.argFloat()))
	case opPop:
		v, err := vm.stack.pop()
		if err != nil {
			return err
		}
		return vm.stack.poke(int(i.argInt()), v)
	case opAddi:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewInt(a + b))
		})
	case opAddf:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewFloat(a + b))
		})
	case opSubi:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewInt(a - b))
		})
	case opSubf:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewFloat(a - b))
		})
	case opMuli:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewInt(a * b))
		})
	case opMulf:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewFloat(a * b))
		})
	case opDivi:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewInt(a / b))
		})
	case opDivf:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewFloat(a / b))
		})
	case opModi:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewInt(a % b))
		})
	case opNegi:
		return withIntSingle(vm, func(a int32) error {
			return vm.stack.push(NewInt(-a))
		})
	case opNegf:
		return withFloatSingle(vm, func(a float32) error {
			return vm.stack.push(NewFloat(-a))
		})
	case opEqi:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewBool(a == b))
		})
	case opEqf:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewBool(a == b))
		})
	case opEqb:
		return withBoolTuple(vm, func(a, b bool) error {
			return vm.stack.push(NewBool(a == b))
		})
	case opEqs:
		return withStringTuple(vm, func(a, b string) error {
			return vm.stack.push(NewBool(a == b))
		})
	case opNei:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewBool(a != b))
		})
	case opNef:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewBool(a != b))
		})
	case opNeb:
		return withBoolTuple(vm, func(a, b bool) error {
			return vm.stack.push(NewBool(a != b))
		})
	case opNes:
		return withStringTuple(vm, func(a, b string) error {
			return vm.stack.push(NewBool(a != b))
		})
	case opGti:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewBool(a > b))
		})
	case opGtf:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewBool(a > b))
		})
	case opGei:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewBool(a >= b))
		})
	case opGef:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewBool(a >= b))
		})
	case opLti:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewBool(a < b))
		})
	case opLtf:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewBool(a < b))
		})
	case opLei:
		return withIntTuple(vm, func(a, b int32) error {
			return vm.stack.push(NewBool(a <= b))
		})
	case opLef:
		return withFloatTuple(vm, func(a, b float32) error {
			return vm.stack.push(NewBool(a <= b))
		})
	default:
		panic("not implemented")
	}
}

func withBoolSingle(vm *VirtualMachine, f func(a bool) error) error {
	a, err := vm.stack.popBool()
	if err != nil {
		return err
	}
	return f(a)
}

func withIntSingle(vm *VirtualMachine, f func(a int32) error) error {
	a, err := vm.stack.popInt()
	if err != nil {
		return err
	}
	return f(a)
}

func withFloatSingle(vm *VirtualMachine, f func(a float32) error) error {
	a, err := vm.stack.popFloat()
	if err != nil {
		return err
	}
	return f(a)
}

func withIntTuple(vm *VirtualMachine, f func(a, b int32) error) error {
	b, err := vm.stack.popInt()
	if err != nil {
		return err
	}
	a, err := vm.stack.popInt()
	if err != nil {
		return err
	}
	return f(a, b)
}

func withFloatTuple(vm *VirtualMachine, f func(a, b float32) error) error {
	b, err := vm.stack.popFloat()
	if err != nil {
		return err
	}
	a, err := vm.stack.popFloat()
	if err != nil {
		return err
	}
	return f(a, b)
}

func withBoolTuple(vm *VirtualMachine, f func(a, b bool) error) error {
	b, err := vm.stack.popBool()
	if err != nil {
		return err
	}
	a, err := vm.stack.popBool()
	if err != nil {
		return err
	}
	return f(a, b)
}

func withStringTuple(vm *VirtualMachine, f func(a, b string) error) error {
	b, err := vm.stack.popString()
	if err != nil {
		return err
	}
	a, err := vm.stack.popString()
	if err != nil {
		return err
	}
	return f(a, b)
}
