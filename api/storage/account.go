package storage

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/api/models"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/ethereum/go-ethereum/common"
)

// GetAccountParams is params to GetAccount transaction
type GetAccountParams struct {
	Address string `json:"address"`
}

// GetAccountResult is result of GetAccount
type GetAccountResult struct {
	Account *models.Account `json:"account"`
}

// GetAccount delivers transaction to blockchain
func (service *Service) GetAccount(r *http.Request, params *GetAccountParams, result *GetAccountResult) error {
	status, _ := service.tAPI.Status()
	appHash := common.BytesToHash(status.SyncInfo.LatestAppHash)
	state, err := storage.New(appHash, service.database)
	if err != nil {
		return err
	}

	address, err := crypto.AddressFromString(params.Address)
	if err != nil {
		return err
	}

	account, err := state.GetAccount(address)
	if err != nil {
		return err
	}
	if account == nil {
		return errors.New("Account not exist")
	}

	var contract *abi.Contract
	if len(account.ContractHash) > 0 {
		contract, err = account.GetContract()
		if err != nil {
			return err
		}
	}

	result.Account = &models.Account{
		Nonce:        account.Nonce,
		ContractHash: hex.EncodeToString(account.ContractHash),
		Contract:     contract,
	}
	return nil
}
