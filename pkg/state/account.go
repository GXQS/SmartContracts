package state

type Account struct {
	Nonce    uint64
	Balance  uint64
	CodeHash [32]byte
}
