package api

import (
	"encoding/hex"
	"fmt"
	"net/http"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// TxByHashArgs is params of GetTxByHash
type TxByHashArgs struct {
	Hash string
}

// TxsByBlockHeightArgs is params of GetTxsByBlockHeight
type TxsByBlockHeightArgs struct {
	Height  int
	Page    int
	PerPage int
}

// TxsByAddressArgs is params of GetTxsByBlockHeight
type TxsByAddressArgs struct {
	Address string
	Page    int
	PerPage int
}

// TxEvent is struct for response of Events of Tx
type TxEvent struct {
	Type       string   `json:"type,omitempty"`
	Attributes []KVPair `json:"attributes,omitempty"`
}

// TxResult is struct for response of Tx
type TxResult struct {
	Data   []byte    `json:"data,omitempty"`
	Events []TxEvent `json:"events,omitempty"`
}

// TxReply is response of TransactionService
type TxReply struct {
	Hash     string   `json:"hash"`
	Height   int64    `json:"height"`
	Index    uint32   `json:"index"`
	Tx       []byte   `json:"tx"`
	TxResult TxResult `json:"tx_result"`
}

// TxsByBlockHeightReply is response for Txs
type TxsByBlockHeightReply struct {
	Transactions []TxReply `json:"transactions"`
	TotalCount   int       `json:"total_count"`
}

// TransactionService is first service
type TransactionService struct {
	client *rpcclient.Client
}

// NewTransactionService returns new instance of TransactionService
func (api *API) NewTransactionService() *TransactionService {
	if api.Client == nil {
		panic("api.NewTransactionService call without api.Client")
	}
	return &TransactionService{api.Client}
}

func mappingTxTypes(tx ctypes.ResultTx) TxReply {
	var txReply TxReply
	txReply.Hash = tx.Hash.String()
	txReply.Height = tx.Height
	txReply.Index = tx.Index
	txReply.Tx = tx.Tx
	txReply.TxResult.Data = tx.TxResult.Data
	for _, event := range tx.TxResult.Events {
		var attributes []KVPair
		for _, cKVPair := range event.Attributes {
			attributes = append(attributes, KVPair{Key: cKVPair.Key, Value: cKVPair.Value})
		}
		e := TxEvent{Type: event.Type, Attributes: attributes}
		txReply.TxResult.Events = append(txReply.TxResult.Events, e)
	}
	return txReply
}

// GetTxByHash is handler of TransactionService
func (service *TransactionService) GetTxByHash(r *http.Request, args *TxByHashArgs, reply *TxReply) error {
	client := *service.client
	hashBytes, err := hex.DecodeString(args.Hash)
	if err != nil {
		return err
	}

	tx, err := client.Tx(hashBytes, false)
	if err != nil {
		return err
	}
	result := mappingTxTypes(*tx)
	reply = &result
	return nil
}

// GetTxsByBlockHeight is handler of TransactionService
func (service *TransactionService) GetTxsByBlockHeight(r *http.Request, args *TxsByBlockHeightArgs, reply *TxsByBlockHeightReply) error {
	client := *service.client
	query := fmt.Sprintf("tx.height='%d'", args.Height)
	txs, err := client.TxSearch(query, false, args.Page, args.PerPage)
	if err != nil {
		return err
	}
	reply.TotalCount = txs.TotalCount
	for _, tx := range txs.Txs {
		transaction := mappingTxTypes(*tx)
		reply.Transactions = append(reply.Transactions, transaction)
	}
	return nil
}

// GetTxsByAccount is handler of TransactionService
func (service *TransactionService) GetTxsByAccount(r *http.Request, args *TxsByAddressArgs, reply *TxsByBlockHeightReply) error {
	client := *service.client
	query := fmt.Sprintf("account.address='%s'", args.Address)

	txs, err := client.TxSearch(query, false, args.Page, args.PerPage)
	if err != nil {
		return err
	}
	reply.TotalCount = txs.TotalCount
	for _, tx := range txs.Txs {
		transaction := mappingTxTypes(*tx)
		reply.Transactions = append(reply.Transactions, transaction)
	}
	return nil
}
