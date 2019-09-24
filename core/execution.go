package core

import (
	"log"
	"strconv"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/QuoineFinancial/vertex/trie"
	"github.com/QuoineFinancial/vertex/vm"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

var events []types.Event

// ApplyTx executes a transaction by either deploying the contract code or invoking a contract method call
func ApplyTx(appHash trie.Hash, tx *crypto.Tx) {
	state := storage.GetState(appHash)
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
		_, results := vm.Call(accountState, data.Method, data.Params...)
		parseEvents(results)
	}
}

// GetEvents returns events for the last VM execution
func GetEvents() []types.Event {
	return events
}

func parseEvents(results [][]byte) {
	attributes := []common.KVPair{}
	for index, result := range results {
		attributes = append(attributes, common.KVPair{
			Key: []byte(strconv.Itoa(index)), Value: result,
		})
	}
	event := types.Event{
		Type:       "result",
		Attributes: attributes,
	}
	// events = make([]types.Event, 0)
	events = append(events, event)
}
