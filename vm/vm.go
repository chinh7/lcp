package vm

import (
	"bytes"
	"log"

	"github.com/vertexdlt/vertex/storage"
	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

var accountState *storage.AccountState

// Call executes a contract given its code, method, and arguments
func Call(state *storage.AccountState, method string, methodArgs ...interface{}) interface{} {
	accountState = state
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
		arg, ok := methodArgs[i].(int64)
		if !ok {
			value := []byte(methodArgs[i].(string))
			log.Println("malloc for", methodArgs[i].(string))
			arg = int64(malloc(int32(len(value))))
			proc.WriteAt(value, arg)
		}
		params[i] = uint64(arg)
	}

	ret, err := vm.ExecCode(funcID, params...)
	if err != nil {
		log.Fatalf("Error executing the default function: %v", err)
	}
	log.Println("return value = ", ret)
	return ret
}
