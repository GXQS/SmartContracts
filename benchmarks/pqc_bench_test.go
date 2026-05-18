package benchmarks

import (
	"testing"

	pq "github.com/GXQS/SmartContracts/pkg/crypto"
)

func BenchmarkEntropyValidation(b *testing.B) {
	v := make([]byte, 32)
	v[0] = 1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pq.ValidateEntropy(v)
	}
}
