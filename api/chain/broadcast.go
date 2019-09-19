package chain

import (
	"encoding/hex"
	"net/http"
)

// BroadcastParams is params to broadcast transaction
type BroadcastParams struct {
	Transaction string
}

// BroadcastResult is result of broadcast
type BroadcastResult struct {
	TransactionHash string
}

// Broadcast delivers transction to blockchain
func (service *Service) Broadcast(
	r *http.Request,
	params *BroadcastParams,
	result *BroadcastResult,
) error {
	bytes, err := hex.DecodeString(params.Transaction)
	if err != nil {
		return err
	}
	broadcastResult, err := service.tAPI.BroadcastTxSync(bytes)
	if err != nil {
		return err
	}
	result.TransactionHash = broadcastResult.Hash.String()
	return nil
}
