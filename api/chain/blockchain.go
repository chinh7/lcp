package chain

import (
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/api/models"
)

// BlockchainParams is params of GetBlockchain
type BlockchainParams struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

// BlockchainResult is response of GetBlockchain
type BlockchainResult struct {
	Blocks []*models.Block `json:"blocks"`
}

// GetBlockchain return the block by height
func (service *Service) GetBlockchain(r *http.Request, params *BlockchainParams, result *BlockchainResult) error {
	return nil
}
