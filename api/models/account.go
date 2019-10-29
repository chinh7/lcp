package models

// Account contains all info of a account on blockchain
type Account struct {
	Nonce        uint64 `json:"nonce"`
	ContractHash string `json:"contractHash"`
	Contract     string `json:"contract"`
}
