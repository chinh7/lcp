package crypto

import (
	"log"
	"time"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"golang.org/x/crypto/blake2b"
)

// GenesisBlock is the first block of liquid chain
var GenesisBlock = Block{
	Header: &BlockHeader{
		Height:          0,
		Time:            time.Unix(0, 0),
		Parent:          common.EmptyHash,
		StateRoot:       common.EmptyHash,
		TransactionRoot: common.EmptyHash,
	},
	Transactions: nil,
}

// BlockHeader contains basic info and root hash of storage, transactions and receipts
type BlockHeader struct {
	Height          uint64
	Time            time.Time
	Parent          common.Hash
	StateRoot       common.Hash
	TransactionRoot common.Hash
}

// Block is unit of Liquid chain
type Block struct {
	Header       *BlockHeader
	Transactions []*Transaction
}

// SetStateRoot set StateRoot of header
func (blockHeader *BlockHeader) SetStateRoot(hash common.Hash) {
	blockHeader.StateRoot = hash
}

// SetTransactionRoot set TransactionRoot of header
func (blockHeader *BlockHeader) SetTransactionRoot(hash common.Hash) {
	blockHeader.TransactionRoot = hash
}

// Hash returns blake2b hash of rlp encoding of block header
func (blockHeader *BlockHeader) Hash() common.Hash {
	encoded, _ := rlp.EncodeToBytes(blockHeader)
	return blake2b.Sum256(encoded)
}

// NewEmptyBlock create empty block
func NewEmptyBlock(parent common.Hash, height uint64, blockTime time.Time) *Block {
	return &Block{
		Header: &BlockHeader{
			Parent: parent,
			Height: height,
			Time:   blockTime,
		},
	}
}

// Encode returns bytes array of block
func (block *Block) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(block)
}

// DecodeBlock returns block from encoded byte array
func DecodeBlock(rawBlock []byte) (*Block, error) {
	var block Block
	if err := rlp.DecodeBytes(rawBlock, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// MustDecodeBlock acts like DecodeBlock, but panic in case of error
func MustDecodeBlock(rawBlock []byte) *Block {
	block, err := DecodeBlock(rawBlock)
	if err != nil {
		log.Fatal(err)
	}
	return block
}
