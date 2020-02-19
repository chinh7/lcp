package gas

import (
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/event"
)

const minimumGasPrice = uint64(18)

// LiquidStation provide a liquid as a gas station
type LiquidStation struct {
	app       App
	policy    Policy
	collector crypto.Address
}

// Sufficient gas of an address is enough for burn
func (station *LiquidStation) Sufficient(addr crypto.Address, fee uint64) bool {
	token := station.app.GetGasContractToken()
	balance, err := token.GetBalance(addr)
	if err != nil {
		panic(err)
	}
	return fee <= balance
}

// Burn gas
func (station *LiquidStation) Burn(addr crypto.Address, fee uint64) []event.Event {
	token := station.app.GetGasContractToken()
	// Move to gas owner
	if fee > 0 {
		events, err := token.Transfer(addr, station.collector, fee)
		if err != nil {
			panic(err)
		}
		return events
	}
	return nil
}

// Switch off fee, never call
func (station *LiquidStation) Switch() bool {
	return false
}

// GetPolicy for liquid token
func (station *LiquidStation) GetPolicy() Policy {
	return station.policy
}

// CheckGasPrice of transaction
func (station *LiquidStation) CheckGasPrice(price uint64) bool {
	return price >= minimumGasPrice
}

// NewLiquidStation with fee
func NewLiquidStation(app App, collector crypto.Address) Station {
	return &LiquidStation{
		app:       app,
		policy:    &AlphaPolicy{},
		collector: collector,
	}
}
