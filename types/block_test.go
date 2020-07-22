package types

import (
	"reflect"
	"testing"
)

func TestBlock_Hash(t *testing.T) {
	type fields struct {
		header *BlockHeader
	}
	tests := []struct {
		name   string
		fields fields
		want   Hash
	}{{
		fields: fields{
			header: &BlockHeader{
				Time:            0,
				Parent:          "2f636344b757343e13e7910eed1b832d769e1d113027424580a2faca232ce015",
				StateRoot:       "572343bcdac17dbae1aba2d1ccde3488adb169b18da8a4ecdffe11c8f1cc1f1f",
				ReceiptRoot:     "497addcfff879adf6ca5c24fcb4d955ad0082eef374fbc7af7c55844594e09b0",
				TransactionRoot: "3e2e21d19f5c3491ea8d5416b44256c401596b184638e63d8ac34f073a686544",
			},
		},
		want: HexToHash("1518212b8c19c5b319d0098330c09ffb6e7b32e6053ebb9ab02c0ea7e370030f"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := &Block{
				header: tt.fields.header,
			}
			if got := block.Hash(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Block.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
