package observability

import "sync/atomic"

var (
	opcodeCount uint64
	gasUsed     uint64
)

func RecordOpcode()      { atomic.AddUint64(&opcodeCount, 1) }
func RecordGas(v uint64) { atomic.AddUint64(&gasUsed, v) }
func SnapshotMetrics() (uint64, uint64) {
	return atomic.LoadUint64(&opcodeCount), atomic.LoadUint64(&gasUsed)
}
