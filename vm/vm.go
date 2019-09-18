package vm

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/QuoineFinancial/vertex/storage"
	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

var accountState *storage.AccountState
var events [][]byte

// Call executes a contract given its code, method, and arguments
func Call(state *storage.AccountState, method string, methodArgs ...interface{}) (interface{}, [][]byte) {
	accountState = state
	events = make([][]byte, 0)
	programReader := bytes.NewReader(accountState.GetCode())

	m, err := wasm.ReadModule(programReader, func(n string) (*wasm.Module, error) { return resolveImports(n) })
	if err != nil {
		log.Fatalf("could not read module: %v", err)
	}
	initMemory(m)

	funcID := int64(m.Export.Entries[method].Index)
	vm, err := exec.NewVM(m)

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
		log.Fatalf("Error executing the default function: %v", err)
	}
	log.Println("return value = ", ret)
	return ret, events
}
