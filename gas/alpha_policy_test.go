package gas

import (
	"github.com/vertexdlt/vertexvm/opcode"
	"testing"
)

func TestAlphaPolicy(t *testing.T) {
	policy := AlphaPolicy{}
	cost := policy.GetCostForOp(opcode.Select)
	if cost != 5 {
		t.Errorf("Expect cost %v, got %v", 5, cost)
	}
	cost = policy.GetCostForStorage(100)
	if cost != 100 {
		t.Errorf("Expect cost %v, got %v", 0, cost)
	}
	cost = policy.GetCostForContract(100)
	if cost != 100 {
		t.Errorf("Expect cost %v, got %v", 0, cost)
	}
}
