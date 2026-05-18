package vm

type Receipt struct {
	GasUsed    uint64
	GasRefund  uint64
	Reverted   bool
	ReturnData []byte
	Logs       []LogEntry
	StateRoot  Hash
}

type LogEntry struct {
	Address Address
	Data    []byte
}
