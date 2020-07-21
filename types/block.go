package types

import (
	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"golang.org/x/crypto/blake2b"
)

// BlockHeader contains basic info and root hash of storage, transactions and reciepts
type BlockHeader struct {
	Time            uint
	Parent          string
	StateHash       string
	ReceiptHash     string
	TransactionHash string
}

// Block is unit of Liquid chain
type Block struct {
	hash         string
	header       *BlockHeader
	transactions []*Transaction
}

// Hash returns blake2b hash of rlp encoding of block header
func (block *Block) Hash() common.Hash {
	encoded, _ := rlp.EncodeToBytes(block.header)
	hash := blake2b.Sum256(encoded)
	return hash
}
