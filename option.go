package stackvm

type settings struct {
	stackLimit int
}

// Option is a function that configures the virtual machine.
type Option func(*settings)

// WithStackLimit sets the limit of the stack.
func WithStackLimit(limit int) Option {
	return func(vm *settings) {
		vm.stackLimit = limit
	}
}

var defaultOpts = []Option{
	WithStackLimit(256),
}
