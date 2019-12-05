package gas

import (
	"github.com/vertexdlt/vertexvm/vm"
)

// Policy for gas cost
type Policy interface {
	vm.GasPolicy
	GetCostForStorage(size int) int64
	GetCostForContract(size int) int64
}
