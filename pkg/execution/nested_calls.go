package execution

import "github.com/GXQS/SmartContracts/pkg/vm"

func NextDepth(parent vm.CallContext) vm.CallContext {
	child := parent
	child.Depth++
	return child
}
