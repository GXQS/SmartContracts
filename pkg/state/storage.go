package state

func StorageKey(account [32]byte, key []byte) []byte {
	out := make([]byte, 32+len(key))
	copy(out, account[:])
	copy(out[32:], key)
	return out
}
