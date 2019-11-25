package models

// Event is emitting from contract execution
type Event struct {
	Name       string `json:"name"`
	Attributes []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"attributes"`
}

// Transaction cointans all transactions info
type Transaction struct {
	Block *Block `json:"block,omitempty"`

	Hash  string `json:"hash"`
	Nonce int64  `json:"nonce"`
	Code  uint32 `json:"code"`
	Data  string `json:"data"`
	Info  string `json:"info"`

	From string `json:"from"`
	To   string `json:"to"`

	GasUsed  int64  `json:"gasUsed"`
	GasLimit int64  `json:"gasLimit"`
	GasPrice string `json:"-"`

	Result map[string]string `json:"result"`
	Events []*Event          `json:"events"`
}
