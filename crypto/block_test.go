package crypto

import (
	"reflect"
	"testing"
	"time"

	"github.com/QuoineFinancial/liquid-chain/common"
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
			if got := block.Header.Hash(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Block.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
