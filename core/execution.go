package core

import (
	"strconv"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/engine"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

// ApplyTx executes a transaction by either deploying the contract code or invoking a contract method call
func ApplyTx(state *storage.State, tx *crypto.Tx) ([]types.Event, error) {
	if (tx.To == crypto.Address{}) {
		contractAddress := tx.From.CreateAddress()
		state.CreateAccount(tx.From.Address(), contractAddress, &tx.Data)
		event := types.Event{
			Type: "deploy",
			Attributes: []common.KVPair{
				common.KVPair{
					Key:   []byte("address"),
					Value: []byte(contractAddress.String()),
				},
			},
		}
		return []types.Event{event}, nil
	}
	data := &crypto.TxData{}
	data.Deserialize(tx.Data)
	execEngine := engine.NewEngine(state, state.GetAccount(tx.To), tx.From.Address())
	_, err := execEngine.Ignite(data.Method, data.Params)
	if err != nil {
		return nil, err
	}
	return execEngine.GetEvents(), nil
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
