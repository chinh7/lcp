package chain

import (
	"encoding/hex"
	"strconv"

	"github.com/QuoineFinancial/vertex/api/models"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

func (service *Service) parseBlock(resultBlock *core_types.ResultBlock) *models.Block {
	block := &models.Block{
		Hash:      resultBlock.BlockMeta.BlockID.Hash.String(),
		Height:    resultBlock.BlockMeta.Header.Height,
		Timestamp: resultBlock.BlockMeta.Header.Time,

		AppHash:           resultBlock.BlockMeta.Header.AppHash.String(),
		ConsensusHash:     resultBlock.BlockMeta.Header.ConsensusHash.String(),
		PreviousBlockHash: resultBlock.BlockMeta.Header.LastBlockID.Hash.String(),
	}
	for _, tx := range resultBlock.Block.Data.Txs {
		block.TxHashes = append(block.TxHashes, hex.EncodeToString(tx.Hash()))
	}
	return block
}

func (service *Service) parseTransaction(resultTx *core_types.ResultTx) *models.Transaction {
	transaction := &models.Transaction{
		Hash:     resultTx.Hash.String(),
		Info:     resultTx.TxResult.Info,
		GasUsed:  resultTx.TxResult.GetGasUsed(),
		GasLimit: resultTx.TxResult.GetGasWanted(),
		Code:     resultTx.TxResult.GetCode(),
		Data:     string(resultTx.TxResult.GetData()),
		Result:   make(map[string]string),
	}

	for _, event := range resultTx.TxResult.GetEvents() {
		switch event.Type {
		case "result":
			for _, attribute := range event.GetAttributes() {
				transaction.Result[string(attribute.GetKey())] = string(attribute.GetValue())
			}
		case "detail":
			for _, attribute := range event.GetAttributes() {
				switch string(attribute.GetKey()) {
				case "to":
					transaction.To = string(attribute.GetValue())
				case "from":
					transaction.From = string(attribute.GetValue())
				case "nonce":
					nonce, _ := strconv.Atoi(string(attribute.GetValue()))
					transaction.Nonce = int64(nonce)
				}
			}
		}
	}

	return transaction
}
