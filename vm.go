package stackvm

import "fmt"

// VirtualMachine is the main struct that represents the virtual machine.
type VirtualMachine struct {
	stack *stack
}

// New creates a new virtual machine.
func New(opts ...Option) *VirtualMachine {
	var s settings
	opts = append(defaultOpts, opts...)
	for _, opt := range opts {
		opt(&s)
	}
	return &VirtualMachine{
		stack: newStack(s.stackLimit),
	}
}

// Run runs the virtual machine with a given function prototype.
func (vm *VirtualMachine) Run(proto *FuncProto, args ...Value) ([]Value, error) {
	if frame := vm.stack.currentFrame(); frame != nil {
		return nil, fmt.Errorf("%w: VM is already running", ErrIllegalState)
	}
	for _, arg := range args {
		vm.stack.push(arg)
	}
	vm.stack.newFrame(proto)

	for {
		frame := vm.stack.currentFrame()
		if frame == nil {
			// Call stack unwind. Return the values on the stack.
			values := vm.stack.popAll()
			return values, nil
		}
		inst, ok := frame.nextInst()
		if !ok {
			return nil, fmt.Errorf("%w: program ended without return", ErrInvalidProgram)
		}
		frame.incIP()
		if err := inst.execute(vm); err != nil {
			return nil, err
		}
	}
}
