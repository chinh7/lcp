package gas

import (
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

// Station interface for check and burn gas
type Station interface {
	Sufficient(addr crypto.Address, gas uint64) bool
	Burn(addr crypto.Address, gas uint64) []event.Event
	CheckGasPrice(price uint32) bool
	Switch() bool
	GetPolicy() Policy
}

// Token interface
type Token interface {
	GetBalance(addr crypto.Address) (uint64, error)
	Transfer(caller crypto.Address, addr crypto.Address, amount uint64) ([]event.Event, error)
	GetContract() *storage.Account
}

// App interface
type App interface {
	SetGasStation(gasStation Station)
	GetGasContractToken() Token
}
