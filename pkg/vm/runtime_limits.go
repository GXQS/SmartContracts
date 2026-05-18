package vm

type RuntimeLimits struct {
	MaxCallDepth int
	MaxMemory    uint64
	MaxSteps     uint64
}
