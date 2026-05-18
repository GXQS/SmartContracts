package vm

type Host interface {
	Call(CallContext, []byte) (Receipt, error)
}

type Dispatcher struct {
	host Host
}

func NewDispatcher(host Host) Dispatcher { return Dispatcher{host: host} }

func (d Dispatcher) Dispatch(ctx CallContext, payload []byte) (Receipt, error) {
	if d.host == nil {
		return Receipt{}, nil
	}
	return d.host.Call(ctx, payload)
}
