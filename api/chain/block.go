package chain

import (
	"net/http"

	"github.com/QuoineFinancial/vertex/api/models"
)

// BlockParams is params of ChainService
type BlockParams struct {
	Height int64 `json:"height"`
}

// BlockResult is response of GetBlock
type BlockResult struct {
	Block *models.Block `json:"block"`
}

// GetBlock return the block by height
func (service *Service) GetBlock(r *http.Request, params *BlockParams, result *BlockResult) error {
	tBlock, err := service.tAPI.Block(&params.Height)
	if err != nil {
		return err
	}
	result.Block = service.parseBlock(tBlock)
	return nil
}
