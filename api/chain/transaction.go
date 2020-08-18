package chain

import (
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/crypto"
)

const defaultTransactionPerPage = int(50)
const startPage = int(0)

// GetTxParams is params of GetTx
type GetTxParams struct {
	Hash string `json:"hash"`
}

// GetTxResult is response of Service
type GetTxResult struct {
	Transaction *crypto.Transaction `json:"tx"`
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
	Transactions []*crypto.Transaction `json:"txs"`
}

// GetTx is handler of Service
func (service *Service) GetTx(r *http.Request, params *GetTxParams, result *GetTxResult) error {
	return nil
}

// GetBlockTxs returns all transactions in given block
func (service *Service) GetBlockTxs(r *http.Request, params *GetBlockTxsParams, result *SearchTransactionResult) error {
	return nil
}

// GetAccountTxs returns all transactions realted to given address
func (service *Service) GetAccountTxs(r *http.Request, params *GetAccountTxsParams, result *SearchTransactionResult) error {
	return nil
}

func (service *Service) GetEventTxs(r *http.Request, params *GetEventTxsParams, result *SearchTransactionResult) error {
	return nil
}

/* TODO: Technical reviews
- Is customizable perPage necessary or not?
- Which value of perPage is suitable? (50, 100, or block capacity?)
*/
func (service *Service) searchTransaction(query string, page *int, result *SearchTransactionResult) error {
	return nil
}
