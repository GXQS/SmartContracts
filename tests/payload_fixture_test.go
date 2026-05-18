package tests

import "github.com/GXQS/SmartContracts/pkg/vm"

func validPayloadFixture() []byte {
	return vm.BuildBoundaryPayload(0)
}
