package chain

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/api/models"
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
	if tx, err := service.tAPI.Tx(hash, false); err != nil {
		return err
	} else if block, err := service.tAPI.Block(&tx.Height); err != nil {
		return err
	} else {
		result.Transaction = service.parseTransaction(tx)
		result.Transaction.Block = service.parseBlock(block)
	}
	return nil
}

// GetBlockTxs returns all transactions in given block
func (service *Service) GetBlockTxs(
	r *http.Request,
	params *GetBlockTxsParams,
	result *SearchTransactionResult,
) error {
	query := fmt.Sprintf("tx.height=%d", params.Height)
	return service.searchTransaction(query, params.Page, result)
}

// GetAccountTxs returns all transactions realted to given address
func (service *Service) GetAccountTxs(
	r *http.Request,
	params *GetAccountTxsParams,
	result *SearchTransactionResult,
) error {
	query := fmt.Sprintf("detail.from='%s'", params.Address)
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
	searchResult, err := service.tAPI.TxSearch(query, false, p, defaultTransactionPerPage)
	if err != nil {
		return err
	}
	result.Pagination = models.Pagination{
		CurrentPage: p,
		LastPage:    searchResult.TotalCount / defaultTransactionPerPage,
		Total:       searchResult.TotalCount,
	}
	for _, tx := range searchResult.Txs {
		transaction := service.parseTransaction(tx)
		result.Transactions = append(result.Transactions, transaction)

	}

	return nil
}
