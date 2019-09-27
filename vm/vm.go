package vm

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/QuoineFinancial/vertex/storage"
	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/tendermint/tendermint/abci/types"
)

// VertexVM is space to execute function
type VertexVM struct {
	event   types.Event
	account *storage.Account
}

// NewVertexVM return new instance of VertexVM
func NewVertexVM(account *storage.Account) *VertexVM {
	return &VertexVM{
		event:   types.Event{},
		account: account,
	}
}

// Call executes a contract given its code, method, and arguments
func (vertexVM *VertexVM) Call(method string, methodArgs ...interface{}) (interface{}, [][]byte, error) {
	events := make([][]byte, 0)
	programReader := bytes.NewReader(vertexVM.account.GetCode())

	m, err := wasm.ReadModule(programReader, func(n string) (*wasm.Module, error) {
		return vertexVM.resolveImports(n)
	})

	if err != nil {
		return nil, [][]byte{}, err
	}
	initMemory(m)

	funcID := int64(m.Export.Entries[method].Index)
	vm, err := exec.NewVM(m)
	if err != nil {
		return nil, [][]byte{}, err
	}

	proc := exec.NewProcess(vm)
	params := make([]uint64, len(methodArgs))
	for i := range methodArgs {
		var arg uint64
		bytes := methodArgs[i].([]byte)
		if len(bytes) <= 8 {
			uintBytes := make([]byte, 8)
			copy(uintBytes[8-len(bytes):], bytes)
			arg = binary.BigEndian.Uint64(uintBytes[:])
		} else {
			value := string(bytes)
			log.Println("malloc for", value)
			arg = uint64(malloc(int32(len(value))))
			proc.WriteAt(bytes, int64(arg))
		}
		params[i] = arg
	}

	ret, err := vm.ExecCode(funcID, params...)
	if err != nil {
		return nil, [][]byte{}, err
	}
	log.Println("return value = ", ret)
	return ret, events, nil
}
