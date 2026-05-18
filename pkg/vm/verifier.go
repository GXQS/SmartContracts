package vm

import (
	"errors"

	"github.com/GXQS/SmartContracts/pkg/state"
)

var (
	ErrMalformedPayloadHeader = errors.New("boundary violation: invalid transaction payload header")
	ErrSignatureSizeMismatch  = errors.New("boundary violation: ML-DSA-65 signature must be exactly 3309 bytes")
	ErrPublicKeySizeMismatch  = errors.New("boundary violation: ML-DSA-65 public key must be exactly 1952 bytes")
)

// TransactionHeader mirrors the exact Protobuf schema defined in GXQS/Manifest.
type TransactionHeader struct {
	Version      uint32
	Nonce        uint64
	Sender       state.Address
	Recipient    state.Address
	GasLimit     uint64
	GasPrice     uint64
	SigOffset    uint32
	PubKeyOffset uint32
	DataOffset   uint32
	DataLen      uint32
}

// ZeroCopyPayloadVerifier processes incoming byte streams using pointer-arithmetic windows.
type ZeroCopyPayloadVerifier struct {
	rawPayload []byte
}

func NewZeroCopyPayloadVerifier(payload []byte) *ZeroCopyPayloadVerifier {
	return &ZeroCopyPayloadVerifier{rawPayload: payload}
}

func readU32BE(raw []byte, offset int) uint32 {
	return uint32(raw[offset])<<24 |
		uint32(raw[offset+1])<<16 |
		uint32(raw[offset+2])<<8 |
		uint32(raw[offset+3])
}

func readU64BE(raw []byte, offset int) uint64 {
	return uint64(raw[offset])<<56 |
		uint64(raw[offset+1])<<48 |
		uint64(raw[offset+2])<<40 |
		uint64(raw[offset+3])<<32 |
		uint64(raw[offset+4])<<24 |
		uint64(raw[offset+5])<<16 |
		uint64(raw[offset+6])<<8 |
		uint64(raw[offset+7])
}

// VerifyPayloadBoundaries inspects the memory map without expensive heap allocations.
func (v *ZeroCopyPayloadVerifier) VerifyPayloadBoundaries() (*TransactionHeader, error) {
	// Minimum header size safety check (4 + 8 + 32 + 32 + 8 + 8 + 4 + 4 + 4 + 4 = 108 bytes).
	if len(v.rawPayload) < 108 {
		return nil, ErrMalformedPayloadHeader
	}

	// Parse header layout boundaries directly from big-endian data stream windows.
	header := &TransactionHeader{
		Version:      readU32BE(v.rawPayload, 0),
		Nonce:        readU64BE(v.rawPayload, 4),
		GasLimit:     readU64BE(v.rawPayload, 76),
		GasPrice:     readU64BE(v.rawPayload, 84),
		SigOffset:    readU32BE(v.rawPayload, 92),
		PubKeyOffset: readU32BE(v.rawPayload, 96),
		DataOffset:   readU32BE(v.rawPayload, 100),
		DataLen:      readU32BE(v.rawPayload, 104),
	}

	// Extract explicit Address segments using deterministic memory windows.
	copy(header.Sender[:], v.rawPayload[12:44])
	copy(header.Recipient[:], v.rawPayload[44:76])

	payloadLen := uint32(len(v.rawPayload))
	if header.SigOffset < 108 || header.PubKeyOffset < 108 || header.DataOffset < 108 {
		return nil, ErrMalformedPayloadHeader
	}

	// Enforce hard structural constraints derived from FIPS 204/203 standards.
	if uint64(header.SigOffset)+3309 > uint64(payloadLen) {
		return nil, ErrSignatureSizeMismatch
	}
	if uint64(header.PubKeyOffset)+1952 > uint64(payloadLen) {
		return nil, ErrPublicKeySizeMismatch
	}
	if uint64(header.DataOffset)+uint64(header.DataLen) > uint64(payloadLen) {
		return nil, ErrMalformedPayloadHeader
	}

	return header, nil
}

type Verifier interface {
	VerifyBytecode([]byte) error
	VerifyABI([]byte) error
	VerifyPayload([]byte) (*TransactionHeader, error)
}

type DefaultVerifier struct{}

func (DefaultVerifier) VerifyBytecode(code []byte) error {
	for _, b := range code {
		if b == 0xFE {
			return errors.New("invalid opcode 0xFE")
		}
	}
	return nil
}

func (DefaultVerifier) VerifyABI(raw []byte) error {
	if len(raw)%4 != 0 {
		return errors.New("abi payload must be 4-byte aligned")
	}
	return nil
}

func (DefaultVerifier) VerifyPayload(raw []byte) (*TransactionHeader, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	return NewZeroCopyPayloadVerifier(raw).VerifyPayloadBoundaries()
}
