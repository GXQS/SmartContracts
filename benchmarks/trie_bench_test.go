package benchmarks

import (
	"strconv"
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
)

func BenchmarkTrieInsert(b *testing.B) {
	tr := state.NewTrie()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strconv.Itoa(i)
		tr.Update([]byte(s), []byte(s))
	}
}
