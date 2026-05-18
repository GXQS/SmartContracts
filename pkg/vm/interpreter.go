package vm

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/GXQS/SmartContracts/pkg/gas"
	"github.com/GXQS/SmartContracts/pkg/state"
)

type Interpreter struct {
	cfg      Config
	schedule gas.Schedule
	db       state.Database
	tracer   Tracer
	limits   RuntimeLimits
}

func NewInterpreter(cfg Config, schedule gas.Schedule, db state.Database, tracer Tracer, limits RuntimeLimits) *Interpreter {
	return &Interpreter{cfg: cfg, schedule: schedule, db: db, tracer: tracer, limits: limits}
}

func (in *Interpreter) Run(ctx CallContext, code []byte) (Receipt, error) {
	if err := ctx.Validate(in.limits); err != nil {
		return Receipt{}, err
	}
	st := NewStack(1024)
	mem := NewMemory(in.cfg.MaxMemory)
	gm := NewGasMeter(ctx.Gas, in.schedule)
	reverter := NewReverter(in.db)
	reverter.Begin()
	defer mem.Zeroize()

	var pc uint64
	var steps uint64
	for pc < uint64(len(code)) {
		if steps > in.limits.MaxSteps {
			reverter.Revert()
			return Receipt{}, errors.New("max steps exceeded")
		}
		steps++
		op := OpCode(code[pc])
		gasBefore := gm.Remaining()
		memBefore := mem.Size()
		pc++

		switch op {
		case STOP:
			reverter.Commit()
			return in.finalize(gm, false, nil), nil
		case PUSH1:
			if pc >= uint64(len(code)) {
				reverter.Revert()
				return Receipt{}, errors.New("push out of bounds")
			}
			if err := st.Push(uint64(code[pc])); err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			pc++
		case POP:
			if _, err := st.Pop(); err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
		case ADD, SUB:
			b, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			a, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			var r uint64
			if op == ADD {
				r = a + b
			} else {
				r = a - b
			}
			if err := st.Push(r); err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
		case MSTORE:
			offset, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			value, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			buf := make([]byte, 8)
			binary.BigEndian.PutUint64(buf, value)
			if err := mem.Store(offset, buf); err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
		case MLOAD:
			offset, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			buf, err := mem.Load(offset, 8)
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			if err := st.Push(binary.BigEndian.Uint64(buf)); err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
		case SSTORE:
			key, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			value, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			k := make([]byte, 8)
			v := make([]byte, 8)
			binary.BigEndian.PutUint64(k, key)
			binary.BigEndian.PutUint64(v, value)
			in.db.SetStorage([32]byte(ctx.Callee), k, v)
		case SLOAD:
			key, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			k := make([]byte, 8)
			binary.BigEndian.PutUint64(k, key)
			value := in.db.GetStorage([32]byte(ctx.Callee), k)
			if len(value) == 0 {
				if err := st.Push(0); err != nil {
					reverter.Revert()
					return Receipt{}, err
				}
			} else if err := st.Push(binary.BigEndian.Uint64(value)); err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
		case RETURN:
			offset, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			size, err := st.Pop()
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			data, err := mem.Load(offset, size)
			if err != nil {
				reverter.Revert()
				return Receipt{}, err
			}
			reverter.Commit()
			return in.finalize(gm, false, data), nil
		case REVERT:
			reverter.Revert()
			return in.finalize(gm, true, nil), nil
		default:
			reverter.Revert()
			return Receipt{}, fmt.Errorf("invalid opcode 0x%02x", byte(op))
		}
		memAfter := mem.Size()
		if !gm.Consume(op, memBefore, memAfter) {
			reverter.Revert()
			return Receipt{}, errors.New("out of gas")
		}
		in.tracer.OnStep(TraceEvent{PC: pc, Op: op, GasBefore: gasBefore, GasAfter: gm.Remaining()})
	}
	reverter.Commit()
	return in.finalize(gm, false, nil), nil
}

func (in *Interpreter) finalize(gm *GasMeter, reverted bool, ret []byte) Receipt {
	root := in.db.Root()
	var h Hash
	copy(h[:], root)
	return Receipt{GasUsed: gm.Used(), GasRefund: gm.Refund(), Reverted: reverted, ReturnData: ret, StateRoot: h}
}
