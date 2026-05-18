package tests

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
	"github.com/GXQS/SmartContracts/pkg/vm"
)

func FuzzInterpreterNeverPanics(f *testing.F) {
	f.Add([]byte{byte(vm.STOP)})
	f.Fuzz(func(t *testing.T, code []byte) {
		engine := vm.New(vm.Config{}, state.NewMemoryDB())
		ctx := vm.NewCallContext(vm.Address{}, vm.Address{}, nil, 100000)
		_, _ = engine.Execute(ctx, code)
	})
}
