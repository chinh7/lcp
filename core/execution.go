package core

import (
	"log"
	"strconv"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/QuoineFinancial/vertex/vm"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

// ApplyTx executes a transaction by either deploying the contract code or invoking a contract method call
func ApplyTx(state *storage.State, tx *crypto.Tx) (types.Event, error) {
	event := types.Event{}
	createContract := tx.To == crypto.Address{}
	if createContract {
		contractAddress := tx.From.CreateAddress()
		log.Println("Deploy contract", contractAddress)
		log.Println(tx.Data)
		state.CreateAccount(contractAddress, &tx.Data)
		event = types.Event{
			Type: "result",
			Attributes: []common.KVPair{
				common.KVPair{
					Key:   []byte("address"),
					Value: []byte(contractAddress.String()),
				},
			},
		}
	} else {
		log.Println("Invoke contract", tx.To)
		data := &crypto.TxData{}
		data.Deserialize(tx.Data)
		vertexVM := vm.NewVertexVM(state.GetAccount(tx.To))
		_, results, err := vertexVM.Call(data.Method, data.Params...)
		event := parseEvent(results)
		if err != nil {
			return event, err
		}
	}
	return event, nil
}

func parseEvent(results [][]byte) types.Event {
	attributes := []common.KVPair{}
	for index, result := range results {
		attributes = append(attributes, common.KVPair{
			Key: []byte(strconv.Itoa(index)), Value: result,
		})
	}
	return types.Event{
		Type:       "result",
		Attributes: attributes,
	}
}
