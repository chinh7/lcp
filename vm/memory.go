package vm

import (
	"log"

	"github.com/go-interpreter/wagon/wasm"
)

const dataEndVal = "__data_end"
const heapBaseVal = "__heap_base"
const stackSize = 16384 // 16Kb
var bufferConfig BufferConfig

// BufferConfig manages the layout of a secondary heap layout used for data sharing between host and wasm
type BufferConfig struct {
	base    int32
	size    int32
	index   int32
	sizeMap map[int32]int32
}

// NewBufferConfig initializes a BufferConfig object
func NewBufferConfig(dataEnd int32, heapBase int32) BufferConfig {
	config := BufferConfig{}
	config.base = dataEnd
	config.index = config.base
	config.size = heapBase - stackSize - config.base
	config.sizeMap = make(map[int32]int32)
	return config
}

func malloc(size int32) (memPtr int32) {
	if bufferConfig.index+size >= bufferConfig.base+bufferConfig.size {
		log.Fatalf("Buffer memory exceeded")
	}
	memPtr = bufferConfig.index
	bufferConfig.index += size
	bufferConfig.sizeMap[memPtr] = size
	return
}

func initMemory(m *wasm.Module) {
	val, err := m.ExecInitExpr(m.GetGlobal(int(m.Export.Entries[dataEndVal].Index)).Init)
	if err != nil {
		log.Fatalf("Could not read data end: %v", err)
	}
	dataEnd := val.(int32)
	val, err = m.ExecInitExpr(m.GetGlobal(int(m.Export.Entries[heapBaseVal].Index)).Init)
	if err != nil {
		log.Fatalf("Could not read heap base: %v", err)
	}
	heapBase := val.(int32)

	bufferConfig = NewBufferConfig(dataEnd, heapBase)
}
