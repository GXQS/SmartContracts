package crypto

import "errors"

type PrecompileID byte

const (
	PrecompileMLDSAVerify PrecompileID = 0x01
	PrecompileMLKEMEncap  PrecompileID = 0x02
	PrecompileMLKEMDecap  PrecompileID = 0x03
)

func ExecutePrecompile(id PrecompileID, payload [][]byte) ([][]byte, error) {
	switch id {
	case PrecompileMLDSAVerify:
		if len(payload) != 3 {
			return nil, errors.New("mldsa verify requires pubkey, message, signature")
		}
		if err := VerifyMLDSA(payload[0], payload[1], payload[2]); err != nil {
			return nil, err
		}
		return [][]byte{[]byte{1}}, nil
	case PrecompileMLKEMEncap:
		if len(payload) != 1 {
			return nil, errors.New("mlkem encapsulate requires public key")
		}
		ct, ss, err := EncapsulateMLKEM(payload[0])
		if err != nil {
			return nil, err
		}
		return [][]byte{ct, ss}, nil
	case PrecompileMLKEMDecap:
		if len(payload) != 2 {
			return nil, errors.New("mlkem decapsulate requires private key and ciphertext")
		}
		ss, err := DecapsulateMLKEM(payload[0], payload[1])
		if err != nil {
			return nil, err
		}
		return [][]byte{ss}, nil
	default:
		return nil, errors.New("unknown precompile")
	}
}
