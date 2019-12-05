package engine

import (
	"encoding/binary"
	"fmt"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/vertexdlt/vertexvm/vm"
)

func readAt(vm *vm.VM, ptr, size int) []byte {
	data := make([]byte, size)
	copy(data, vm.GetMemory()[ptr:ptr+size])
	return data
}

func (engine *Engine) chainPrintBytes(vm *vm.VM, args ...uint64) uint64 {
	ptr := int(uint32(args[0]))
	size := int(uint32(args[1]))
	bytes := readAt(vm, ptr, size)
	fmt.Println(string(bytes))
	return 0
}

func (engine *Engine) chainStorageSet(vm *vm.VM, args ...uint64) uint64 {
	keyPtr := int(uint32(args[0]))
	keySize := int(uint32(args[1]))
	valuePtr := int(uint32(args[2]))
	valueSize := int(uint32(args[3]))
	// Burn gas before actually execute
	cost := engine.gasPolicy.GetCostForStorage(valueSize)
	vm.BurnGas(cost)
	key := readAt(vm, keyPtr, keySize)
	value := readAt(vm, valuePtr, valueSize)
	engine.account.SetStorage(key, value)
	return 0
}

func (engine *Engine) chainStorageGet(vm *vm.VM, args ...uint64) uint64 {
	keyPtr := int(uint32(args[0]))
	keySize := int(uint32(args[1]))
	key := readAt(vm, keyPtr, keySize)
	valuePtr := int(uint32(args[2]))
	value, err := engine.account.GetStorage(key)
	if err == nil && len(value) > 0 {
		copy(vm.GetMemory()[valuePtr:], value)
	} else {
		valuePtr = 0
	}
	return uint64(valuePtr)
}

func (engine *Engine) chainStorageSizeGet(vm *vm.VM, args ...uint64) uint64 {
	keyPtr := int(uint32(args[0]))
	keySize := int(uint32(args[1]))
	key := readAt(vm, keyPtr, keySize)
	value, _ := engine.account.GetStorage(key)
	return uint64(len(value))
}

func (engine *Engine) chainGetCaller(vm *vm.VM, args ...uint64) uint64 {
	ptr := int(uint32(args[0]))
	copy(vm.GetMemory()[ptr:], engine.caller[:])
	return 0
}

func (engine *Engine) chainGetCreator(vm *vm.VM, args ...uint64) uint64 {
	ptr := int(uint32(args[0]))
	creator := engine.account.Creator
	copy(vm.GetMemory()[ptr:], creator[:])
	return 0
}

func (engine *Engine) handleEmitEvent(event *abi.Event, vm *vm.VM, args ...uint64) uint64 {
	attributes := common.KVPairs{}
	for i, param := range event.Parameters {
		var value []byte
		if param.Type.IsPointer() {
			paramPtr := int(uint32(args[i]))
			size, _ := param.Type.GetMemorySize()
			value = readAt(vm, paramPtr, size)
		} else {
			size, _ := abi.Uint64.GetMemorySize()
			value = make([]byte, size)
			binary.BigEndian.PutUint64(value, args[i])
		}
		attributes = append(attributes, common.KVPair{
			Key:   []byte(param.Name),
			Value: value,
		})
	}
	engine.events = append(engine.events, types.Event{
		Type:       EventPrefix + event.Name,
		Attributes: attributes,
	})
	return 0
}

// GetFunction get host function for WebAssembly
func (engine *Engine) GetFunction(module, name string) vm.HostFunction {
	switch module {
	case "env":
		switch name {
		case "chain_print_bytes":
			return engine.chainPrintBytes
		case "chain_storage_set":
			return engine.chainStorageSet
		case "chain_storage_get":
			return engine.chainStorageGet
		case "chain_storage_size_get":
			return engine.chainStorageSizeGet
		case "chain_get_caller":
			return engine.chainGetCaller
		case "chain_get_creator":
			return engine.chainGetCreator
		default:
			contract, _ := engine.account.GetContract()
			if event, err := contract.Header.GetEvent(name); err == nil {
				return func(vm *vm.VM, args ...uint64) uint64 {
					return engine.handleEmitEvent(event, vm, args...)
				}
			}
			panic(fmt.Errorf("unknown import resolved: %s", name))
		}
	case "wasi_unstable":
		return wasiDefaultHandler
	}
	return nil
}
