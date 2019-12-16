package core

import (
	"errors"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/engine"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

// ApplyTx executes a transaction by either deploying the contract code or invoking a contract method call
func ApplyTx(state *storage.State, tx *crypto.Tx, gasStation gas.Station) (uint64, []event.Event, uint64, error) {
	policy := gasStation.GetPolicy()
	gasLimit := tx.GasLimit
	if (tx.To == crypto.Address{}) {
		contractSize := len(tx.Data)
		gasUsed := policy.GetCostForContract(contractSize)
		if gasLimit < gasUsed {
			return 0, nil, 0, errors.New("out of gas")
		}
		contractAddress := tx.From.CreateAddress()
		_, err := state.CreateAccount(tx.From.Address(), contractAddress, tx.Data)
		if err != nil {
			return 0, nil, 0, err
		}
		gasEvents := gasStation.Burn(tx.From.Address(), gasUsed*tx.GasPrice)
		events := append([]event.Event{event.NewDeploymentEvent(contractAddress)}, gasEvents...)
		state.Commit()
		return 0, events, gasUsed, nil
	}
	data := &crypto.TxData{}
	_ = data.Deserialize(tx.Data) // deserialize error is already checked in checkTx
	contractAccount, err := state.GetAccount(tx.To)
	if err != nil {
		return 0, nil, 0, err
	}

	// Create new account if fromAddress is not exist
	fromAddress := tx.From.Address()
	fromAccount, _ := state.GetAccount(fromAddress)
	if fromAccount == nil {
		fromAccount, err = state.CreateAccount(fromAddress, fromAddress, nil)
		if err != nil {
			return 0, nil, 0, err
		}
	}

	nonce := fromAccount.Nonce
	fromAccount.SetNonce(nonce + 1)

	execEngine := engine.NewEngine(state, contractAccount, tx.From.Address(), policy, gasLimit)
	result, err := execEngine.Ignite(data.Method, data.Params)
	engineEvents := []event.Event{}
	if err != nil {
		state.Revert()
	} else {
		engineEvents = execEngine.GetEvents()
	}

	gasUsed := execEngine.GetGasUsed()
	gasEvents := gasStation.Burn(tx.From.Address(), gasUsed)
	events := append(engineEvents, gasEvents...)
	state.Commit()
	return result, events, gasUsed, err
}
