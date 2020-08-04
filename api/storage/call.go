package storage

import (
	"net/http"

	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// CallParams is params to execute Call
type CallParams struct {
	Height  int64    `json:"height"`
	Address string   `json:"address"`
	Method  string   `json:"method"`
	Args    []string `json:"args"`
}

// CallResult is result of Call
type CallResult struct {
	Events []*crypto.TxEvent `json:"events"`
	Return interface{}       `json:"value"`
}

// Call to execute function without tx creation in blockchain
func (service *Service) Call(r *http.Request, params *CallParams, result *CallResult) error {
	// TODO: Add Call API
	return nil
}
