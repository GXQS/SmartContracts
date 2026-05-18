package vm

import "errors"

type CallContext struct {
	Caller Address
	Callee Address
	Input  []byte
	Gas    uint64
	Value  uint64
	Depth  int
}

func NewCallContext(caller, callee Address, input []byte, gas uint64) CallContext {
	return CallContext{Caller: caller, Callee: callee, Input: append([]byte(nil), input...), Gas: gas}
}

func (c CallContext) Validate(limits RuntimeLimits) error {
	if c.Depth > limits.MaxCallDepth {
		return errors.New("max call depth exceeded")
	}
	return nil
}
