package tests

import (
	"testing"

	pq "github.com/GXQS/SmartContracts/pkg/crypto"
)

func TestEntropyValidation(t *testing.T) {
	if err := pq.ValidateEntropy(make([]byte, 32)); err == nil {
		t.Fatal("expected all-zero entropy rejection")
	}
	good := make([]byte, 32)
	good[0] = 1
	if err := pq.ValidateEntropy(good); err != nil {
		t.Fatalf("unexpected entropy validation error: %v", err)
	}
}
