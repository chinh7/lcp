package gas

import (
	"github.com/vertexdlt/vertexvm/opcode"
)

// FreePolicy is a simple policy for first version
type FreePolicy struct {
	Policy
}

// GetCostForOp get cost from table
func (p *FreePolicy) GetCostForOp(op opcode.Opcode) int64 {
	return 0
}

// GetCostForStorage size of data
func (p *FreePolicy) GetCostForStorage(size int) uint64 {
	return 0
}

// GetCostForContract creation
func (p *FreePolicy) GetCostForContract(size int) uint64 {
	return 0
}
