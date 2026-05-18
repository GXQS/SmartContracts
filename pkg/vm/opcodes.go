package vm

type OpCode byte

const (
	STOP     OpCode = 0x00
	ADD      OpCode = 0x01
	SUB      OpCode = 0x03
	POP      OpCode = 0x50
	MLOAD    OpCode = 0x51
	MSTORE   OpCode = 0x52
	SLOAD    OpCode = 0x54
	SSTORE   OpCode = 0x55
	JUMPDEST OpCode = 0x5B
	PUSH1    OpCode = 0x60
	RETURN   OpCode = 0xF3
	REVERT   OpCode = 0xFD
)
