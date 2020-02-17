package chain

import (
	"encoding/base64"
	"net/http"
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

// Broadcast delivers transction to blockchain
func (service *Service) Broadcast(r *http.Request, params *BroadcastParams, result *BroadcastResult) error {
	bytes, err := base64.StdEncoding.DecodeString(params.RawTransaction)
	if err != nil {
		return err
	}
	broadcastResult, err := service.tAPI.BroadcastTxSync(bytes)
	if err != nil {
		return err
	}
	result.TransactionHash = broadcastResult.Hash.String()
	result.Code = broadcastResult.Code
	result.Log = broadcastResult.Log
	return nil
}
