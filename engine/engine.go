package engine

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/vertexdlt/vertexvm/vm"
	vertex "github.com/vertexdlt/vertexvm/vm"
)

const (
	// ExportSecDataEnd is wasm export section key for __data_end
	ExportSecDataEnd = "__data_end"

	// EventPrefix is prefix of Type for all events emitting by engine
	EventPrefix = "engine."
)

type foreignMethod struct {
	contractAddress crypto.Address
	name            string
}

// Engine is space to execute function
type Engine struct {
	state         *storage.State
	account       *storage.Account
	caller        crypto.Address
	gasPolicy     gas.Policy
	callDepth     int
	memAggr       int
	events        []event.Event
	methodLookup  map[string]*foreignMethod
	ptrArgSizeMap map[int]int
	gas           *vm.Gas
	parent        *Engine
}

// NewEngine return new instance of Engine
func NewEngine(state *storage.State, account *storage.Account, caller crypto.Address, gasPolicy gas.Policy, gasLimit uint64) *Engine {
	return &Engine{
		state:         state,
		account:       account,
		caller:        caller,
		gasPolicy:     gasPolicy,
		events:        []event.Event{},
		methodLookup:  make(map[string]*foreignMethod),
		ptrArgSizeMap: make(map[int]int),
		gas:           &vm.Gas{Limit: gasLimit},
		parent:        nil,
	}
}

// GetEvents return the event of engine
func (engine *Engine) GetEvents() []event.Event {
	return engine.events
}

// GetGasUsed return gas used by vm
func (engine *Engine) GetGasUsed() uint64 {
	return engine.gas.Used
}

// NewChildEngine share with parent state except caller is contract itself
func (engine *Engine) NewChildEngine(account *storage.Account) *Engine {
	return &Engine{
		account:       account,
		state:         engine.state,
		caller:        engine.account.GetAddress(),
		gasPolicy:     engine.gasPolicy,
		events:        []event.Event{},
		methodLookup:  make(map[string]*foreignMethod),
		ptrArgSizeMap: make(map[int]int),
		gas:           engine.gas,
		parent:        engine,
	}
}

// Ignite executes a contract given its code, method, and arguments
func (engine *Engine) Ignite(method string, methodArgs []byte) (uint64, error) {
	contract, err := engine.account.GetContract()
	if err != nil {
		return 0, err
	}
	vm, err := vertex.NewVM(contract.Code, engine.gasPolicy, engine.gas, engine)
	if err != nil {
		return 0, err
	}
	funcID, ok := vm.GetFunctionIndex(method)
	if !ok {
		return 0, errors.New("Cannot find invoke function")
	}

	val, _ := vm.Module.ExecInitExpr(vm.Module.GetGlobal(int(vm.Module.ExportSec.ExportMap[ExportSecDataEnd].Desc.Idx)).Init)
	offset := int(val.(int32))

	function, err := contract.Header.GetFunction(method)
	if err != nil {
		return 0, err
	}

	decodedBytes, err := abi.DecodeToBytes(function.Parameters, methodArgs)
	if err != nil {
		return 0, err
	}

	arguments, err := engine.loadArguments(vm, decodedBytes, function.Parameters, offset)
	if err != nil {
		return 0, err
	}
	ret, err := vm.Invoke(funcID, arguments...)
	return ret, err
}

func (engine *Engine) setStats(callDepth, memAggr int) {
	engine.callDepth = callDepth
	engine.memAggr = memAggr
}

func (engine *Engine) loadArguments(vm *vm.VM, byteArgs [][]byte, params []*abi.Parameter, offset int) ([]uint64, error) {
	var args = make([]uint64, len(byteArgs))
	byteSize := 0
	for _, bytes := range byteArgs {
		byteSize += len(bytes)
	}
	if byteSize > 1024 {
		return []uint64{}, fmt.Errorf("arguments byte size exceeds limit")
	}
	for i, bytes := range byteArgs {
		isArray := params[i].IsArray || params[i].Type.IsAddress()
		if isArray {
			if params[i].Type.IsAddress() {
				if _, err := crypto.AddressFromBytes(bytes); err != nil {
					return nil, err
				}
			}
			if _, err := vm.MemWrite(bytes, offset); err != nil {
				return nil, err
			}
			args[i] = uint64(offset)
			engine.ptrArgSizeMap[offset] = len(bytes)
			offset += len(bytes)
		} else {
			buffer := make([]byte, 8)
			copy(buffer, bytes)
			args[i] = binary.LittleEndian.Uint64(buffer)
		}
	}
	return args, nil
}

func (engine *Engine) ptrArgSizeGet(ptr int) (int, error) {
	size, ok := engine.ptrArgSizeMap[ptr]
	if !ok {
		return 0, errors.New("pointer size not found")
	}
	return size, nil
}

func (engine *Engine) pushEvent(event event.Event) {
	if engine.parent != nil {
		engine.parent.pushEvent(event)
	} else {
		engine.events = append(engine.events, event)
	}
}
