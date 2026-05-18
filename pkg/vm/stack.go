package vm

import "errors"

type Stack struct {
	data []uint64
	max  int
}

func NewStack(max int) *Stack {
	if max <= 0 {
		max = 1024
	}
	return &Stack{max: max}
}

func (s *Stack) Push(v uint64) error {
	if len(s.data) >= s.max {
		return errors.New("stack overflow")
	}
	s.data = append(s.data, v)
	return nil
}

func (s *Stack) Pop() (uint64, error) {
	if len(s.data) == 0 {
		return 0, errors.New("stack underflow")
	}
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v, nil
}

func (s *Stack) Len() int { return len(s.data) }
