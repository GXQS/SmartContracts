package gas

type Schedule struct {
	opcodes map[byte]uint64
}

func DefaultSchedule() Schedule {
	return Schedule{opcodes: map[byte]uint64{
		0x00: 0,
		0x01: 3,
		0x03: 3,
		0x50: 2,
		0x51: 3,
		0x52: 3,
		0x54: 100,
		0x55: 200,
		0x5B: 1,
		0x60: 3,
		0xF3: 0,
		0xFD: 0,
	}}
}

func (s Schedule) OpcodeCost(op byte) uint64 {
	if c, ok := s.opcodes[op]; ok {
		return c
	}
	return 1000
}
