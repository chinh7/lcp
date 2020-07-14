package chain

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/api/models"
	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

const defaultTransactionPerPage = int(50)
const startPage = int(0)

// GetTxParams is params of GetTx
type GetTxParams struct {
	Hash string `json:"hash"`
}

// GetTxResult is response of Service
type GetTxResult struct {
	Transaction *models.Transaction `json:"tx"`
}

// GetBlockTxsParams is params of GetTxsByBlockHeight
type GetBlockTxsParams struct {
	Height int  `json:"height"`
	Page   *int `json:"page"`
}

// GetEventTxsParams is params of GetEventTxs
type GetEventTxsParams struct {
	Contract string `json:"contract"`
	Event    Event  `json:"event"`
	Page     *int   `json:"page"`
}

// Event is including in GetEventTxsParams
type Event struct {
	Name  string     `json:"name"`
	Param EventParam `json:"param"`
}

// EventParam is param of Event
type EventParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetAccountTxsParams is params of GetTxsByBlockHeight
type GetAccountTxsParams struct {
	Address string `json:"address"`
	Page    *int   `json:"page"`
}

// SearchTransactionResult is response of query request
type SearchTransactionResult struct {
	Transactions []*models.Transaction `json:"txs"`
	Pagination   models.Pagination     `json:"pagination"`
}

// GetTx is handler of Service
func (service *Service) GetTx(
	r *http.Request,
	params *GetTxParams,
	result *GetTxResult,
) error {
	hash, _ := hex.DecodeString(params.Hash)
	tx, err := service.tAPI.Tx(hash, false)
	if err != nil {
		return err
	}
	parsedTransaction, err := service.parseTransaction(tx)
	if err != nil {
		return err
	}
	result.Transaction = parsedTransaction
	return nil
}

// GetBlockTxs returns all transactions in given block
func (service *Service) GetBlockTxs(
	r *http.Request,
	params *GetBlockTxsParams,
	result *SearchTransactionResult,
) error {
	heightBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(heightBytes, uint64(params.Height))
	heightHex := hex.EncodeToString(heightBytes)
	detailString := event.Detail.String()
	var heightParam string
	for index, param := range event.Detail.GetEvent().Parameters {
		if param.Name == "height" {
			heightParam = hex.EncodeToString([]byte{byte(index)})
		}
	}
	query := fmt.Sprintf("%s.%s='%s'", detailString, heightParam, heightHex)
	return service.searchTransaction(query, params.Page, result)
}

// GetAccountTxs returns all transactions realted to given address
func (service *Service) GetAccountTxs(
	r *http.Request,
	params *GetAccountTxsParams,
	result *SearchTransactionResult,
) error {
	detailString := event.Detail.String()
	var fromParam string
	for index, param := range event.Detail.GetEvent().Parameters {
		if param.Name == "from" {
			fromParam = hex.EncodeToString([]byte{byte(index)})
		}
	}
	address, _ := crypto.AddressFromString(params.Address)
	query := fmt.Sprintf("%s.%s='%s'", detailString, fromParam, hex.EncodeToString(address[:]))
	return service.searchTransaction(query, params.Page, result)
}

// GetEventTxs search for transactions based on an event
func (service *Service) GetEventTxs(
	r *http.Request,
	params *GetEventTxsParams,
	result *SearchTransactionResult,
) error {
	contractAddress, err := crypto.AddressFromString(params.Contract)
	if err != nil {
		return err
	}
	status, err := service.tAPI.Status()
	if err != nil {
		return err
	}
	appHash := common.BytesToHash(status.SyncInfo.LatestAppHash)
	state, err := storage.New(appHash, service.database)
	if err != nil {
		return err
	}
	account, err := state.GetAccount(contractAddress)
	if err != nil {
		return err
	}
	contract, err := account.GetContract()
	if err != nil {
		return err
	}
	event, err := contract.Header.GetEvent(params.Event.Name)
	if err != nil {
		return err
	}
	eventName := hex.EncodeToString(append(contractAddress[:], event.GetIndexByte()...))
	var eventParamKey, eventParamValue string
	for index, param := range event.Parameters {
		if param.Name == params.Event.Param.Name {
			eventParamKey = hex.EncodeToString([]byte{byte(index)})
			encode, err := abi.EncodeFromString([]*abi.Parameter{param}, []string{params.Event.Param.Value})
			if err != nil {
				return err
			}
			eventParamValue = hex.EncodeToString(encode[2:])
		}
	}

	query := fmt.Sprintf("%s.%s='%s'", eventName, eventParamKey, eventParamValue)
	return service.searchTransaction(query, params.Page, result)
}

/* TODO: Technical reviews
- Is customizable perPage necessary or not?
- Which value of perPage is suitable? (50, 100, or block capacity?)
*/
func (service *Service) searchTransaction(query string, page *int, result *SearchTransactionResult) error {
	p := startPage
	if page != nil {
		p = *page
	}
	searchResult, err := service.tAPI.TxSearch(query, false, p, defaultTransactionPerPage, "desc")
	if err != nil {
		return err
	}
	result.Pagination = models.Pagination{
		CurrentPage: p,
		LastPage:    searchResult.TotalCount / defaultTransactionPerPage,
		Total:       searchResult.TotalCount,
	}
	for _, tx := range searchResult.Txs {
		parsedTransaction, err := service.parseTransaction(tx)
		if err != nil {
			return err
		}
		result.Transactions = append(result.Transactions, parsedTransaction)
	}

	return nil
}
