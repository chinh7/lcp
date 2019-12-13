package engine

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/vertexdlt/vertexvm/vm"
)

const maxEngineCallDepth = 16

func readAt(vm *vm.VM, ptr, size int) ([]byte, error) {
	data := make([]byte, size)
	_, err := vm.MemRead(data, ptr)
	return data, err
}

func (engine *Engine) chainPrintBytes(vm *vm.VM, args ...uint64) (uint64, error) {
	ptr, size := int(args[0]), int(args[1])
	bytes, err := readAt(vm, ptr, size)
	log.Println(string(bytes))
	return 0, err
}

func (engine *Engine) chainStorageSet(vm *vm.VM, args ...uint64) (uint64, error) {
	keyPtr, keySize := int(args[0]), int(args[1])
	valuePtr, valueSize := int(args[2]), int(args[3])
	// Burn gas before actually execute
	cost := engine.gasPolicy.GetCostForStorage(valueSize)
	err := vm.BurnGas(cost)
	if err != nil {
		return 0, err
	}
	key, err := readAt(vm, keyPtr, keySize)
	if err != nil {
		return 0, err
	}
	value, err := readAt(vm, valuePtr, valueSize)
	if err != nil {
		return 0, err
	}
	err = engine.account.SetStorage(key, value)
	return 0, err
}

func (engine *Engine) chainStorageGet(vm *vm.VM, args ...uint64) (uint64, error) {
	keyPtr, keySize := int(args[0]), int(args[1])
	key, err := readAt(vm, keyPtr, keySize)
	if err != nil {
		return 0, err
	}
	valuePtr := int(uint32(args[2]))
	value, err := engine.account.GetStorage(key)
	if err == nil {
		if len(value) == 0 {
			valuePtr = 0
		} else {
			_, err = vm.MemWrite(value, valuePtr)
		}
	}
	return uint64(valuePtr), err
}

func (engine *Engine) chainStorageSizeGet(vm *vm.VM, args ...uint64) (uint64, error) {
	keyPtr, keySize := int(args[0]), int(args[1])
	key, err := readAt(vm, keyPtr, keySize)
	if err != nil {
		return 0, err
	}
	value, err := engine.account.GetStorage(key)
	return uint64(len(value)), err
}

func (engine *Engine) chainGetCaller(vm *vm.VM, args ...uint64) (uint64, error) {
	_, err := vm.MemWrite(engine.caller[:], int(args[0]))
	return 0, err
}

func (engine *Engine) chainGetCreator(vm *vm.VM, args ...uint64) (uint64, error) {
	creator := engine.account.Creator
	_, err := vm.MemWrite(creator[:], int(args[0]))
	return 0, err
}

func (engine *Engine) chainPtrArgSizeGet(vm *vm.VM, args ...uint64) (uint64, error) {
	size, err := engine.ptrArgSizeGet(int(args[0]))
	return uint64(size), err
}

func (engine *Engine) chainPtrArgSizeSet(vm *vm.VM, args ...uint64) (uint64, error) {
	engine.ptrArgSizeMap[int(args[0])] = int(args[1])
	return 0, nil
}

func (engine *Engine) chainMethodBind(vm *vm.VM, args ...uint64) (uint64, error) {
	contractAddrBytes, err := readAt(vm, int(args[0]), crypto.AddressLength)
	if err != nil {
		return 0, err
	}
	contractAddr := crypto.AddressFromBytes(contractAddrBytes)

	invokedMethodBytes, err := readAt(vm, int(args[1]), int(args[2]))
	if err != nil {
		return 0, err
	}

	invokedMethod := string(invokedMethodBytes[:len(invokedMethodBytes)-1])
	aliasMethodBytes, err := readAt(vm, int(args[3]), int(args[4]))
	if err != nil {
		return 0, err
	}
	aliasMethod := string(aliasMethodBytes[:len(aliasMethodBytes)-1])
	engine.methodLookup[aliasMethod] = &foreignMethod{contractAddr, invokedMethod}
	return 0, nil
}

