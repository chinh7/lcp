package core

import (
	"errors"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/engine"
	"github.com/QuoineFinancial/vertex/gas"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

// ApplyTx executes a transaction by either deploying the contract code or invoking a contract method call
func ApplyTx(state *storage.State, tx *crypto.Tx, gasStation gas.Station) ([]types.Event, uint64, error) {
	policy := gasStation.GetPolicy()
	gasLimit := int64(tx.GasLimit)
	if (tx.To == crypto.Address{}) {
		contractSize := len(tx.Data)

		gasUsed := policy.GetCostForContract(contractSize)
		if !gasStation.Sufficient(tx.From.Address(), gasUsed) {
			return nil, 0, errors.New("out of gas")
		}
		contractAddress := tx.From.CreateAddress()
		state.CreateAccount(tx.From.Address(), contractAddress, tx.Data)
		event := types.Event{
			Type: "deploy",
			Attributes: []common.KVPair{
				common.KVPair{
					Key:   []byte("address"),
					Value: contractAddress[:],
				},
			},
		}
		return []types.Event{event}, gasUsed, nil
	}
	data := &crypto.TxData{}
	data.Deserialize(tx.Data)
	contractAccount, err := state.GetAccount(tx.To)
	if err != nil {
		return nil, 0, err
	}
	execEngine := engine.NewEngine(state, contractAccount, tx.From.Address(), policy, gasLimit)
	_, gasUsed, err := execEngine.Ignite(data.Method, data.Params)
	gasStation.Burn(tx.From.Address(), gasUsed)
	if err != nil {
		return nil, gasUsed, err
	}
	return execEngine.GetEvents(), gasUsed, nil
}
