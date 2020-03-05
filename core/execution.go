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
	var err error
	fromAddress := tx.From.Address()
	fromAccount, _ := state.GetAccount(fromAddress)
	// Make sure fromAccount is created
	if fromAccount == nil {
		fromAccount, err = state.CreateAccount(fromAddress, fromAddress, nil)
		if err != nil {
			return 0, nil, 0, err
		}
	}
	fromAccount.SetNonce(fromAccount.Nonce + 1)
	state.Commit()

	if (tx.To == crypto.Address{}) {
		return applyDeployContractTx(state, tx, gasStation)
	}
	return applyInvokeTx(state, tx, gasStation)
}

func applyDeployContractTx(state *storage.State, tx *crypto.Tx, gasStation gas.Station) (uint64, []event.Event, uint64, error) {
	contractSize := len(tx.Data)
	policy := gasStation.GetPolicy()
	gasUsed := policy.GetCostForContract(contractSize)
	if uint64(tx.GasLimit) < gasUsed {
		return 0, nil, 0, errors.New("out of gas")
	}

	// Create contract account
	contractAddress := tx.From.CreateAddress()
	_, err := state.CreateAccount(tx.From.Address(), contractAddress, tx.Data)
	if err != nil {
		return 0, nil, 0, err
	}

	gasEvents := gasStation.Burn(tx.From.Address(), gasUsed*uint64(tx.GasPrice))
	events := append([]event.Event{event.NewDeploymentEvent(contractAddress)}, gasEvents...)
	state.Commit()
	return 0, events, gasUsed, nil
}

func applyInvokeTx(state *storage.State, tx *crypto.Tx, gasStation gas.Station) (uint64, []event.Event, uint64, error) {
	var err error
	// contract account not found is checked before apply tx
	contractAccount, err := state.GetAccount(tx.To)
	if err != nil {
		return 0, nil, 0, err
	}

	policy := gasStation.GetPolicy()
	execEngine := engine.NewEngine(state, contractAccount, tx.From.Address(), policy, uint64(tx.GasLimit))

	data := &crypto.TxData{}
	_ = data.Deserialize(tx.Data) // deserialize error is already checked in checkTx
	result, err := execEngine.Ignite(data.Method, data.Params)
	engineEvents := []event.Event{}
	if err != nil {
		state.Revert()
	} else {
		engineEvents = execEngine.GetEvents()
	}

	gasUsed := execEngine.GetGasUsed()
	gasEvents := gasStation.Burn(tx.From.Address(), gasUsed*uint64(tx.GasPrice))
	events := append(engineEvents, gasEvents...)
	state.Commit()
	return result, events, gasUsed, err
}
