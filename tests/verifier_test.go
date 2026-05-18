package tests

import (
	"errors"
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
	"github.com/GXQS/SmartContracts/pkg/vm"
)

func buildPayloadWithData() []byte {
	payload := vm.BuildBoundaryPayload(64)
	vm.WriteU32BE(payload, vm.HeaderVersionOffset, 1)
	vm.WriteU64BE(payload, vm.HeaderNonceOffset, 7)

	var sender state.Address
	var recipient state.Address
	for i := 0; i < 32; i++ {
		sender[i] = byte(i + 1)
		recipient[i] = byte(i + 33)
	}
	copy(payload[vm.HeaderSenderOffset:vm.HeaderRecipientOffset], sender[:])
	copy(payload[vm.HeaderRecipientOffset:vm.HeaderGasLimitOffset], recipient[:])

	vm.WriteU64BE(payload, vm.HeaderGasLimitOffset, 500000)
	vm.WriteU64BE(payload, vm.HeaderGasPriceOffset, 2)
	vm.WriteU32BE(payload, vm.HeaderSigOffsetOffset, vm.PayloadHeaderSize)
	vm.WriteU32BE(payload, vm.HeaderPubOffsetOffset, vm.PayloadHeaderSize+vm.MLDSA65SignatureSize)
	vm.WriteU32BE(payload, vm.HeaderDataOffset, vm.PayloadHeaderSize+vm.MLDSA65SignatureSize+vm.MLDSA65PublicKeySize)
	vm.WriteU32BE(payload, vm.HeaderDataLenOffset, 64)
	return payload
}

func TestZeroCopyPayloadVerifierValidPayload(t *testing.T) {
	payload := buildPayloadWithData()
	verifier := vm.NewZeroCopyPayloadVerifier(payload)
	header, err := verifier.VerifyPayloadBoundaries()
	if err != nil {
		t.Fatalf("expected payload to pass verification: %v", err)
	}
	if header.Version != 1 || header.Nonce != 7 || header.GasLimit != 500000 || header.GasPrice != 2 {
		t.Fatal("unexpected parsed header values")
	}
	for i := 0; i < 32; i++ {
		if header.Sender[i] != payload[vm.HeaderSenderOffset+i] || header.Recipient[i] != payload[vm.HeaderRecipientOffset+i] {
			t.Fatal("address windows parsed incorrectly")
		}
	}
}

func TestZeroCopyPayloadVerifierRejectsMalformedHeader(t *testing.T) {
	_, err := vm.NewZeroCopyPayloadVerifier(make([]byte, vm.PayloadHeaderSize-1)).VerifyPayloadBoundaries()
	if !errors.Is(err, vm.ErrMalformedPayloadHeader) {
		t.Fatalf("expected malformed header error, got %v", err)
	}
}

func TestZeroCopyPayloadVerifierRejectsSignatureWindowOverflow(t *testing.T) {
	payload := buildPayloadWithData()
	vm.WriteU32BE(payload, vm.HeaderSigOffsetOffset, uint32(len(payload)-vm.MLDSA65SignatureSize+1))
	_, err := vm.NewZeroCopyPayloadVerifier(payload).VerifyPayloadBoundaries()
	if !errors.Is(err, vm.ErrSignatureSizeMismatch) {
		t.Fatalf("expected signature mismatch error, got %v", err)
	}
}

func TestZeroCopyPayloadVerifierRejectsPubKeyWindowOverflow(t *testing.T) {
	payload := buildPayloadWithData()
	vm.WriteU32BE(payload, vm.HeaderPubOffsetOffset, uint32(len(payload)-vm.MLDSA65PublicKeySize+1))
	_, err := vm.NewZeroCopyPayloadVerifier(payload).VerifyPayloadBoundaries()
	if !errors.Is(err, vm.ErrPublicKeySizeMismatch) {
		t.Fatalf("expected public key mismatch error, got %v", err)
	}
}

func TestZeroCopyPayloadVerifierRejectsDataWindowOverflow(t *testing.T) {
	payload := buildPayloadWithData()
	vm.WriteU32BE(payload, vm.HeaderDataOffset, uint32(len(payload)-10))
	vm.WriteU32BE(payload, vm.HeaderDataLenOffset, 11)
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

	validCtx := vm.NewCallContext(vm.Address{}, vm.Address{}, validPayloadFixture(), 100000)
	if _, err := engine.Execute(validCtx, []byte{byte(vm.STOP)}); err != nil {
		t.Fatalf("expected valid ctx.Input to pass payload verifier, got %v", err)
	}
}
