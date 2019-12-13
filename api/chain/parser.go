package chain

import (
	"encoding/hex"
	"strconv"

	"github.com/QuoineFinancial/liquid-chain/api/models"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

func (service *Service) parseBlockMeta(resultBlockMeta *types.BlockMeta) *models.Block {
	return &models.Block{
		Hash:              resultBlockMeta.BlockID.Hash.String(),
		Time:              resultBlockMeta.Header.Time,
		Height:            resultBlockMeta.Header.Height,
		AppHash:           resultBlockMeta.Header.AppHash.String(),
		ConsensusHash:     resultBlockMeta.Header.ConsensusHash.String(),
		PreviousBlockHash: resultBlockMeta.Header.LastBlockID.Hash.String(),
	}
}

func (service *Service) parseBlock(resultBlock *core_types.ResultBlock) *models.Block {
	block := service.parseBlockMeta(resultBlock.BlockMeta)
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
		Events:   []*models.Event{},
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
		default:
			vEvent := models.Event{Name: event.GetType()}
			for _, attribute := range event.GetAttributes() {
				vEvent.Attributes = append(vEvent.Attributes, struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				}{
					Key:   string(attribute.Key),
					Value: hex.EncodeToString(attribute.Value),
				})
			}
			transaction.Events = append(transaction.Events, &vEvent)

		}
	}

	return transaction
}
