package consensus

import (
	"fmt"

	"github.com/QuoineFinancial/liquid-chain/crypto"
)

func (app *App) validateTx(tx *crypto.Transaction) (uint32, error) {
	nonce := uint64(0)
	address := crypto.AddressFromPubKey(tx.Sender.PublicKey)
	account, err := app.state.GetAccount(address)
	if err != nil {
		return CodeTypeUnknownError, err
	}
	if account != nil {
		nonce = account.Nonce
	}

	// Validate tx nonce
	if tx.Sender.Nonce != nonce {
		err := fmt.Errorf("Invalid nonce. Expected %v, got %v", nonce, tx.Sender.Nonce)
		return CodeTypeBadNonce, err
	}

	// Validate tx signature
	signingHash := crypto.GetSigHash(tx)
	if valid := crypto.VerifySignature(tx.Sender.PublicKey, signingHash[:], tx.Signature); !valid {
		return CodeTypeInvalidSignature, fmt.Errorf("Invalid signature")
	}

	// Validate Non-existent contract invoke
	if tx.Receiver != crypto.EmptyAddress {
		// invoke transaction
		account, err := app.state.GetAccount(tx.Receiver)
		if err != nil {
			return CodeTypeUnknownError, err
		}
		if !account.IsContract() {
			return CodeTypeNonContractAccount, fmt.Errorf("Invoke a non-contract account")
		}
	}

	// Validate gas limit
	fee := uint64(tx.GasLimit) * uint64(tx.GasPrice)
	if !app.gasStation.Sufficient(address, fee) {
		return CodeTypeInsufficientFee, fmt.Errorf("Insufficient fee")
	}

	// Validate gas price
	if !app.gasStation.CheckGasPrice(tx.GasPrice) {
		return CodeTypeInvalidGasPrice, fmt.Errorf("Invalid gas price")
	}

	return CodeTypeOK, nil
}
