package tests

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/gas"
)

func TestMemoryExpansionIsQuadratic(t *testing.T) {
	s := gas.DefaultSchedule()
	c1 := s.MemoryExpansionCost(0, 32)
	c2 := s.MemoryExpansionCost(32, 64)
	if c2 <= c1 {
		t.Fatal("expected later expansion to be more expensive")
	}
}
