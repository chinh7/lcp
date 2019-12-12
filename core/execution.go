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
	gasLimit := tx.GasLimit
	if !gasStation.Sufficient(tx.From.Address(), gasLimit) {
		return nil, 0, errors.New("out of gas")
	}
	if (tx.To == crypto.Address{}) {
		contractSize := len(tx.Data)
		gasUsed := policy.GetCostForContract(contractSize)
		if gasLimit < gasUsed {
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
		gasEvents := gasStation.Burn(tx.From.Address(), gasUsed)
		events := append([]types.Event{event}, gasEvents...)
		state.Commit()
		return events, gasUsed, nil
	}
	data := &crypto.TxData{}
	data.Deserialize(tx.Data)
	contractAccount, err := state.GetAccount(tx.To)
	if err != nil {
		return nil, 0, err
	}

	execEngine := engine.NewEngine(state, contractAccount, tx.From.Address(), policy, gasLimit)
	if _, err = execEngine.Ignite(data.Method, data.Params); err != nil {
		state.Revert()
	}

	gasUsed := execEngine.GetGasUsed()
	gasEvents := gasStation.Burn(tx.From.Address(), gasUsed)
	events := append(execEngine.GetEvents(), gasEvents...)
	state.Commit()
	return events, gasUsed, err
}
