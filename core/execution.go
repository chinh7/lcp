package core

import (
	"log"

	"github.com/vertexdlt/vertex/crypto"
	"github.com/vertexdlt/vertex/storage"
	"github.com/vertexdlt/vertex/vm"
)

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
		vm.Call(accountState, data.Method, data.Params...)
	}
}
