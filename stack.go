package main

import "errors"

type Stack struct {
	elements []Component
}

func (s *Stack) Push(group Component) {
	s.elements = append(s.elements, group)
}

func (s *Stack) Pop() (Component, error) {
	if s.IsEmpty() {
		return nil, errors.New("Stack is empty, so no element can be removed from it")
	}

	n := s.Size() - 1
	str := s.elements[n]
	s.elements = s.elements[:n]
	return str, nil
}

func (s *Stack) Peek() (Component, error) {
	if s.IsEmpty() {
		return nil, errors.New("Stack is empty, so it cannot be peeked at")
	}
	return s.elements[len(s.elements)-1], nil
}

func (s *Stack) Size() int {
	return len(s.elements)
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}
