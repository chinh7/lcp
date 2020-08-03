package storage

import (
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/storage"
)

// GetAccountParams is params to GetAccount transaction
type GetAccountParams struct {
	Address string `json:"address"`
}

// GetAccountResult is result of GetAccount
type GetAccountResult struct {
	Account *storage.Account `json:"account"`
}

// GetAccount delivers transaction to blockchain
func (service *Service) GetAccount(r *http.Request, params *GetAccountParams, result *GetAccountResult) error {
	// TODO: Add GetAccount API
	return nil
}
