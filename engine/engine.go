package engine

import (
	"encoding/binary"
	"errors"
	"log"

	"github.com/QuoineFinancial/vertex/storage"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/vertexdlt/vertexvm/vm"
)

// Engine is space to execute function
type Engine struct {
	event   types.Event
	account *storage.Account
}

// NewEngine return new instance of Engine
func NewEngine(account *storage.Account) *Engine {
	return &Engine{
		event:   types.Event{},
		account: account,
	}
}

// Ignite executes a contract given its code, method, and arguments
func (engine *Engine) Ignite(method string, methodArgs ...interface{}) (interface{}, [][]byte, error) {
	events := make([][]byte, 0)

	vm, err := vm.NewVM(engine.account.GetCode(), engine)
	if err != nil {
		return nil, [][]byte{}, err
	}
	funcID, ok := vm.GetFunctionIndex(method)
	if !ok {
		return nil, [][]byte{}, errors.New("Cannot find invoke function")
	}

	val, _ := vm.Module.ExecInitExpr(vm.Module.GetGlobal(int(vm.Module.Export.Entries["__data_end"].Desc.Idx)).Init)
	offset := int(val.(int32))
	params := make([]uint64, len(methodArgs))
	for i := range methodArgs {
		var arg uint64
		bytes := methodArgs[i].([]byte)
		if len(bytes) <= 8 {
			uintBytes := make([]byte, 8)
			copy(uintBytes[8-len(bytes):], bytes)
			arg = binary.BigEndian.Uint64(uintBytes[:])
		} else {
			copy(vm.GetMemory()[offset:], bytes)
			arg = uint64(offset)
			offset += len(bytes)
		}
		params[i] = arg
	}
	ret := vm.Invoke(funcID, params...)
	log.Println("return value = ", ret)
	return ret, events, nil
}
