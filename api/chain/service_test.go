package chain

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/stretchr/testify/assert"
)

var testResourceInstance *testResource

func TestMain(m *testing.M) {
	testResourceInstance = newTestResource()
	testResourceInstance.seed()
	code := m.Run()
	testResourceInstance.tearDown()
	os.Exit(code)
}

func TestGetLatestBlock(t *testing.T) {
	var result BlockResult
	testResourceInstance.service.GetLatestBlock(nil, &LatestBlockParams{}, &result)
	assert.Equal(t, block{
		Hash:            common.HexToHash("7ad9917bf11abdffc2be47e966d1cfabe4d573e2cdd5e83db1baa5b94ada21e1"),
		Height:          4,
		Time:            4,
		Parent:          common.HexToHash("49ec9a40849711c6207c458024c150b1dd306b4ae351902b1a049f5e90fe1ff7"),
		StateRoot:       common.HexToHash("4932482ebd2cc8031f0a44de9158ec1fd7c9d388a4ad2de02e914c1babb823f7"),
		TransactionRoot: common.HexToHash("45b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c0"),
		ReceiptRoot:     common.HexToHash("45b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c0"),
		Transactions:    []transaction{},
		Receipts:        []receipt{},
	}, *result.Block)
}

func TestGetBlockByHeight(t *testing.T) {
	var result BlockResult
	testResourceInstance.service.GetBlockByHeight(nil, &BlockByHeightParams{
		Height: 2,
	}, &result)

	sender, _ := crypto.AddressFromString("LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT")
	receiver, _ := crypto.AddressFromString("LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH")
	signature, _ := base64.StdEncoding.DecodeString("OttqA4/C5Bk/04EMSXsBZ8U8bNWb4ErsBwStsdo4gDuV9kKEdb2Z/TEr9WQb100e7gj3g1meyKVinI2ZbjGcBg==")

	assert.Equal(t, block{
		Time:            2,
		Height:          2,
		Hash:            common.HexToHash("c150cbb67250266d77d573c8603ccd88f13cecb527e9e4fde208acbe4d078601"),
		Parent:          common.HexToHash("838330cdd2952a26233f30dd805e449d35bd14eec1cc4f53b0af7229dfbc2c51"),
		StateRoot:       common.HexToHash("50f3d1d9c48e5d967600fcc19ba7c6fd124bcef6d0c3701134122af51cd2e6e2"),
		TransactionRoot: common.HexToHash("eb9f7258591dca0b1d890ecb5e506d67c9e3d16f1fc6e3be40eeb084d78df501"),
		ReceiptRoot:     common.HexToHash("215b491384a59b0aef1a4764f8869a4da55686ff457191b50fb72519583fd588"),

		Transactions: []transaction{{
			Hash:        common.HexToHash("b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc"),
			Type:        "invoke",
			BlockHeight: 2,
			Version:     1,
			Sender:      sender,
			Nonce:       1,
			Receiver:    receiver,
			GasPrice:    1,
			GasLimit:    0,
			Signature:   signature,
			Payload: call{
				Name: "mint",
				Args: []argument{{
					Type:  "uint64",
					Name:  "amount",
					Value: "1000",
				}},
			},
		}},

		Receipts: []receipt{{
			Transaction: common.HexToHash("b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc"),
			Result:      "0",
			GasUsed:     0,
			Code:        0,
			Events: []call{{
				Contract: receiver.String(),
				Name:     "Mint",
				Args: []argument{{
					Type:  "address",
					Name:  "to",
					Value: "LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT",
				}, {
					Type:  "uint64",
					Name:  "amount",
					Value: "1000",
				}},
			}},
			PostState: common.HexToHash("50f3d1d9c48e5d967600fcc19ba7c6fd124bcef6d0c3701134122af51cd2e6e2"),
		}},
	}, *result.Block)
}

func TestGetTransaction(t *testing.T) {
	var result GetTransactionResult
	testResourceInstance.service.GetTransaction(nil, &GetTransactionParams{
		Hash: "b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc",
	}, &result)

	sender, _ := crypto.AddressFromString("LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT")
	receiver, _ := crypto.AddressFromString("LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH")
	signature, _ := base64.StdEncoding.DecodeString("OttqA4/C5Bk/04EMSXsBZ8U8bNWb4ErsBwStsdo4gDuV9kKEdb2Z/TEr9WQb100e7gj3g1meyKVinI2ZbjGcBg==")

	assert.Equal(t, transaction{
		Hash:        common.HexToHash("b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc"),
		Type:        "invoke",
		BlockHeight: 2,
		Version:     1,
		Sender:      sender,
		Nonce:       1,
		Receiver:    receiver,
		GasPrice:    1,
		GasLimit:    0,
		Signature:   signature,
		Payload: call{
			Name: "mint",
			Args: []argument{{
				Type:  "uint64",
				Name:  "amount",
				Value: "1000",
			}},
		},
	}, *result.Transaction)

	assert.Equal(t, receipt{
		Transaction: common.HexToHash("b3fef26e5cb52f0681a06bb9c9ff78acb53ef79daa0c46fecbf1e212e9a67ddc"),
		Result:      "0",
		GasUsed:     0,
		Code:        0,
		Events: []call{{
			Name:     "Mint",
			Contract: receiver.String(),
			Args: []argument{{
				Type:  "address",
				Name:  "to",
				Value: "LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT",
			}, {
				Type:  "uint64",
				Name:  "amount",
				Value: "1000",
			}},
		}},
		PostState: common.HexToHash("50f3d1d9c48e5d967600fcc19ba7c6fd124bcef6d0c3701134122af51cd2e6e2"),
	}, *result.Receipt)
}

