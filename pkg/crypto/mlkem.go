package crypto

import (
	"crypto/rand"
	"errors"

	"github.com/cloudflare/circl/kem/mlkem/mlkem768"
)

const MLKEMCiphertextBytes = 1088

func EncapsulateMLKEM(publicKey []byte) (ciphertext, sharedSecret []byte, err error) {
	pk, err := mlkem768.Scheme().UnmarshalBinaryPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	seed := make([]byte, mlkem768.EncapsulationSeedSize)
	if _, err = rand.Read(seed); err != nil {
		return nil, nil, err
	}
	ct, ss, err := mlkem768.Scheme().EncapsulateDeterministically(pk, seed)
	if err != nil {
		return nil, nil, err
	}
	return ct, ss, nil
}

func DecapsulateMLKEM(privateKey, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) != MLKEMCiphertextBytes {
		return nil, errors.New("invalid ML-KEM ciphertext size")
	}
	sk, err := mlkem768.Scheme().UnmarshalBinaryPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return mlkem768.Scheme().Decapsulate(sk, ciphertext)
}
