package models

import "time"

// KVPair is struct for key value pair
type KVPair struct {
	Key   []byte `json:"key,omitempty"`
	Value []byte `json:"value,omitempty"`
}

// Block contains block data
type Block struct {
	Hash      string    `json:"hash"`
	Height    int64     `json:"height"`
	Timestamp time.Time `json:"timestamp"`

	AppHash           string `json:"appHash"`
	ConsensusHash     string `json:"consensusHash"`
	PreviousBlockHash string `json:"previousBlockHash"`

	TxHashes []string       `json:"txHashes,omitempty"`
	Txs      []*Transaction `json:"txs,omitempty"`
}
