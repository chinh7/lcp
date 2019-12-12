package storage

import (
	"fmt"
	"net/http"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/engine"
	"github.com/QuoineFinancial/vertex/gas"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/ethereum/go-ethereum/common"
)

// CallParams is params to execute Call
type CallParams struct {
	Height  int64    `json:"height"`
	Address string   `json:"address"`
	Method  string   `json:"method"`
	Args    []string `json:"args"`
}

// CallResult is result of Call
type CallResult struct {
	Return interface{} `json:"value"`
}

// Call to execute function without tx creation in blockchain
func (service *Service) Call(r *http.Request, params *CallParams, result *CallResult) error {
	var appHash common.Hash
	if params.Height == 0 {
		status, _ := service.tAPI.Status()
		appHash = common.BytesToHash(status.SyncInfo.LatestAppHash)
	} else {
		block, err := service.tAPI.Block(&params.Height)
		if err != nil {
			return fmt.Errorf("Block %d not found", params.Height)
		}
		appHash = common.BytesToHash(block.BlockMeta.Header.AppHash)
	}

	state, err := storage.New(appHash, service.database)
	if err != nil {
		return fmt.Errorf("Could not init state for app hash %s", appHash.String())
	}

	account, err := state.GetAccount(crypto.AddressFromString(params.Address))
	if err != nil {
		return fmt.Errorf("Account not found for address %s", params.Address)
	}

	contract, err := account.GetContract()
	if err != nil {
		return fmt.Errorf("Contract not found for address %s", params.Address)
	}

	function, err := contract.Header.GetFunction(params.Method)
	if err != nil {
		return fmt.Errorf("Function for method %s not found", params.Method)
	}

	data, err := abi.EncodeFromString(function.Parameters, params.Args)
	if err != nil {
		return fmt.Errorf("Invalid params for method %s", params.Method)
	}

	engine := engine.NewEngine(state, account, crypto.AddressFromString(params.Address), &gas.FreePolicy{}, 0)
	ret, err := engine.Ignite(params.Method, data)
	if err != nil {
		return err
	}

	result.Return = ret
	return nil
}
