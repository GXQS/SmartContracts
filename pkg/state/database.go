package state

import (
	"crypto/sha256"
	"sort"
)

type SnapshotID uint64

type Database interface {
	SetStorage(account [32]byte, key, value []byte)
	GetStorage(account [32]byte, key []byte) []byte
	Root() []byte
	Snapshot() SnapshotID
	RevertToSnapshot(SnapshotID)
}

type MemoryDB struct {
	storage   map[[32]byte]map[string][]byte
	history   []snapshotState
	snapshotN SnapshotID
}

type snapshotState struct {
	id      SnapshotID
	storage map[[32]byte]map[string][]byte
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{storage: map[[32]byte]map[string][]byte{}}
}

func cloneStorage(in map[[32]byte]map[string][]byte) map[[32]byte]map[string][]byte {
	out := make(map[[32]byte]map[string][]byte, len(in))
	for addr, slots := range in {
		copySlots := make(map[string][]byte, len(slots))
		for k, v := range slots {
			buf := make([]byte, len(v))
			copy(buf, v)
			copySlots[k] = buf
		}
		out[addr] = copySlots
	}
	return out
}

func (db *MemoryDB) Snapshot() SnapshotID {
	db.snapshotN++
	db.history = append(db.history, snapshotState{id: db.snapshotN, storage: cloneStorage(db.storage)})
	return db.snapshotN
}

func (db *MemoryDB) RevertToSnapshot(id SnapshotID) {
	for i := len(db.history) - 1; i >= 0; i-- {
		if db.history[i].id == id {
			db.storage = cloneStorage(db.history[i].storage)
			db.history = db.history[:i]
			return
		}
	}
}

func (db *MemoryDB) SetStorage(account [32]byte, key, value []byte) {
	slots, ok := db.storage[account]
	if !ok {
		slots = map[string][]byte{}
		db.storage[account] = slots
	}
	buf := make([]byte, len(value))
	copy(buf, value)
	slots[string(key)] = buf
}

func (db *MemoryDB) GetStorage(account [32]byte, key []byte) []byte {
	slots, ok := db.storage[account]
	if !ok {
		return nil
	}
	v := slots[string(key)]
	out := make([]byte, len(v))
	copy(out, v)
	return out
}

func (db *MemoryDB) Root() []byte {
	t := NewTrie()
	accounts := make([][32]byte, 0, len(db.storage))
	for a := range db.storage {
		accounts = append(accounts, a)
	}
	sort.Slice(accounts, func(i, j int) bool { return string(accounts[i][:]) < string(accounts[j][:]) })
	for _, addr := range accounts {
		slots := db.storage[addr]
		keys := make([]string, 0, len(slots))
		for k := range slots {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			h := sha256.Sum256(append(addr[:], []byte(k)...))
			t.Update(h[:], slots[k])
		}
	}
	return t.RootHash()
}