func newUint64(value uint64) *uint64 {
	return &value
}

func TestCall(t *testing.T) {
	tests := []struct {
		name    string
		params  CallParams
		result  CallResult
		wantErr bool
	}{{
		name: "valid",
		params: CallParams{
			Height:  nil,
			Address: "LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
			Method:  "get_balance",
			Args:    []string{"LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT"},
		},
		result: CallResult{
			Result: "1000",
			Code:   crypto.ReceiptCodeOK,
			Events: []*call{},
		},
		wantErr: false,
	}, {
		name: "invalid address",
		params: CallParams{
			Height:  newUint64(1),
			Address: "invalid_address",
		},
		wantErr: true,
	}, {
		name: "call nil contract",
		params: CallParams{
			Height:  newUint64(1),
			Address: "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53",
		},
		wantErr: true,
	}, {
		name: "call not a contract",
		params: CallParams{
			Height:  newUint64(1),
			Address: "LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT",
		},
		wantErr: true,
	}, {
		name: "invalid function",
		params: CallParams{
			Height:  newUint64(1),
			Address: "LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
			Method:  "invalid_function",
		},
		wantErr: true,
	}, {
		name: "invalid params",
		params: CallParams{
			Height:  newUint64(1),
			Address: "LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
			Method:  "get_balance",
			Args:    []string{},
		},
		wantErr: true,
	}, {
		name: "ignite with events",
		params: CallParams{
			Address: "LA3K6XGDQXAZN6J22J5VCEFIU25PE4BEZRZE5K76WDGUIRV3HLKJALPV",
			Method:  "say",
			Args:    []string{"1"},
		},
		result: CallResult{
			Result: "1",
			Code:   0,
			Events: []*call{{
				Contract: "LA3K6XGDQXAZN6J22J5VCEFIU25PE4BEZRZE5K76WDGUIRV3HLKJALPV",
				Name:     "Say",
				Args: []argument{{
					Type:  "lparray",
					Name:  "message",
					Value: "Q2hlY2tpbmc=",
				}},
			}},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result CallResult
			if err := testResourceInstance.service.Call(nil, &tt.params, &result); (err != nil) != tt.wantErr {
				t.Errorf("Service.Call() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestGetAccount(t *testing.T) {
	sender, _ := crypto.AddressFromString("LA5WUJ54Z23KILLCUOUNAKTPBVZWKMQVO4O6EQ5GHLAERIMLLHNCTXXT")

	tests := []struct {
		name    string
		params  GetAccountParams
		result  GetAccountResult
		wantErr bool
	}{{
		name: "valid",
		params: GetAccountParams{
			Address: "LBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
		},
		result: GetAccountResult{
			Account: &storage.Account{
				Nonce:        0,
				Creator:      sender,
				StorageHash:  common.Hash{0x29, 0xc3, 0x5c, 0xda, 0xdc, 0x63, 0x49, 0xf, 0xb9, 0x2d, 0xdf, 0x18, 0x80, 0xc0, 0xb2, 0x98, 0x29, 0xb2, 0xab, 0x82, 0x1d, 0xf9, 0x18, 0x58, 0x2f, 0xef, 0x98, 0x9, 0x5, 0xf1, 0x88, 0x5c},
				ContractHash: common.Hash{0xd8, 0x9a, 0xb7, 0x4c, 0xc7, 0xf9, 0x5c, 0x3, 0xd5, 0x7d, 0xc6, 0x76, 0xee, 0xeb, 0x9d, 0xfc, 0x78, 0x15, 0xde, 0xe8, 0xc0, 0x5d, 0x7b, 0x2a, 0xe2, 0x8b, 0x7, 0xee, 0x5f, 0x6a, 0xa1, 0x4},
			},
		},
		wantErr: false,
	}, {
		name: "invalid address",
		params: GetAccountParams{
			Address: "MBAPQ4LVHFYZQXRSS3CCN6VUZ2EEC6IN5S2RGQLHS3RNNOIBNP4B6XNH",
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result GetAccountResult
			err := testResourceInstance.service.GetAccount(nil, &tt.params, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				assert.Equal(t, tt.result.Account.Nonce, result.Account.Nonce, "Nonce")
				assert.Equal(t, tt.result.Account.Creator, result.Account.Creator, "Creator")
				assert.Equal(t, tt.result.Account.StorageHash, result.Account.StorageHash, "StorageHash")
				assert.Equal(t, tt.result.Account.ContractHash, result.Account.ContractHash, "ContractHash")
			}
		})
	}
}
