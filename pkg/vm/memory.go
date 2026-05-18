package vm

import (
	"errors"
)

type Memory struct {
	buf      []byte
	maxBytes uint64
}

func NewMemory(max uint64) *Memory {
	if max == 0 {
		max = 1 << 20
	}
	return &Memory{maxBytes: max}
}

func (m *Memory) ensure(size uint64) error {
	if size > m.maxBytes {
		return errors.New("memory limit exceeded")
	}
	if size <= uint64(len(m.buf)) {
		return nil
	}
	newBuf := make([]byte, size)
	copy(newBuf, m.buf)
	m.buf = newBuf
	return nil
}

func (m *Memory) Store(offset uint64, data []byte) error {
	end := offset + uint64(len(data))
	if end < offset {
		return errors.New("memory overflow")
	}
	if err := m.ensure(end); err != nil {
		return err
	}
	copy(m.buf[offset:end], data)
	return nil
}

func (m *Memory) Load(offset, size uint64) ([]byte, error) {
	end := offset + size
	if end < offset {
		return nil, errors.New("memory overflow")
	}
	if err := m.ensure(end); err != nil {
		return nil, err
	}
	out := make([]byte, size)
	copy(out, m.buf[offset:end])
	return out, nil
}

func (m *Memory) Zeroize() {
	for i := range m.buf {
		m.buf[i] = 0
	}
}

func (m *Memory) Size() uint64 { return uint64(len(m.buf)) }
