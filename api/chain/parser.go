package chain

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/api/models"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/ethereum/go-ethereum/common"
	abciTypes "github.com/tendermint/tendermint/abci/types"
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

func (service *Service) parseTransaction(resultTx *core_types.ResultTx) (*models.Transaction, error) {
	transaction := &models.Transaction{
		Hash:     resultTx.Hash.String(),
		Info:     resultTx.TxResult.Info,
		GasUsed:  uint32(resultTx.TxResult.GetGasUsed()),
		GasLimit: uint32(resultTx.TxResult.GetGasWanted()),
		Code:     resultTx.TxResult.GetCode(),
		Data:     string(resultTx.TxResult.GetData()),
		Events:   []*models.Event{},
	}

	for _, e := range resultTx.TxResult.GetEvents() {
		parsedEvent, err := service.parseEvent(transaction, e)
		if err != nil {
			return nil, err
		}
		if parsedEvent != nil {
			transaction.Events = append(transaction.Events, parsedEvent)
		}
	}

	return transaction, nil
}

func parseEventName(name []byte) (*crypto.Address, uint32, error) {
	address := crypto.AddressFromBytes(name[0:crypto.AddressLength])
	index := binary.LittleEndian.Uint32(name[crypto.AddressLength:])
	return &address, index, nil
}

func (service *Service) parseEvent(tx *models.Transaction, tmEvent abciTypes.Event) (*models.Event, error) {
	name, err := hex.DecodeString(tmEvent.Type)
	var result models.Event
	if err != nil {
		return nil, err
	}

	if len(name) == 1 {
		eventCode := event.SystemEventCode(name[0])
		switch eventCode {
		case event.Detail:
			detailEvent := event.LoadDetailEvent(tmEvent)
			tx.From = detailEvent.From.String()
			tx.To = detailEvent.To.String()
			tx.Nonce = detailEvent.Nonce
			tx.Result = detailEvent.Result
			tx.GasPrice = uint32(detailEvent.GasPrice)
		case event.Deployment:
			tx.Contract = event.LoadDeploymentEvent(tmEvent).Address.String()
		}
		return nil, nil
	} else {
		contractAddress, index, err := event.ParseCustomEventName(name)
		if err != nil {
			return nil, err
		}

		status, _ := service.tAPI.Status()
		appHash := common.BytesToHash(status.SyncInfo.LatestAppHash)
		state, err := storage.New(appHash, service.database)
		if err != nil {
			return nil, err
		}
		account, err := state.GetAccount(*contractAddress)
		if err != nil {
			return nil, err
		}

		contract, err := account.GetContract()
		if err != nil {
			return nil, err
		}

		abiEvent, err := contract.Header.GetEventByIndex(index)
		if err != nil {
			return nil, err
		}

		result.Name = abiEvent.Name
		result.Contract = contractAddress.String()
		for index, param := range abiEvent.Parameters {
			valueByte, _ := hex.DecodeString(string(tmEvent.Attributes[index].Value))
			var value string
			if param.Type == abi.Address {
				address := crypto.AddressFromBytes(valueByte)
				value = address.String()
			} else {
				value = strconv.FormatUint(binary.LittleEndian.Uint64(valueByte), 10)
			}
			result.Attributes = append(result.Attributes, models.EventAttribute{
				Key:   param.Name,
				Type:  param.Type.String(),
				Value: value,
			})
		}
	}
	return &result, nil
}
