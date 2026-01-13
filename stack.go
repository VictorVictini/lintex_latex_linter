package main

import "errors"

type Stack struct {
	elements []string
}

func (s *Stack) Push(str string) {
	s.elements = append(s.elements, str)
}

func (s *Stack) Pop() (string, error) {
	if s.IsEmpty() {
		return "", errors.New("Stack is empty, so no element can be removed from it")
	}

	n := s.Size() - 1
	str := s.elements[n]
	s.elements = s.elements[:n]
	return str, nil
}

func (s *Stack) Peek() (string, error) {
	if s.IsEmpty() {
		return "", errors.New("Stack is empty, so it cannot be peeked at")
	}
	return s.elements[len(s.elements)-1], nil
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Stack) Size() int {
	return len(s.elements)
}
