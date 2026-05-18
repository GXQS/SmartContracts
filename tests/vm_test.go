package tests

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
	"github.com/GXQS/SmartContracts/pkg/vm"
)

func TestVMExecutesDeterministicBytecode(t *testing.T) {
	engine := vm.New(vm.Config{}, state.NewMemoryDB())
	ctx := vm.NewCallContext(vm.Address{}, vm.Address{}, nil, 100000)
	code := []byte{byte(vm.PUSH1), 0x02, byte(vm.PUSH1), 0x03, byte(vm.ADD), byte(vm.STOP)}
	receipt, err := engine.Execute(ctx, code)
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if receipt.GasUsed == 0 {
		t.Fatal("expected gas usage")
	}
}
