package vm

import (
	"log"
	"reflect"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/tendermint/tendermint/libs/common"
)

const wasmPageSize = 65536 // 64Kb

func syscall(proc *exec.Process, idx int32, args ...int32) int32 {
	if idx == 45 {
		return 0
	} else if idx == 192 {
		requested := args[1]
		log.Println("mmap2", requested)
		return int32(0)
	} else {
		log.Printf("syscall %d: NYI\n", idx)
	}
	return -1
}

func readAt(proc *exec.Process, ptr int32) []byte {
	size := bufferConfig.sizeMap[ptr]
	if size == 0 {
		size = 16
	}
	data := make([]byte, size)
	proc.ReadAt(data, int64(ptr))
	return data
}

func syscall0(proc *exec.Process, idx int32) int32 {
	return syscall(proc, idx)
}

func syscall1(proc *exec.Process, idx, a int32) int32 {
	return syscall(proc, idx, a)
}

func syscall2(proc *exec.Process, idx, a, b int32) int32 {
	return syscall(proc, idx, a, b)
}

func syscall3(proc *exec.Process, idx, a, b, c int32) int32 {
	return syscall(proc, idx, a, b, c)
}

func syscall4(proc *exec.Process, idx, a, b, c, d int32) int32 {
	return syscall(proc, idx, a, b, c, d)
}

func syscall5(proc *exec.Process, idx, a, b, c, d, e int32) int32 {
	return syscall(proc, idx, a, b, c, d, e)
}

func syscall6(proc *exec.Process, idx, a, b, c, d, e, f int32) int32 {
	return syscall(proc, idx, a, b, c, d, e, f)
}

func printBytes(proc *exec.Process, size int32, ptr int32) {
	key := readAt(proc, ptr)
	log.Println("printBytes", string(key))
}

func (vertexVM *VertexVM) emitEvent(proc *exec.Process, ptr int32) {
	data := readAt(proc, ptr)
	vertexVM.event.Attributes = append(vertexVM.event.Attributes, common.KVPair{
		Key:   []byte(""), // TODO: Decide the key of event
		Value: data,
	})
}

func (vertexVM *VertexVM) getStorage(proc *exec.Process, keyPtr int32) (valuePtr int32) {
	key := readAt(proc, keyPtr)
	value, err := vertexVM.account.GetStorage(key)
	if err == nil && len(value) > 0 {
		valuePtr = malloc(int32(len(value)))
		proc.WriteAt(value, int64(valuePtr))
	} else {
		valuePtr = 0
	}
	return
}

func (vertexVM *VertexVM) setStorage(proc *exec.Process, keyPtr int32, valuePtr int32) {
	key := readAt(proc, keyPtr)
	value := readAt(proc, valuePtr)
	vertexVM.account.SetStorage(key, value)
}

func (vertexVM *VertexVM) resolveImports(name string) (*wasm.Module, error) {
	m := wasm.NewModule()

	m.Types = &wasm.SectionTypes{
		// All function types in this module
		Entries: []wasm.FunctionSig{
			{
				Form:        0,
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:        0,
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:        0,
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:        0,
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:        0,
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:        0,
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:        0,
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:        0,
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{},
			},
			{
				Form:        0,
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{},
			},
		},
	}
	m.FunctionIndexSpace = []wasm.Function{
		{
			Sig:  &m.Types.Entries[0],
			Host: reflect.ValueOf(syscall0),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[1],
			Host: reflect.ValueOf(syscall1),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[2],
			Host: reflect.ValueOf(syscall2),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[3],
			Host: reflect.ValueOf(syscall3),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[4],
			Host: reflect.ValueOf(syscall4),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[5],
			Host: reflect.ValueOf(syscall5),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[6],
			Host: reflect.ValueOf(syscall6),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[0],
			Host: reflect.ValueOf(vertexVM.getStorage),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[7],
			Host: reflect.ValueOf(vertexVM.setStorage),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[7],
			Host: reflect.ValueOf(printBytes),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[8],
			Host: reflect.ValueOf(vertexVM.emitEvent),
			Body: &wasm.FunctionBody{},
		},
	}

	m.Export = &wasm.SectionExports{
		Entries: map[string]wasm.ExportEntry{
			"__syscall0": {
				FieldStr: "__syscall0",
				Kind:     wasm.ExternalFunction,
				Index:    0,
			},
			"__syscall1": {
				FieldStr: "__syscall1",
				Kind:     wasm.ExternalFunction,
				Index:    1,
			},
			"__syscall2": {
				FieldStr: "__syscall2",
				Kind:     wasm.ExternalFunction,
				Index:    2,
			},
			"__syscall3": {
				FieldStr: "__syscall3",
				Kind:     wasm.ExternalFunction,
				Index:    3,
			},
			"__syscall4": {
				FieldStr: "__syscall4",
				Kind:     wasm.ExternalFunction,
				Index:    4,
			},
			"__syscall5": {
				FieldStr: "__syscall5",
				Kind:     wasm.ExternalFunction,
				Index:    5,
			},
			"__syscall6": {
				FieldStr: "__syscall6",
				Kind:     wasm.ExternalFunction,
				Index:    6,
			},
			"get_storage": {
				FieldStr: "get_storage",
				Kind:     wasm.ExternalFunction,
				Index:    7,
			},
			"set_storage": {
				FieldStr: "set_storage",
				Kind:     wasm.ExternalFunction,
				Index:    8,
			},
			"print_bytes": {
				FieldStr: "print_bytes",
				Kind:     wasm.ExternalFunction,
				Index:    9,
			},
			"emit_event": {
				FieldStr: "emit_event",
				Kind:     wasm.ExternalFunction,
				Index:    10,
			},
		},
	}
	return m, nil
}
