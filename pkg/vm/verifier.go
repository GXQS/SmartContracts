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

const (
	PayloadHeaderSize     = 108
	MLDSA65SignatureSize  = 3309
	MLDSA65PublicKeySize  = 1952
	HeaderVersionOffset   = 0
	HeaderNonceOffset     = 4
	HeaderSenderOffset    = 12
	HeaderRecipientOffset = 44
	HeaderGasLimitOffset  = 76
	HeaderGasPriceOffset  = 84
	HeaderSigOffsetOffset = 92
	HeaderPubOffsetOffset = 96
	HeaderDataOffset      = 100
	HeaderDataLenOffset   = 104
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

func WriteU32BE(raw []byte, offset int, value uint32) {
	raw[offset] = byte(value >> 24)
	raw[offset+1] = byte(value >> 16)
	raw[offset+2] = byte(value >> 8)
	raw[offset+3] = byte(value)
}

func WriteU64BE(raw []byte, offset int, value uint64) {
	raw[offset] = byte(value >> 56)
	raw[offset+1] = byte(value >> 48)
	raw[offset+2] = byte(value >> 40)
	raw[offset+3] = byte(value >> 32)
	raw[offset+4] = byte(value >> 24)
	raw[offset+5] = byte(value >> 16)
	raw[offset+6] = byte(value >> 8)
	raw[offset+7] = byte(value)
}

func BuildBoundaryPayload(dataLen uint32) []byte {
	sigOffset := PayloadHeaderSize
	pubOffset := sigOffset + MLDSA65SignatureSize
	dataOffset := pubOffset + MLDSA65PublicKeySize

	payload := make([]byte, dataOffset+int(dataLen))
	WriteU32BE(payload, HeaderVersionOffset, 1)
	WriteU64BE(payload, HeaderNonceOffset, 1)
	WriteU64BE(payload, HeaderGasLimitOffset, 1_000_000)
	WriteU64BE(payload, HeaderGasPriceOffset, 1)
	WriteU32BE(payload, HeaderSigOffsetOffset, uint32(sigOffset))
	WriteU32BE(payload, HeaderPubOffsetOffset, uint32(pubOffset))
	WriteU32BE(payload, HeaderDataOffset, uint32(dataOffset))
	WriteU32BE(payload, HeaderDataLenOffset, dataLen)
	return payload
}

// VerifyPayloadBoundaries inspects the memory map without expensive heap allocations.
func (v *ZeroCopyPayloadVerifier) VerifyPayloadBoundaries() (*TransactionHeader, error) {
	// Minimum header size safety check: Version(4) + Nonce(8) + Sender(32) + Recipient(32) + GasLimit(8) + GasPrice(8) + SigOffset(4) + PubKeyOffset(4) + DataOffset(4) + DataLen(4) = PayloadHeaderSize.
	if len(v.rawPayload) < PayloadHeaderSize {
		return nil, ErrMalformedPayloadHeader
	}

	// Parse header layout boundaries directly from big-endian data stream windows.
	header := &TransactionHeader{
		Version:      readU32BE(v.rawPayload, HeaderVersionOffset),
		Nonce:        readU64BE(v.rawPayload, HeaderNonceOffset),
		GasLimit:     readU64BE(v.rawPayload, HeaderGasLimitOffset),
		GasPrice:     readU64BE(v.rawPayload, HeaderGasPriceOffset),
		SigOffset:    readU32BE(v.rawPayload, HeaderSigOffsetOffset),
		PubKeyOffset: readU32BE(v.rawPayload, HeaderPubOffsetOffset),
		DataOffset:   readU32BE(v.rawPayload, HeaderDataOffset),
		DataLen:      readU32BE(v.rawPayload, HeaderDataLenOffset),
	}

	// Extract explicit Address segments using deterministic memory windows.
	copy(header.Sender[:], v.rawPayload[HeaderSenderOffset:HeaderRecipientOffset])
	copy(header.Recipient[:], v.rawPayload[HeaderRecipientOffset:HeaderGasLimitOffset])

	payloadLen := uint32(len(v.rawPayload))
	if header.SigOffset < PayloadHeaderSize || header.PubKeyOffset < PayloadHeaderSize || header.DataOffset < PayloadHeaderSize {
		return nil, ErrMalformedPayloadHeader
	}

	// Enforce hard structural constraints derived from FIPS 204/203 standards.
	if header.SigOffset > payloadLen {
		return nil, ErrSignatureSizeMismatch
	}
	if payloadLen-header.SigOffset < MLDSA65SignatureSize {
		return nil, ErrSignatureSizeMismatch
	}
	if header.PubKeyOffset > payloadLen {
		return nil, ErrPublicKeySizeMismatch
	}
	if payloadLen-header.PubKeyOffset < MLDSA65PublicKeySize {
		return nil, ErrPublicKeySizeMismatch
	}
	if header.DataOffset > payloadLen {
		return nil, ErrMalformedPayloadHeader
	}
	if payloadLen-header.DataOffset < header.DataLen {
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
		return nil, ErrMalformedPayloadHeader
	}
	return NewZeroCopyPayloadVerifier(raw).VerifyPayloadBoundaries()
}
