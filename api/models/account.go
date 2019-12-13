package models

import "github.com/QuoineFinancial/liquid-chain/abi"

// Account contains all info of a account on blockchain
type Account struct {
	Nonce        uint64      `json:"nonce"`
	ContractHash string      `json:"contractHash"`
	Contract     *abi.Contract `json:"contract"`
}
