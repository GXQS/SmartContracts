package execution

import "github.com/GXQS/SmartContracts/pkg/vm"

type ReceiptEnvelope struct {
	TxHash  vm.Hash
	Receipt vm.Receipt
}
