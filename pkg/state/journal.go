package state

type JournalEntry struct {
	Account [32]byte
	Key     []byte
	Before  []byte
	After   []byte
}
