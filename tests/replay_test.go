package tests

import (
	"bytes"
	"testing"

	"github.com/GXQS/SmartContracts/pkg/execution"
	"github.com/GXQS/SmartContracts/pkg/state"
	"github.com/GXQS/SmartContracts/pkg/vm"
)

func TestDeterministicReplay(t *testing.T) {
	db1 := state.NewMemoryDB()
	db2 := state.NewMemoryDB()
	code := []byte{byte(vm.PUSH1), 0x01, byte(vm.PUSH1), 0x02, byte(vm.ADD), byte(vm.STOP)}
	ctx := vm.NewCallContext(vm.Address{}, vm.Address{}, validPayloadFixture(), 100000)
	v1 := vm.New(vm.Config{}, db1)
	v2 := vm.New(vm.Config{}, db2)
	r1, err := execution.DeterministicReplay(v1, []vm.CallContext{ctx}, [][]byte{code})
	if err != nil {
		t.Fatal(err)
	}
	r2, err := execution.DeterministicReplay(v2, []vm.CallContext{ctx}, [][]byte{code})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(r1[0].StateRoot[:], r2[0].StateRoot[:]) {
		t.Fatal("state root mismatch")
	}
}
