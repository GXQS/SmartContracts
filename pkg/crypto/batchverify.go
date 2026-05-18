package crypto

type VerifyRequest struct {
	PublicKey []byte
	Message   []byte
	Signature []byte
}

func BatchVerifyMLDSA(requests []VerifyRequest) error {
	for _, req := range requests {
		if err := VerifyMLDSA(req.PublicKey, req.Message, req.Signature); err != nil {
			return err
		}
	}
	return nil
}
