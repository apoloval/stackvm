package stackvm

import "fmt"

// Value is a value that can be stored in the stack and manipulated by the virtual machine.
type Value struct {
	t typeTag
	v any
}

// NoValue is a value that represents no value.
var NoValue = Value{t: TypeNone}

// NewInt creates a new int value.
func NewInt(v int32) Value {
	return newValue(TypeInt, v)
}

// NewFloat creates a new float value.
func NewFloat(v float32) Value {
	return newValue(TypeFloat, v)
}

// NewBool creates a new bool value.
func NewBool(v bool) Value {
	return newValue(TypeBool, v)
}

// NewString creates a new string value.
func NewString(v string) Value {
	return newValue(TypeString, v)
}

// NewFunction creates a new function value.
func NewFunction(vm *VirtualMachine, proto *FuncProto) Value {
	f := &Function{proto: proto}
	return newValue(TypeFunction, f)
}

// AsInt returns the value as an int.
func (v Value) AsInt() (int32, error) {
	if err := v.ensureType(TypeInt); err != nil {
		return 0, err
	}
	return v.v.(int32), nil
}

// AsFloat returns the value as a float.
func (v Value) AsFloat() (float32, error) {
	if err := v.ensureType(TypeFloat); err != nil {
		return 0, err
	}
	return v.v.(float32), nil
}

// AsBool returns the value as a bool.
func (v Value) AsBool() (bool, error) {
	if err := v.ensureType(TypeBool); err != nil {
		return false, err
	}
	return v.v.(bool), nil
}

// AsString returns the value as a string.
func (v Value) AsString() (string, error) {
	if err := v.ensureType(TypeString); err != nil {
		return "", err
	}
	return v.v.(string), nil
}

func newValue(t typeTag, v any) Value {
	return Value{t: t, v: v}
}

func (v Value) ensureType(t typeTag) error {
	if v.t == t {
		return nil
	}
	return fmt.Errorf("%w: expected %s, got %s", ErrTypeMismatch, typeNames[t], typeNames[v.t])
}

type typeTag uint8

const (
	TypeNone typeTag = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeString
	TypeFunction
)

var typeNames = map[typeTag]string{
	TypeNone:     "none",
	TypeInt:      "int",
	TypeFloat:    "float",
	TypeBool:     "bool",
	TypeString:   "string",
	TypeFunction: "function",
}
