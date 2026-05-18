package state

import "errors"

type Proof struct {
	Root  []byte
	Key   []byte
	Value []byte
}

func (t *Trie) GenerateProof(key []byte) (Proof, error) {
	v, ok := t.Get(key)
	if !ok {
		return Proof{}, errors.New("key not found")
	}
	return Proof{Root: t.RootHash(), Key: append([]byte(nil), key...), Value: v}, nil
}

func VerifyProof(p Proof) bool {
	t := NewTrie()
	t.Update(p.Key, p.Value)
	root := t.RootHash()
	if len(root) != len(p.Root) {
		return false
	}
	for i := range root {
		if root[i] != p.Root[i] {
			return false
		}
	}
	return true
}
