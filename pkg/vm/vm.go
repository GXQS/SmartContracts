package vm

import (
	"errors"

	"github.com/GXQS/SmartContracts/pkg/gas"
	"github.com/GXQS/SmartContracts/pkg/state"
)

type Hash [32]byte
type Address = state.Address

type Config struct {
	MaxCallDepth int
	MaxMemory    uint64
}

type VM struct {
	cfg      Config
	stateDB  state.Database
	tracer   Tracer
	verifier Verifier
	limits   RuntimeLimits
}

func New(cfg Config, db state.Database) *VM {
	if cfg.MaxCallDepth <= 0 {
		cfg.MaxCallDepth = 16
	}
	if cfg.MaxMemory == 0 {
		cfg.MaxMemory = 1 << 20
	}
	return &VM{
		cfg:      cfg,
		stateDB:  db,
		tracer:   NoopTracer{},
		verifier: DefaultVerifier{},
		limits: RuntimeLimits{
			MaxCallDepth: cfg.MaxCallDepth,
			MaxMemory:    cfg.MaxMemory,
			MaxSteps:     100000,
		},
	}
}

func (v *VM) WithTracer(t Tracer) *VM {
	if t != nil {
		v.tracer = t
	}
	return v
}

func (v *VM) Execute(ctx CallContext, code []byte) (Receipt, error) {
	if _, err := v.verifier.VerifyPayload(ctx.Input); err != nil {
		return Receipt{}, err
	}
	if err := v.verifier.VerifyBytecode(code); err != nil {
		return Receipt{}, err
	}
	if ctx.Gas <= 0 {
		return Receipt{}, errors.New("insufficient gas")
	}
	interpreter := NewInterpreter(v.cfg, gas.DefaultSchedule(), v.stateDB, v.tracer, v.limits)
	return interpreter.Run(ctx, code)
}
