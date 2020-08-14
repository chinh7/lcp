package crypto

import (
	"bytes"
	"testing"
	"time"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/ed25519"
)

func TestBlock_Hash(t *testing.T) {
	type fields struct {
		header *BlockHeader
	}
	tests := []struct {
		name   string
		fields fields
		want   common.Hash
	}{{
		fields: fields{
			header: &BlockHeader{
				Time:            time.Unix(123, 0),
				Height:          1,
				Parent:          common.HexToHash("2f636344b757343e13e7910eed1b832d769e1d113027424580a2faca232ce015"),
				StateRoot:       common.HexToHash("572343bcdac17dbae1aba2d1ccde3488adb169b18da8a4ecdffe11c8f1cc1f1f"),
				TransactionRoot: common.HexToHash("3e2e21d19f5c3491ea8d5416b44256c401596b184638e63d8ac34f073a686544"),
			},
		},
		want: common.HexToHash("145eb771aa5c5f66132971083bcc5a2db139d83fa6b4ebf441f446c8b3ee0bef"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := &Block{
				Header: tt.fields.header,
			}
			if got := block.Header.Hash(); !cmp.Equal(got, tt.want) {
				t.Errorf("Block.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock(t *testing.T) {
	block := NewEmptyBlock(common.EmptyHash, 0, time.Unix(0, 0))
	block.Header.SetStateRoot(common.BytesToHash([]byte{1, 2, 3}))
	block.Header.SetTransactionRoot(common.BytesToHash([]byte{1, 2, 3}))
	block.Transactions = []*Transaction{{
		Version: 1,
		Sender: &TxSender{
			Nonce:     uint64(0),
			PublicKey: ed25519.NewKeyFromSeed(make([]byte, 32)).Public().(ed25519.PublicKey),
		},
		Receiver: Address{},
		Payload: &TxPayload{
			Contract: []byte{1, 2, 3},
			Method:   "Transfer",
			Params:   []byte{4, 5, 6},
		},
		GasPrice:  1,
		GasLimit:  2,
		Signature: []byte{7, 8, 9},
		Receipt: &TxReceipt{
			Result:  1,
			GasUsed: 2,
			Code:    ReceiptCodeOK,
			Events: []*TxEvent{{
				Contract: Address{},
				Data:     []byte{10, 11, 12},
			}},
		},
	}, {
		Version: 1,
		Sender: &TxSender{
			Nonce:     uint64(0),
			PublicKey: ed25519.NewKeyFromSeed(make([]byte, 32)).Public().(ed25519.PublicKey),
		},
		Receiver: Address{},
		Payload: &TxPayload{
			Contract: []byte{1, 1, 1},
			Method:   "Mint",
			Params:   []byte{4, 4, 4},
		},
		GasPrice:  1,
		GasLimit:  10,
		Signature: []byte{7, 8, 9},
		Receipt: &TxReceipt{
			Result:  1,
			GasUsed: 2,
			Code:    ReceiptCodeOK,
			Events: []*TxEvent{{
				Contract: Address{},
				Data:     []byte{10, 11, 12},
			}},
		},
	}}

	encoded, _ := block.Encode()
	decodedBlock := MustDecodeBlock(encoded)
	if decodedBlock.Header.Hash() != block.Header.Hash() {
		t.Errorf("Got block hash after decoded = %v, want %v", decodedBlock.Header.Hash(), block.Header.Hash())
	}

	if len(decodedBlock.Transactions) != len(block.Transactions) {
		t.Errorf("Encode transaction in block error")
	} else {
		for i := range decodedBlock.Transactions {
			if decodedBlock.Transactions[i].Hash() != block.Transactions[i].Hash() {
				t.Errorf("Encode transaction in block error")
			}
		}
	}

	encodedNew, _ := decodedBlock.Encode()
	if !bytes.Equal(encoded, encodedNew) {
		t.Errorf("Encode not equal, got = %v, want %v", encodedNew, encoded)
	}
}

func TestMustDecodeBlock(t *testing.T) {
	// This decoding should panic
	defer func() { recover() }()
	MustDecodeBlock([]byte{1, 2, 3})
	t.Errorf("did not panic")
}
