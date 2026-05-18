package tests

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/vm"
)

func TestMemoryBounds(t *testing.T) {
	m := vm.NewMemory(8)
	if err := m.Store(4, []byte{1, 2, 3, 4, 5}); err == nil {
		t.Fatal("expected memory limit exceeded")
	}
}
