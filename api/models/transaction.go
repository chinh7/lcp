package models

// Transaction cointans all transactions info
type Transaction struct {
	Block *Block `json:"block,omitempty"`

	Hash  string `json:"hash"`
	Nonce string `json:"nonce"`

	From string `json:"from"`
	To   string `json:"to"`

	GasUsed  string `json:"gasUsed"`
	GasPrice string `json:"gasPrice"`
	GasLimit string `json:"gasLimit"`

	Result string `json:"result"`
}
