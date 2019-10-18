package models

import "time"

// Block contains block data
type Block struct {
	Hash   string    `json:"hash"`
	Time   time.Time `json:"time"`
	Height int64     `json:"height"`

	AppHash           string `json:"appHash"`
	ConsensusHash     string `json:"consensusHash"`
	PreviousBlockHash string `json:"previousBlockHash"`

	TxHashes []string       `json:"txHashes,omitempty"`
	Txs      []*Transaction `json:"txs,omitempty"`
}
