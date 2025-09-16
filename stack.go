package stackvm

type stack struct {
	data   []Value
	frames []frame
	limit  int
}

func newStack(limit int) *stack {
	s := &stack{
		data:  make([]Value, 0, limit),
		limit: limit,
	}
	return s
}

func (s *stack) newFrame(proto *FuncProto) *frame {
	s.frames = append(s.frames, frame{
		proto:     proto,
		stackBase: len(s.data) - proto.nargs,
		ip:        0,
	})
	return &s.frames[len(s.frames)-1]
}

func (s *stack) unwindFrame(nres int) (f frame, err error) {
	if frame := s.currentFrame(); frame == nil {
		err = ErrStackUnderflow
		return
	} else {
		f = *frame
	}
	s.frames = s.frames[:len(s.frames)-1]

	// displace the last nres items from the stack to the base
	for i := 0; i < nres; i++ {
		s.data[f.stackBase+i] = s.data[len(s.data)-1-i]
	}
	s.data = s.data[:f.stackBase+nres]

	return
}

func (s *stack) currentFrame() *frame {
	if len(s.frames) == 0 {
		return nil
	}
	return &s.frames[len(s.frames)-1]
}

func (s *stack) push(item Value) error {
	if len(s.data) >= s.limit {
		return ErrStackOverflow
	}
	s.data = append(s.data, item)
	return nil
}

func (s *stack) pop() (Value, error) {
	frame := s.currentFrame()
	if len(s.data) == 0 || (frame != nil && len(s.data) <= frame.stackBase) {
		return NoValue, ErrStackUnderflow
	}
	item := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return item, nil
}

func (s *stack) popInt() (int, error) {
	item, err := s.pop()
	if err != nil {
		return 0, err
	}
	return item.AsInt()
}

func (s *stack) popAll() []Value {
	var base int
	if frame := s.currentFrame(); frame != nil {
		base = frame.stackBase
	}
	values := make([]Value, len(s.data)-base)
	copy(values, s.data[base:])
	s.data = s.data[:base]
	return values
}

type frame struct {
	proto     *FuncProto
	stackBase int
	ip        InstPtr
}

func (f *frame) nextInst() (Inst, bool) {
	if f.ip >= InstPtr(len(f.proto.bytecode)) {
		return 0, false
	}
	return f.proto.bytecode[f.ip], true
}

func (f *frame) incIP() {
	f.ip++
}
