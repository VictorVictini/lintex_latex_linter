package main

type Stack struct {
	elements []Component
}

func (s *Stack) Push(group Component) {
	s.elements = append(s.elements, group)
}

func (s *Stack) Pop() (Component, CreateError) {
	if s.IsEmpty() {
		return nil, SERVER_RESPONSIBLE_STACK_EMPTY
	}

	n := s.Size() - 1
	str := s.elements[n]
	s.elements = s.elements[:n]
	return str, nil
}

func (s *Stack) Peek() (Component, CreateError) {
	if s.IsEmpty() {
		return nil, SERVER_RESPONSIBLE_STACK_EMPTY
	}
	return s.elements[len(s.elements)-1], nil
}

func (s *Stack) Size() int {
	return len(s.elements)
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}
