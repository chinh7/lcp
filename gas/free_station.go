package gas

import (
	"log"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/event"
)

// FreeStation provide a free gas station
type FreeStation struct {
	app    App
	policy Policy
}

// Sufficient gas of an address is enough for burn
func (station *FreeStation) Sufficient(addr crypto.Address, gas uint64) bool {
	return true
}

// Burn gas, do nothing
func (station *FreeStation) Burn(addr crypto.Address, gas uint64) []event.Event {
	return nil
}

// Switch on fee
func (station *FreeStation) Switch() bool {
	app := station.app
	token := app.GetGasContractToken()
	if token != nil {
		contract := token.GetContract()
		creator := contract.GetCreator()
		balance, err := token.GetBalance(creator)
		if err != nil {
			panic(err)
		}
		// Only activate if creator balance > 0 aka minted
		if balance > 0 {
			log.Println("Change to liquid station")
			app.SetGasStation(NewLiquidStation(app, contract.GetAddress()))
			return true
		}
	}
	return false
}

// GetPolicy free
func (station *FreeStation) GetPolicy() Policy {
	return station.policy
}

// NewFreeStation constructor
func NewFreeStation(app App) Station {
	return &FreeStation{
		app:    app,
		policy: &FreePolicy{},
	}
}
