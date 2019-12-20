package gas

import (
	"github.com/vertexdlt/vertexvm/opcode"
	"testing"
)

func TestFreePolicy(t *testing.T) {
	policy := FreePolicy{}
	cost := policy.GetCostForOp(opcode.Select)
	if cost != 0 {
		t.Errorf("Expect cost %v, got %v", 0, cost)
	}
	cost = policy.GetCostForStorage(100)
	if cost != 0 {
		t.Errorf("Expect cost %v, got %v", 0, cost)
	}
	cost = policy.GetCostForContract(100)
	if cost != 0 {
		t.Errorf("Expect cost %v, got %v", 0, cost)
	}
}
