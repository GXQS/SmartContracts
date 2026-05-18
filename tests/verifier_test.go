package tests

import (
	"errors"
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
	"github.com/GXQS/SmartContracts/pkg/vm"
)

func putU32BE(dst []byte, offset int, value uint32) {
	dst[offset] = byte(value >> 24)
	dst[offset+1] = byte(value >> 16)
	dst[offset+2] = byte(value >> 8)
	dst[offset+3] = byte(value)
}

func putU64BE(dst []byte, offset int, value uint64) {
	dst[offset] = byte(value >> 56)
	dst[offset+1] = byte(value >> 48)
	dst[offset+2] = byte(value >> 40)
	dst[offset+3] = byte(value >> 32)
	dst[offset+4] = byte(value >> 24)
	dst[offset+5] = byte(value >> 16)
	dst[offset+6] = byte(value >> 8)
	dst[offset+7] = byte(value)
}

func validPayload() []byte {
	payload := make([]byte, 108+3309+1952+64)
	putU32BE(payload, 0, 1)
	putU64BE(payload, 4, 7)

	var sender state.Address
	var recipient state.Address
	for i := 0; i < 32; i++ {
		sender[i] = byte(i + 1)
		recipient[i] = byte(i + 33)
	}
	copy(payload[12:44], sender[:])
	copy(payload[44:76], recipient[:])

	putU64BE(payload, 76, 500000)
	putU64BE(payload, 84, 2)
	putU32BE(payload, 92, 108)
	putU32BE(payload, 96, 108+3309)
	putU32BE(payload, 100, 108+3309+1952)
	putU32BE(payload, 104, 64)
	return payload
}

func TestZeroCopyPayloadVerifierValidPayload(t *testing.T) {
	payload := validPayload()
	verifier := vm.NewZeroCopyPayloadVerifier(payload)
	header, err := verifier.VerifyPayloadBoundaries()
	if err != nil {
		t.Fatalf("expected payload to pass verification: %v", err)
	}
	if header.Version != 1 || header.Nonce != 7 || header.GasLimit != 500000 || header.GasPrice != 2 {
		t.Fatal("unexpected parsed header values")
	}
	for i := 0; i < 32; i++ {
		if header.Sender[i] != payload[12+i] || header.Recipient[i] != payload[44+i] {
			t.Fatal("address windows parsed incorrectly")
		}
	}
}

func TestZeroCopyPayloadVerifierRejectsMalformedHeader(t *testing.T) {
	_, err := vm.NewZeroCopyPayloadVerifier(make([]byte, 107)).VerifyPayloadBoundaries()
	if !errors.Is(err, vm.ErrMalformedPayloadHeader) {
		t.Fatalf("expected malformed header error, got %v", err)
	}
}

func TestZeroCopyPayloadVerifierRejectsSignatureWindowOverflow(t *testing.T) {
	payload := validPayload()
	putU32BE(payload, 92, uint32(len(payload)-3308))
	_, err := vm.NewZeroCopyPayloadVerifier(payload).VerifyPayloadBoundaries()
	if !errors.Is(err, vm.ErrSignatureSizeMismatch) {
		t.Fatalf("expected signature mismatch error, got %v", err)
	}
}

func TestZeroCopyPayloadVerifierRejectsPubKeyWindowOverflow(t *testing.T) {
	payload := validPayload()
	putU32BE(payload, 96, uint32(len(payload)-1951))
	_, err := vm.NewZeroCopyPayloadVerifier(payload).VerifyPayloadBoundaries()
	if !errors.Is(err, vm.ErrPublicKeySizeMismatch) {
		t.Fatalf("expected public key mismatch error, got %v", err)
	}
}

func TestZeroCopyPayloadVerifierRejectsDataWindowOverflow(t *testing.T) {
	payload := validPayload()
	putU32BE(payload, 100, uint32(len(payload)-10))
	putU32BE(payload, 104, 11)
	_, err := vm.NewZeroCopyPayloadVerifier(payload).VerifyPayloadBoundaries()
	if !errors.Is(err, vm.ErrMalformedPayloadHeader) {
		t.Fatalf("expected malformed payload error, got %v", err)
	}
}

func TestVMExecuteRunsPayloadVerifierBeforeInterpreter(t *testing.T) {
	engine := vm.New(vm.Config{}, state.NewMemoryDB())
	ctx := vm.NewCallContext(vm.Address{}, vm.Address{}, []byte{1, 2, 3}, 100000)
	_, err := engine.Execute(ctx, []byte{byte(vm.STOP)})
	if !errors.Is(err, vm.ErrMalformedPayloadHeader) {
		t.Fatalf("expected payload verifier error, got %v", err)
	}
}
