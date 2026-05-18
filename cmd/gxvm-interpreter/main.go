package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/GXQS/SmartContracts/pkg/state"
	"github.com/GXQS/SmartContracts/pkg/vm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: gxvm-interpreter <hex-bytecode>")
		os.Exit(2)
	}

	code, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Println("invalid hex bytecode:", err)
		os.Exit(2)
	}

	engine := vm.New(vm.Config{MaxCallDepth: 16, MaxMemory: 1 << 20}, state.NewMemoryDB())
	r := vm.NewCallContext(vm.Address{}, vm.Address{}, nil, 1_000_000)
	receipt, err := engine.Execute(r, code)
	if err != nil {
		fmt.Println("execution error:", err)
		os.Exit(1)
	}

	fmt.Printf("gas_used=%d reverted=%t return=%x state_root=%x\n", receipt.GasUsed, receipt.Reverted, receipt.ReturnData, receipt.StateRoot)
}
