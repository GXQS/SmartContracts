package tests

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
)

func TestRootDeterminism(t *testing.T) {
	db1 := state.NewMemoryDB()
	db2 := state.NewMemoryDB()
	var acc [32]byte
	db1.SetStorage(acc, []byte("a"), []byte("1"))
	db2.SetStorage(acc, []byte("a"), []byte("1"))
	r1 := db1.Root()
	r2 := db2.Root()
	if string(r1) != string(r2) {
		t.Fatal("roots should match")
	}
}
