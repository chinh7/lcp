package types

import (
	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"golang.org/x/crypto/blake2b"
)

// BlockHeader contains basic info and root hash of storage, transactions and reciepts
type BlockHeader struct {
	Time            uint
	Parent          Hash
	StateRoot       Hash
	ReceiptRoot     Hash
	TransactionRoot Hash
}

// Block is unit of Liquid chain
type Block struct {
	hash         string
	header       *BlockHeader
	transactions []*Transaction
}

// Hash returns blake2b hash of rlp encoding of block header
func (block *Block) Hash() Hash {
	encoded, _ := rlp.EncodeToBytes(block.header)
	hash := blake2b.Sum256(encoded)
	return hash
}
