package vm

import "errors"

type Verifier interface {
	VerifyBytecode([]byte) error
	VerifyABI([]byte) error
}

type DefaultVerifier struct{}

func (DefaultVerifier) VerifyBytecode(code []byte) error {
	for _, b := range code {
		if b == 0xFE {
			return errors.New("invalid opcode 0xFE")
		}
	}
	return nil
}

func (DefaultVerifier) VerifyABI(raw []byte) error {
	if len(raw)%4 != 0 {
		return errors.New("abi payload must be 4-byte aligned")
	}
	return nil
}
