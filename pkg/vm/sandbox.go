package vm

import "errors"

type Sandbox struct {
	allowHostCalls bool
}

func NewSandbox() Sandbox { return Sandbox{} }

func (s Sandbox) CheckHostAccess() error {
	if !s.allowHostCalls {
		return errors.New("host interface denied in sandbox")
	}
	return nil
}
