package benchmarks

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
	"github.com/GXQS/SmartContracts/pkg/vm"
)

func BenchmarkVMExecution(b *testing.B) {
	engine := vm.New(vm.Config{}, state.NewMemoryDB())
	ctx := vm.NewCallContext(vm.Address{}, vm.Address{}, nil, 100000)
	code := []byte{byte(vm.PUSH1), 0x02, byte(vm.PUSH1), 0x03, byte(vm.ADD), byte(vm.STOP)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Execute(ctx, code)
	}
}
