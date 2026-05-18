package crypto

func Zeroize(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}
}
