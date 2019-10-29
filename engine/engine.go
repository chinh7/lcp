package engine

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/vertexdlt/vertexvm/vm"
)

const (
	// ExportSecDataEnd is wasm export section key for __data_end
	ExportSecDataEnd = "__data_end"
)

// Engine is space to execute function
type Engine struct {
	event   types.Event
	account *storage.Account
	caller  crypto.Address
}

// NewEngine return new instance of Engine
func NewEngine(account *storage.Account, caller crypto.Address) *Engine {
	return &Engine{
		event:   types.Event{},
		account: account,
		caller:  caller,
	}
}

// Ignite executes a contract given its code, method, and arguments
func (engine *Engine) Ignite(method string, methodArgs []byte) (interface{}, [][]byte, error) {
	events := make([][]byte, 0)
	contract := engine.account.GetContract()
	header, dataIndex, err := abi.DecodeHeader(contract)
	vm, err := vm.NewVM(contract[dataIndex:], engine)
	if err != nil {
		return nil, [][]byte{}, err
	}
	funcID, ok := vm.GetFunctionIndex(method)
	if !ok {
		return nil, [][]byte{}, errors.New("Cannot find invoke function")
	}

	val, _ := vm.Module.ExecInitExpr(vm.Module.GetGlobal(int(vm.Module.Export.Entries[ExportSecDataEnd].Desc.Idx)).Init)
	offset := int(val.(int32))

	function, err := header.GetFunction(method)
	if err != nil {
		return nil, [][]byte{}, err
	}

	decodedBytes, err := abi.DecodeToBytes(function.Parameters, methodArgs)
	if err != nil {
		return nil, [][]byte{}, err
	}

	arguments, err := loadArguments(vm, decodedBytes, function.Parameters, offset)
	if err != nil {
		return nil, [][]byte{}, err
	}

	ret := vm.Invoke(funcID, arguments...)
	log.Println("return value = ", ret)
	return ret, events, nil
}

func loadArguments(vm *vm.VM, byteArgs [][]byte, params []abi.Parameter, offset int) ([]uint64, error) {
	var args = make([]uint64, len(byteArgs))
	byteSize := 0
	for _, bytes := range byteArgs {
		byteSize += len(bytes)
	}
	if byteSize > 1024 {
		return []uint64{}, fmt.Errorf("arguments byte size exceeds limit")
	}
	for i, bytes := range byteArgs {
		isArray := params[i].IsArray || params[i].Type.String() == "address"
		if isArray {
			copy(vm.GetMemory()[offset:], bytes)
			args[i] = uint64(offset)
			offset += len(bytes)
		} else {
			buffer := make([]byte, 8)
			copy(buffer, bytes)
			args[i] = binary.LittleEndian.Uint64(buffer)
		}
	}
	return args, nil
}
