package engine

import (
	"fmt"

	"github.com/vertexdlt/vertexvm/vm"
)

func wasiUnstableHandler(name string) vm.HostFunction {
	switch name {
	case "proc_exit":
		return wasiProcExit
	case "proc_raise":
		return wasiProcRaise
	default:
		return wasiDefaultHandler
	}
}

func wasiDefaultHandler(vm *vm.VM, args ...uint64) (uint64, error) {
	return 52, nil // __WASI_ENOSYS
}

func wasiProcExit(vm *vm.VM, args ...uint64) (uint64, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("invalid proc_exit argument")
	}
	return 0, fmt.Errorf("process exit with code: %d", args[0])
}

func wasiProcRaise(vm *vm.VM, args ...uint64) (uint64, error) {
	return wasiProcExit(vm, args...)
}
