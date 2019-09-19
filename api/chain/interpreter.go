package chain

import (
	"github.com/QuoineFinancial/vertex/api/models"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

func (service *Service) parseBlock(resultBlock *core_types.ResultBlock) *models.Block {
	block := &models.Block{
		Hash:      resultBlock.BlockMeta.BlockID.String(),
		Height:    resultBlock.BlockMeta.Header.Height,
		Timestamp: resultBlock.BlockMeta.Header.Time,

		AppHash:           resultBlock.BlockMeta.Header.AppHash.String(),
		ConsensusHash:     resultBlock.BlockMeta.Header.ConsensusHash.String(),
		PreviousBlockHash: resultBlock.BlockMeta.Header.LastBlockID.Hash.String(),
	}
	for _, tx := range resultBlock.Block.Data.Txs {
		block.TxHashes = append(block.TxHashes, tx.String())
	}
	return block
}

func (service *Service) parseTransaction(resultTx *core_types.ResultTx) *models.Transaction {
	return &models.Transaction{
		Hash:   resultTx.Hash.String(),
		Result: resultTx.TxResult.String(),
	}
}
