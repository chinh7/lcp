package chain

import (
	"encoding/base64"
	"net/http"

	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

// BroadcastParams is params to broadcast transaction
type BroadcastParams struct {
	RawTransaction string `json:"rawTx"`
}

// BroadcastResult is result of broadcast
type BroadcastResult struct {
	TransactionHash string `json:"hash"`
	Code            uint32 `json:"code"`
	Log             string `json:"log"`
}

func formatBroadcastResponse(result *BroadcastResult, res *core_types.ResultBroadcastTx) {
	result.TransactionHash = res.Hash.String()
	result.Code = res.Code
	result.Log = res.Log
}

// Broadcast delivers transaction to blockchain
func (service *Service) Broadcast(r *http.Request, params *BroadcastParams, result *BroadcastResult) error {
	bytes, err := base64.StdEncoding.DecodeString(params.RawTransaction)
	if err != nil {
		return err
	}
	broadcastResult, err := service.tAPI.BroadcastTxSync(bytes)
	if err != nil {
		return err
	}

	formatBroadcastResponse(result, broadcastResult)
	return nil
}

// BroadcastAsync broadcast but wont wait for transaction's checkTx result
func (service *Service) BroadcastAsync(r *http.Request, params *BroadcastParams, result *BroadcastResult) error {
	bytes, err := base64.StdEncoding.DecodeString(params.RawTransaction)
	if err != nil {
		return err
	}
	broadcastResult, err := service.tAPI.BroadcastTxAsync(bytes)
	if err != nil {
		return err
	}
	formatBroadcastResponse(result, broadcastResult)
	return nil
}

// BroadcastCommit broadcast and wait until the transaction is committed in a block or fail to pass checkTx
func (service *Service) BroadcastCommit(r *http.Request, params *BroadcastParams, result *BroadcastResult) error {
	bytes, err := base64.StdEncoding.DecodeString(params.RawTransaction)
	if err != nil {
		return err
	}
	broadcastResult, err := service.tAPI.BroadcastTxCommit(bytes)
	if err != nil {
		return err
	}

	result.TransactionHash = broadcastResult.Hash.String()

	// Handle if tx is rejected by CheckTx
	if broadcastResult.CheckTx.IsErr() {
		result.Code = broadcastResult.CheckTx.Code
		result.Log = broadcastResult.CheckTx.Log
		return nil
	}

	result.Code = broadcastResult.DeliverTx.Code
	result.Log = broadcastResult.DeliverTx.Log
	return nil
}
