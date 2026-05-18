package tests

import (
	"testing"

	"github.com/GXQS/SmartContracts/pkg/state"
)

func TestTrieProofRoundTrip(t *testing.T) {
	tr := state.NewTrie()
	tr.Update([]byte("k"), []byte("v"))
	proof, err := tr.GenerateProof([]byte("k"))
	if err != nil {
		t.Fatalf("proof generation failed: %v", err)
	}
	if !state.VerifyProof(proof) {
		t.Fatal("proof verification failed")
	}
}
