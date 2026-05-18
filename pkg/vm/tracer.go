package vm

type TraceEvent struct {
	PC        uint64
	Op        OpCode
	GasBefore uint64
	GasAfter  uint64
}

type Tracer interface {
	OnStep(event TraceEvent)
}

type NoopTracer struct{}

func (NoopTracer) OnStep(TraceEvent) {}
