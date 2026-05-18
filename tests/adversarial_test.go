package tests

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
	"github.com/GXQS/SmartContracts/pkg/vm"
)

func TestInvalidOpcodeRejected(t *testing.T) {
	engine := vm.New(vm.Config{}, state.NewMemoryDB())
	ctx := vm.NewCallContext(vm.Address{}, vm.Address{}, nil, 100000)
	_, err := engine.Execute(ctx, []byte{0xFE})
	if err == nil {
		t.Fatal("expected invalid opcode error")
	}
}
