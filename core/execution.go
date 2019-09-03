package core

import (
	"log"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/QuoineFinancial/vertex/vm"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

var events []types.Event

// ApplyTx executes a transaction by either deploying the contract code or invoking a contract method call
func ApplyTx(tx *crypto.Tx) {
	state := storage.GetState()
	createContract := tx.To == crypto.Address{}
	if createContract {
		contractAddress := tx.From.CreateAddress()
		log.Println("Deploy contract", contractAddress)
		accountState := state.CreateAccountState(contractAddress)
		accountState.SetCode(tx.Data)
	} else {
		log.Println("Invoke contract", tx.To)
		accountState := state.GetAccountState(tx.To)
		data := &crypto.TxData{}
		data.Deserialize(tx.Data)
		_, eventBytes := vm.Call(accountState, data.Method, data.Params...)
		parseEvents(eventBytes)
	}
}

// GetEvents returns events for the last VM execution
func GetEvents() []types.Event {
	return events
}

func parseEvents(eventBytes [][]byte) {
	events = make([]types.Event, 0)
	for _, bytes := range eventBytes {
		attributes := []common.KVPair{common.KVPair{Key: make([]byte, 0), Value: bytes}}
		event := types.Event{Attributes: attributes}
		events = append(events, event)
	}
}
