package crypto

import "errors"

func ValidateEntropy(entropy []byte) error {
	if len(entropy) < 32 {
		return errors.New("entropy too short")
	}
	var nonZero byte
	for _, b := range entropy {
		nonZero |= b
	}
	if nonZero == 0 {
		return errors.New("entropy is all zeros")
	}
	return nil
}
