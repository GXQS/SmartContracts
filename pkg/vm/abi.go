package vm

import (
	"encoding/binary"
	"errors"
)

type ABICall struct {
	Selector [4]byte
	Args     []byte
}

func DecodeABICall(input []byte) (ABICall, error) {
	if len(input) < 4 {
		return ABICall{}, errors.New("abi input too short")
	}
	var sel [4]byte
	copy(sel[:], input[:4])
	args := make([]byte, len(input)-4)
	copy(args, input[4:])
	return ABICall{Selector: sel, Args: args}, nil
}

func EncodeUint64(v uint64) []byte {
	out := make([]byte, 32)
	binary.BigEndian.PutUint64(out[24:], v)
	return out
}
