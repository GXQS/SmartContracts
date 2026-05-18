package state

func CompactEncode(nibbles []byte, terminator bool) []byte {
	flags := byte(0)
	if terminator {
		flags |= 0x20
	}
	odd := len(nibbles)%2 == 1
	if odd {
		flags |= 0x10
	}
	encoded := make([]byte, 0, len(nibbles)/2+1)
	if odd {
		encoded = append(encoded, flags|nibbles[0])
		nibbles = nibbles[1:]
	} else {
		encoded = append(encoded, flags)
	}
	for i := 0; i < len(nibbles); i += 2 {
		encoded = append(encoded, (nibbles[i]<<4)|nibbles[i+1])
	}
	return encoded
}
