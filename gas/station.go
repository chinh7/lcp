package gas

import (
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/tendermint/tendermint/abci/types"
)

// Station interface for check and burn gas
type Station interface {
	Sufficient(addr crypto.Address, gas uint64) bool
	Burn(addr crypto.Address, gas uint64) []types.Event
	Switch() bool
	GetPolicy() Policy
}

// Token interface
type Token interface {
	GetBalance(addr crypto.Address) (uint64, error)
	Transfer(caller crypto.Address, addr crypto.Address, amount uint64) ([]types.Event, error)
	GetContract() *storage.Account
}

// App interface
type App interface {
	SetGasStation(gasStation Station)
	GetGasContractToken() Token
}
