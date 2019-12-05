package gas

import (
	"github.com/QuoineFinancial/vertex/crypto"
)

// LiquidStation provide a liquid as a gas station
type LiquidStation struct {
	app       App
	policy    Policy
	collector crypto.Address
}

// Sufficient gas of an address is enough for burn
func (station *LiquidStation) Sufficient(addr crypto.Address, gas int64) bool {
	token := station.app.GetGasContractToken()
	balance, err := token.GetBalance(addr)
	if err != nil {
		panic(err)
	}
	return uint64(gas) <= balance
}

// Burn gas
func (station *LiquidStation) Burn(addr crypto.Address, gas int64) {
	token := station.app.GetGasContractToken()
	// Move to gas owner
	if gas > 0 {
		token.Transfer(addr, station.collector, uint64(gas))
	}
}

// Switch off fee, never call
func (station *LiquidStation) Switch() bool {
	return false
}

// GetPolicy for liquid token
func (station *LiquidStation) GetPolicy() Policy {
	return station.policy
}

// NewLiquidStation with fee
func NewLiquidStation(app App, collector crypto.Address) Station {
	return &LiquidStation{
		app:       app,
		policy:    &AlphaPolicy{},
		collector: collector,
	}
}
