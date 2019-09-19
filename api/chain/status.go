package chain

import (
	"net/http"
)

// GetStatusResult is response of StatusService
type GetStatusResult struct {
	ChainID           string `json:"chainId"`
	LatestBlockHash   string `json:"latestBlockHash"`
	LatestBlockHeight int64  `json:"latestBlockHeight"`
}

// GetStatus returns current status of chain
func (service *Service) GetStatus(r *http.Request, _ interface{}, result *GetStatusResult) error {
	status, err := service.tAPI.Status()
	if err != nil {
		return err
	}
	result.LatestBlockHash = status.SyncInfo.LatestBlockHash.String()
	result.ChainID = status.NodeInfo.Network
	result.LatestBlockHeight = status.SyncInfo.LatestBlockHeight
	return nil
}
