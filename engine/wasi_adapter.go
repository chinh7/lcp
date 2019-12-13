package engine

import (
	"encoding/binary"
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

func wasiEnvironSizesGet(vm *vm.VM, args ...uint64) uint64 {
	countPtr := args[0]
	bufSizePtr := args[1]
	env := map[string]string{}
	totalSize := 0
	for key, value := range env {
		totalSize += len(key) + len(value) + 2
	}

	// wasm32 size_t = 32bit
	countBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(countBytes, uint32(len(env)))
	vm.MemWrite(countBytes, int(countPtr))

	bufSizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufSizeBytes, uint32(totalSize))
	vm.MemWrite(bufSizeBytes, int(bufSizePtr))

	return 0 // __WASI_ESUCCESS
}

func wasiEnvironGet(vm *vm.VM, args ...uint64) uint64 {
	pointersPtr := args[0]
	envPtr := uint32(args[1])
	env := map[string]string{}
	for key, value := range env {
		envPtrBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(envPtrBytes, envPtr)
		vm.MemWrite(envPtrBytes, int(pointersPtr))
		pointersPtr += 4 // 32 bytes advancement
		envBytes := []byte(key + "=" + value)
		vm.MemWrite(envBytes, int(envPtr))
		envPtr += uint32(len(envBytes))
	}
	zero := make([]byte, 4)
	vm.MemWrite(zero, int(pointersPtr))
	return 0 // __WASI_ESUCCESS
}

func wasiProcExit(vm *vm.VM, args ...uint64) (uint64, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("invalid proc_exit argument")
	}
	return 0, fmt.Errorf("process exit with code: %d", args[0])
}

func wasiProcRaise(vm *vm.VM, args ...uint64) (uint64, error) {
	return wasiProcRaise(vm, args...)
}
