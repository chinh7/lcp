package chain

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/QuoineFinancial/vertex/api/models"
)

// TransactionPerPage is number of transaction returns per page
const TransactionPerPage = 100

// GetTransactionByHashParams is params of GetTransactionByHash
type GetTransactionByHashParams struct {
	Hash string `json:"hash"`
}

// GetTransactionByHashResult is response of Service
type GetTransactionByHashResult struct {
	transaction *models.Transaction
}

// QueryTransactionsByBlockHeightParams is params of GetTxsByBlockHeight
type QueryTransactionsByBlockHeightParams struct {
	Height int
	Page   *int
}

// QueryTransactionsByAddressParams is params of GetTxsByBlockHeight
type QueryTransactionsByAddressParams struct {
	Address string
	Page    *int
}

// QueryTransactionsResult is response of query request
type QueryTransactionsResult struct {
	Transactions []*models.Transaction `json:"transactions"`
	Total        int                   `json:"total"`
}

// GetTransactionByHash is handler of Service
func (service *Service) GetTransactionByHash(
	r *http.Request,
	params *GetTransactionByHashParams,
	result *GetTransactionByHashResult,
) error {
	hash, _ := hex.DecodeString(params.Hash)
	if tx, err := service.tAPI.Tx(hash, false); err != nil {
		return err
	} else if block, err := service.tAPI.Block(&tx.Height); err != nil {
		return err
	} else {
		result.transaction = service.parseTransaction(tx)
		result.transaction.Block = service.parseBlock(block)
	}
	return nil
}

// QueryTransactionsByBlockHeight returns all transactions in given block
func (service *Service) QueryTransactionsByBlockHeight(
	r *http.Request,
	params *QueryTransactionsByBlockHeightParams,
	result *QueryTransactionsResult,
) error {
	query := fmt.Sprintf("tx.height='%d'", params.Height)
	return service.queryTransaction(query, params.Page, result)
}

// QueryTransactionsByAddress returns all transactions realted to given address
func (service *Service) QueryTransactionsByAddress(
	r *http.Request,
	params *QueryTransactionsByAddressParams,
	result *QueryTransactionsResult,
) error {
	query := fmt.Sprintf("account.address='%s'", params.Address)
	return service.queryTransaction(query, params.Page, result)
}

func (service *Service) queryTransaction(
	query string,
	page *int,
	result *QueryTransactionsResult,
) error {
	p := 0
	if page == nil {
		p = *page
	}
	searchResult, err := service.tAPI.TxSearch(query, false, p, TransactionPerPage)
	if err != nil {
		return err
	}
	for _, tx := range searchResult.Txs {
		transaction := service.parseTransaction(tx)
		result.Transactions = append(result.Transactions, transaction)
	}

	return nil
}
