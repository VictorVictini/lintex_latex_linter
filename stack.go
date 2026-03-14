package main

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Push(group T) {
	s.elements = append(s.elements, group)
}

func (s *Stack[T]) Pop() (res T, _ CreateError) {
	if s.IsEmpty() {
		return res, SERVER_RESPONSIBLE_STACK_EMPTY
	}

	n := s.Size() - 1
	res = s.elements[n]
	s.elements = s.elements[:n]
	return res, nil
}

func (s *Stack[T]) Peek() (res T, _ CreateError) {
	if s.IsEmpty() {
		return res, SERVER_RESPONSIBLE_STACK_EMPTY
	}
	return s.elements[len(s.elements)-1], nil
}

func (s *Stack[T]) Size() int {
	return len(s.elements)
}

func (s *Stack[T]) IsEmpty() bool {
	return s.Size() == 0
}
