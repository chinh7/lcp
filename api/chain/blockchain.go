package chain

import (
	"net/http"

	"github.com/QuoineFinancial/vertex/api/models"
)

const blocksPerPage = int(20)

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
	tBlockchain, err := service.tAPI.BlockchainInfo(params.Min, params.Max)
	if err != nil {
		return err
	}
	for _, blockMeta := range tBlockchain.BlockMetas {
		result.Blocks = append(result.Blocks, service.parseBlockMeta(blockMeta))
	}
	return nil
}
