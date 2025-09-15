package stackvm

type stack struct {
	items []any
	limit int
}

func newStack(limit int) *stack {
	return &stack{
		items: make([]any, 0, limit),
		limit: limit,
	}
}

func (s *stack) push(item any) error {
	if len(s.items) >= s.limit {
		return ErrStackOverflow
	}
	s.items = append(s.items, item)
	return nil
}

func (s *stack) pop() (any, error) {
	if len(s.items) == 0 {
		return nil, ErrStackUnderflow
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}
