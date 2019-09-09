package vm

import (
	"bytes"
	"log"
	"strconv"

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
		stringArg := string(methodArgs[i].([]byte))
		arg, err := strconv.ParseInt(stringArg, 10, 64)
		if err != nil {
			value := methodArgs[i].([]byte)
			log.Println("malloc for", string(methodArgs[i].([]byte)))
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
	return ret, events
}
