package vm

import "github.com/GXQS/SmartContracts/pkg/gas"

type GasMeter struct {
	remaining uint64
	schedule  gas.Schedule
	used      uint64
	refund    uint64
}

func NewGasMeter(limit uint64, schedule gas.Schedule) *GasMeter {
	return &GasMeter{remaining: limit, schedule: schedule}
}

func (g *GasMeter) Consume(op OpCode, memBefore, memAfter uint64) bool {
	cost := g.schedule.OpcodeCost(byte(op)) + g.schedule.MemoryExpansionCost(memBefore, memAfter)
	if cost > g.remaining {
		return false
	}
	g.remaining -= cost
	g.used += cost
	return true
}

func (g *GasMeter) AddRefund(value uint64) { g.refund += value }
func (g *GasMeter) Used() uint64           { return g.used }
func (g *GasMeter) Remaining() uint64      { return g.remaining }
func (g *GasMeter) Refund() uint64         { return g.refund }
