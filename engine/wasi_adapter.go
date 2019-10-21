package engine

import (
	"encoding/binary"

	"github.com/vertexdlt/vertexvm/vm"
)

func wasiDefaultHandler(vm *vm.VM, args ...uint64) uint64 {
	return 52 // __WASI_ENOSYS
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
	binary.LittleEndian.PutUint32(vm.GetMemory()[countPtr:], uint32(len(env)))
	binary.LittleEndian.PutUint32(vm.GetMemory()[bufSizePtr:], uint32(totalSize))
	return 0 // __WASI_ESUCCESS
}

func wasiEnvironGet(vm *vm.VM, args ...uint64) uint64 {
	pointersPtr := args[0]
	envPtr := uint32(args[1])
	env := map[string]string{}
	for key, value := range env {
		binary.LittleEndian.PutUint32(vm.GetMemory()[pointersPtr:], envPtr)
		pointersPtr += 4 // 32 bytes advancement
		envBytes := []byte(key + "=" + value)
		copy(vm.GetMemory()[envPtr:], envBytes)
		envPtr += uint32(len(envBytes))
	}
	binary.LittleEndian.PutUint32(vm.GetMemory()[pointersPtr:], 0)
	return 0 // __WASI_ESUCCESS
}

func wasiProcExit(vm *vm.VM, args ...uint64) uint64 {
	return 0
}

func wasiProcRaise(vm *vm.VM, args ...uint64) uint64 {
	return 0
}
