package models

// Account contains all info of a account on blockchain
type Account struct {
	Nonce    uint64 `json:"nonce"`
	CodeHash string `json:"codeHash"`
	Code     string `json:"code"`
}
