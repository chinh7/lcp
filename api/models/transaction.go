package models

// EventAttribute is attribute of event
type EventAttribute struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Event is emitting from contract execution
type Event struct {
	Name       string           `json:"name"`
	Contract   string           `json:"contract"`
	Attributes []EventAttribute `json:"attributes"`
}

// Transaction cointans all transactions info
type Transaction struct {
	Block *Block `json:"block,omitempty"`

	Hash  string `json:"hash"`
	Nonce uint64 `json:"nonce"`
	Code  uint32 `json:"code"`
	Data  string `json:"data"`
	Info  string `json:"info"`

	Contract string `json:"contract"`
	From     string `json:"from"`
	To       string `json:"to"`

	GasUsed  uint32 `json:"gasUsed"`
	GasLimit uint32 `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`

	Result uint64   `json:"result"`
	Events []*Event `json:"events"`
}
