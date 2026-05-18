package state

import (
	"bytes"
	"crypto/sha256"
)

type node interface{ hash() []byte }

type branchNode struct {
	children [16]node
	value    []byte
}

type extensionNode struct {
	path []byte
	next node
}

type leafNode struct {
	path  []byte
	value []byte
}

func (n *branchNode) hash() []byte {
	h := sha256.New()
	for _, c := range n.children {
		if c == nil {
			h.Write([]byte{0})
			continue
		}
		h.Write(c.hash())
	}
	h.Write(n.value)
	return h.Sum(nil)
}

func (n *extensionNode) hash() []byte {
	h := sha256.New()
	h.Write([]byte{1})
	h.Write(CompactEncode(n.path, false))
	if n.next != nil {
		h.Write(n.next.hash())
	}
	return h.Sum(nil)
}

func (n *leafNode) hash() []byte {
	h := sha256.New()
	h.Write([]byte{2})
	h.Write(CompactEncode(n.path, true))
	h.Write(n.value)
	return h.Sum(nil)
}

type Trie struct{ root node }

func NewTrie() *Trie { return &Trie{} }

func Nibbles(key []byte) []byte {
	out := make([]byte, len(key)*2)
	for i, b := range key {
		out[i*2] = b >> 4
		out[i*2+1] = b & 0x0F
	}
	return out
}

func commonPrefix(a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}

func (t *Trie) Update(key, value []byte) {
	t.root = insert(t.root, Nibbles(key), value)
}

func insert(cur node, path, value []byte) node {
	if cur == nil {
		return &leafNode{path: append([]byte(nil), path...), value: append([]byte(nil), value...)}
	}
	switch n := cur.(type) {
	case *leafNode:
		p := commonPrefix(n.path, path)
		if p == len(n.path) && p == len(path) {
			n.value = append([]byte(nil), value...)
			return n
		}
		b := &branchNode{}
		if p == len(n.path) {
			b.value = append([]byte(nil), n.value...)
		} else {
			b.children[n.path[p]] = &leafNode{path: append([]byte(nil), n.path[p+1:]...), value: append([]byte(nil), n.value...)}
		}
		if p == len(path) {
			b.value = append([]byte(nil), value...)
		} else {
			b.children[path[p]] = &leafNode{path: append([]byte(nil), path[p+1:]...), value: append([]byte(nil), value...)}
		}
		if p == 0 {
			return b
		}
		return &extensionNode{path: append([]byte(nil), path[:p]...), next: b}
	case *extensionNode:
		p := commonPrefix(n.path, path)
		if p == len(n.path) {
			n.next = insert(n.next, path[p:], value)
			return n
		}
		b := &branchNode{}
		if p+1 <= len(n.path) {
			remaining := n.path[p+1:]
			if len(remaining) == 0 {
				b.children[n.path[p]] = n.next
			} else {
				b.children[n.path[p]] = &extensionNode{path: append([]byte(nil), remaining...), next: n.next}
			}
		}
		if p == len(path) {
			b.value = append([]byte(nil), value...)
		} else {
			b.children[path[p]] = &leafNode{path: append([]byte(nil), path[p+1:]...), value: append([]byte(nil), value...)}
		}
		if p == 0 {
			return b
		}
		return &extensionNode{path: append([]byte(nil), path[:p]...), next: b}
	case *branchNode:
		if len(path) == 0 {
			n.value = append([]byte(nil), value...)
			return n
		}
		idx := path[0]
		n.children[idx] = insert(n.children[idx], path[1:], value)
		return n
	default:
		return cur
	}
}

func (t *Trie) Get(key []byte) ([]byte, bool) {
	return get(t.root, Nibbles(key))
}

func get(cur node, path []byte) ([]byte, bool) {
	if cur == nil {
		return nil, false
	}
	switch n := cur.(type) {
	case *leafNode:
		if bytes.Equal(n.path, path) {
			return append([]byte(nil), n.value...), true
		}
		return nil, false
	case *extensionNode:
		if len(path) < len(n.path) || !bytes.Equal(n.path, path[:len(n.path)]) {
			return nil, false
		}
		return get(n.next, path[len(n.path):])
	case *branchNode:
		if len(path) == 0 {
			return append([]byte(nil), n.value...), n.value != nil
		}
		return get(n.children[path[0]], path[1:])
	default:
		return nil, false
	}
}

func (t *Trie) RootHash() []byte {
	if t.root == nil {
		return make([]byte, 32)
	}
	return t.root.hash()
}
