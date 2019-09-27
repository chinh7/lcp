package storage

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/ethereum/go-ethereum/common"
)

// GetAccountParams is params to GetAccount transaction
type GetAccountParams struct {
	Address string `json:"address"`
}

// GetAccountResult is result of GetAccount
type GetAccountResult struct {
	Nonce    uint64 `json:"nonce"`
	CodeHash string `json:"codeHash"`
}

// GetAccount delivers transction to blockchain
func (service *Service) GetAccount(r *http.Request, params *GetAccountParams, result *GetAccountResult) error {
	status, _ := service.tAPI.Status()
	appHash := common.BytesToHash(status.SyncInfo.LatestAppHash)
	state := storage.GetState(appHash)
	account := state.GetAccount(crypto.AddressFromString(params.Address))
	fmt.Println("ACCOUNT", account)
	result = &GetAccountResult{
		Nonce:    account.Nonce,
		CodeHash: hex.EncodeToString(account.CodeHash),
	}

	return nil
}
