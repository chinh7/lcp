package consensus

import (
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/engine"
)

const (
	// InitFunctionName is default init function name
	InitFunctionName = "init"
)

func (app *App) applyTransaction(tx *crypto.Transaction) (*crypto.TxReceipt, error) {
	if tx.Receiver == crypto.EmptyAddress {
		return app.deployContract(tx)
	}
	return app.invokeContract(tx)
}

func (app *App) deployContract(tx *crypto.Transaction) (*crypto.TxReceipt, error) {
	var receipt crypto.TxReceipt

	contractSize := len(tx.Payload.Contract)
	policy := app.gasStation.GetPolicy()
	receipt.GasUsed = uint32(policy.GetCostForContract(contractSize))
	if tx.GasLimit < receipt.GasUsed {
		receipt.Error = "Out of gas"
		receipt.Success = false
		return &receipt, nil
	}

	// Create contract account
	senderAddress := crypto.AddressFromPubKey(tx.Sender.PublicKey)
	contractAddress := crypto.NewDeploymentAddress(senderAddress, tx.Sender.Nonce)
	contractAccount, err := app.state.CreateAccount(senderAddress, contractAddress, tx.Payload.Contract)

	if err != nil {
		return nil, err
	}

	if tx.Payload.Method == InitFunctionName {
		execEngine := engine.NewEngine(app.state, contractAccount, senderAddress, policy, uint64(tx.GasLimit-receipt.GasUsed))
		if result, err := execEngine.Ignite(tx.Payload.Method, tx.Payload.Params); err != nil {
			receipt.Error = err.Error()
			receipt.Success = false
			app.state.Revert()
		} else {
			receipt.Result = result
			receipt.Success = true
		}
		receipt.GasUsed += uint32(execEngine.GetGasUsed())
	}

	// Create account for creator and increase nonce by 1
	if err := app.increaseNonce(senderAddress); err != nil {
		return nil, err
	}

	gasEvents := app.gasStation.Burn(senderAddress, uint64(receipt.GasUsed)*uint64(tx.GasPrice))
	receipt.Events = append(receipt.Events, gasEvents...)
	return &receipt, nil
}

func (app *App) invokeContract(tx *crypto.Transaction) (*crypto.TxReceipt, error) {
	var receipt crypto.TxReceipt

	// contract account not found is checked before apply tx
	contractAccount, err := app.state.GetAccount(tx.Receiver)
	if err != nil {
		panic(err)
	}

	policy := app.gasStation.GetPolicy()
	senderAddress := crypto.AddressFromPubKey(tx.Sender.PublicKey)
	execEngine := engine.NewEngine(app.state, contractAccount, senderAddress, policy, uint64(tx.GasLimit))

	if result, err := execEngine.Ignite(tx.Payload.Method, tx.Payload.Params); err != nil {
		receipt.Error = err.Error()
		receipt.Success = false
		app.state.Revert()
	} else {
		receipt.Result = result
		receipt.Events = append(receipt.Events, execEngine.GetEvents()...)
	}

	// Create/get account for creator and increase nonce by 1
	if err := app.increaseNonce(senderAddress); err != nil {
		return nil, err
	}

	receipt.GasUsed = uint32(execEngine.GetGasUsed())
	gasEvents := app.gasStation.Burn(senderAddress, uint64(receipt.GasUsed)*uint64(tx.GasPrice))
	receipt.Events = append(receipt.Events, gasEvents...)

	return &receipt, nil
}

func (app *App) increaseNonce(address crypto.Address) error {
	account, err := app.state.GetAccount(address)
	if err != nil {
		return err
	}

	// Make sure account is created
	if account == nil {
		account, err = app.state.CreateAccount(address, address, nil)
		if err != nil {
			return err
		}
	}

	account.SetNonce(account.Nonce + 1)
	return nil
}
