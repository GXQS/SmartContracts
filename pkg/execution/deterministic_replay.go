package execution

import "github.com/GXQS/SmartContracts/pkg/vm"

func DeterministicReplay(vmEngine *vm.VM, contexts []vm.CallContext, payloads [][]byte) ([]vm.Receipt, error) {
	out := make([]vm.Receipt, 0, len(contexts))
	for i := range contexts {
		receipt, err := vmEngine.Execute(contexts[i], payloads[i])
		if err != nil {
			return nil, err
		}
		out = append(out, receipt)
	}
	return out, nil
}
