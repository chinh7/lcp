package chain

import (
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// BlockByHeightParams contains query height
type BlockByHeightParams struct {
	Height uint64 `json:"height"`
}

// BlockResult is response of GetBlock
type BlockResult struct {
	Block *crypto.Block `json:"block"`
}

// GetLatestBlock return the block by height
func (service *Service) GetLatestBlock(r *http.Request, _ interface{}, result *BlockResult) error {
	return nil
}

// GetBlockByHeight return block by its height
func (service *Service) GetBlockByHeight(r *http.Request, params *BlockByHeightParams, result *BlockResult) error {
	return nil
}
