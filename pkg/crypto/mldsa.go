package crypto

import (
	"errors"

	"github.com/cloudflare/circl/sign/mldsa/mldsa65"
)

const (
	MLDSASignatureBytes = 3309
	MLDSAPublicKeyBytes = 1952
)

func VerifyMLDSA(publicKey, message, signature []byte) error {
	if len(publicKey) != MLDSAPublicKeyBytes {
		return errors.New("invalid ML-DSA public key size")
	}
	if len(signature) != MLDSASignatureBytes {
		return errors.New("invalid ML-DSA signature size")
	}
	var pk mldsa65.PublicKey
	if err := pk.UnmarshalBinary(publicKey); err != nil {
		return err
	}
	if !mldsa65.Verify(&pk, message, nil, signature) {
		return errors.New("ML-DSA verification failed")
	}
	return nil
}
