package tests

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
	"github.com/GXQS/SmartContracts/pkg/vm"
)

func TestIntegrationStorageRoundTrip(t *testing.T) {
	engine := vm.New(vm.Config{}, state.NewMemoryDB())
	ctx := vm.NewCallContext(vm.Address{}, vm.Address{}, nil, 100000)
	code := []byte{byte(vm.PUSH1), 0x01, byte(vm.PUSH1), 0x2A, byte(vm.SSTORE), byte(vm.PUSH1), 0x01, byte(vm.SLOAD), byte(vm.STOP)}
	_, err := engine.Execute(ctx, code)
	if err != nil {
		t.Fatalf("integration execute failed: %v", err)
	}
}