func (engine *Engine) handleInvokeAlias(foreignMethod *foreignMethod, vm *vm.VM, args ...uint64) (uint64, error) {
	if engine.callDepth+1 > maxEngineCallDepth {
		return 0, errors.New("call depth limit reached")
	}

	foreignAccount, err := engine.state.GetAccount(foreignMethod.contractAddress)
	if err != nil {
		return 0, err
	}
	contract, err := foreignAccount.GetContract()
	if err != nil {
		return 0, err
	}
	function, err := contract.Header.GetFunction(foreignMethod.name)
	if err != nil {
		return 0, err
	}
	var values [][]byte
	var bytes []byte
	for i, param := range function.Parameters {
		if param.IsArray {
			argPtr := int(args[i])
			size, _ := engine.ptrArgSizeGet(int(args[i]))
			bytes, err = readAt(vm, argPtr, size)
			if err != nil {
				return 0, err
			}
		} else {
			if param.Type.IsPointer() {
				argPtr := int(args[i])
				size := param.Type.GetMemorySize()
				bytes, err = readAt(vm, argPtr, size)
				if err != nil {
					return 0, err
				}
			} else {
				binary.LittleEndian.PutUint64(bytes, args[i])
				size := param.Type.GetMemorySize()
				bytes = bytes[:size]
			}
		}
		values = append(values, bytes)

	}
	methodArgs, err := abi.EncodeFromBytes(function.Parameters, values)
	if err != nil {
		return 0, err
	}

	account, err := engine.state.GetAccount(foreignMethod.contractAddress)
	if err != nil {
		return 0, err
	}
	// TODO memcheck
	childEngine := engine.NewChildEngine(account)
	childEngine.setStats(engine.callDepth+1, engine.memAggr+vm.MemSize())
	return childEngine.Ignite(foreignMethod.name, methodArgs)
}

func (engine *Engine) handleEmitEvent(event *abi.Event, vm *vm.VM, args ...uint64) (uint64, error) {
	attributes := common.KVPairs{}
	address := engine.account.GetAddress()
	attributes = append(attributes, common.KVPair{
		Key:   []byte("address"),
		Value: address[:],
	})
	for i, param := range event.Parameters {
		var value []byte
		var err error
		if param.Type.IsPointer() {
			paramPtr := int(uint32(args[i]))
			size := param.Type.GetMemorySize()
			value, err = readAt(vm, paramPtr, size)
			if err != nil {
				return 0, err
			}
		} else {
			size := abi.Uint64.GetMemorySize()
			value = make([]byte, size)
			binary.BigEndian.PutUint64(value, args[i])
		}
		if param.Name != "address" {
			attributes = append(attributes, common.KVPair{
				Key:   []byte(param.Name),
				Value: value,
			})
		}
	}
	engine.pushEvent(types.Event{
		Type:       EventPrefix + event.Name,
		Attributes: attributes,
	})
	return 0, nil
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
		case "chain_method_bind":
			return engine.chainMethodBind
		case "chain_arg_size_get":
			return engine.chainPtrArgSizeGet
		case "chain_arg_size_set":
			return engine.chainPtrArgSizeSet
		default:
			contract, _ := engine.account.GetContract()
			if event, err := contract.Header.GetEvent(name); err == nil {
				return func(vm *vm.VM, args ...uint64) (uint64, error) {
					return engine.handleEmitEvent(event, vm, args...)
				}
			}

			if foreignMethod, ok := engine.methodLookup[name]; ok {
				return func(vm *vm.VM, args ...uint64) (uint64, error) {
					return engine.handleInvokeAlias(foreignMethod, vm, args...)
				}
			}
			panic(fmt.Errorf("unknown import resolved: %s", name))
		}
	case "wasi_unstable":
		return wasiUnstableHandler(name)
	}
	return nil
}
