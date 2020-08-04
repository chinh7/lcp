package chain

import (
	"net/http"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/consensus"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// BlockParams is params of ChainService
type BlockParams struct {
}

// BlockResult is response of GetBlock
type BlockResult struct {
	Block *crypto.Block `json:"block"`
}

// GetLatestBlock return the block by height
func (service *Service) GetLatestBlock(r *http.Request, params *BlockParams, result *BlockResult) error {
	lastBlockHash := service.app.InfoDB.Get([]byte(consensus.LastBlockHashKey))
	rawBlock := service.app.BlockDB.Get(lastBlockHash)
	var block crypto.Block
	rlp.DecodeBytes(rawBlock, &block)
	result.Block = &block
	return nil
}
