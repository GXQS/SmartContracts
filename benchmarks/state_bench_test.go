package benchmarks

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
)

func BenchmarkStateRoot(b *testing.B) {
	db := state.NewMemoryDB()
	var a [32]byte
	db.SetStorage(a, []byte("k"), []byte("v"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = db.Root()
	}
}
