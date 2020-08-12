package engine

import (
	"io/ioutil"
	"math"
	"time"

	"testing"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

func loadContract(abiPath, wasmPath string) *abi.Contract {
	header, err := abi.LoadHeaderFromFile(abiPath)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		panic(err)
	}

	return &abi.Contract{
		Header: header,
		Code:   data,
	}
}

func TestEngineIgnite(t *testing.T) {
	contractCreator, _ := crypto.AddressFromString("LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R")
	mathAddress, _ := crypto.AddressFromString("LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53")
	utilAddress, _ := crypto.AddressFromString("LCR57ROUHIQ2AV4D3E3D7ZBTR6YXMKZQWTI4KSHSWCUCRXBKNJKKBCNY")
	state := storage.NewStateStorage(db.NewMemoryDB())
	if err := state.LoadState(&crypto.BlockHeader{
		Height: 1,
		Time:   time.Unix(1578905663, 0),
	}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		callee        *abi.Contract
		calleeAddress crypto.Address
		caller        *abi.Contract
		callerAddress crypto.Address
		funcName      string
		args          []string
		want          uint64
		wantErr       bool
	}{
		{
			name:          "chained ignite",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "hypotenuse",
			args:          []string{"3", "4"},
			want:          math.Float64bits(5),
		},
		{
			name:          "chained ignite with array param",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "variance",
			args:          []string{"[1,2,3,4,5]"},
			want:          2,
		},
		{
			name:          "chained ignite with events",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "xor_checksum",
			args:          []string{"LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"},
			want:          149,
		},
		{
			name:          "ignite unknown imported function",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "average",
			args:          []string{"[1,2,3,4,5]"},
			wantErr:       true,
		},
		{
			name:          "chained ignite overflow",
			callee:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			calleeAddress: utilAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "mean",
			args:          []string{"[1,2,3,4,5]"},
			wantErr:       true,
		},
		{
			name:          "ignite block time",
			caller:        loadContract("testdata/blockinfo-abi.json", "testdata/blockinfo.wasm"),
			callerAddress: utilAddress,
			funcName:      "block_time",
			args:          []string{},
			want:          1578905663,
		},
		{
			name:          "ignite block height",
			caller:        loadContract("testdata/blockinfo-abi.json", "testdata/blockinfo.wasm"),
			callerAddress: utilAddress,
			funcName:      "block_height",
			args:          []string{},
			want:          1,
		},
		{
			name:          "chained ignite with invoke address validation",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "mod_invoke",
			args:          []string{"LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"},
			wantErr:       true,
		},
		{
			name:          "chained ignite with event address validation",
			callee:        loadContract("testdata/math-abi.json", "testdata/math.wasm"),
			calleeAddress: mathAddress,
			caller:        loadContract("testdata/util-abi.json", "testdata/util.wasm"),
			callerAddress: utilAddress,
			funcName:      "mod_emit",
			args:          []string{"LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"},
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contractBytes, _ := rlp.EncodeToBytes(&tt.caller)
			callerAccount, _ := state.CreateAccount(contractCreator, tt.callerAddress, contractBytes)
			execEngine := NewEngine(state, callerAccount, contractCreator, &gas.FreePolicy{}, 0)
			if tt.callee != nil {
				if tt.calleeAddress.String() != tt.callerAddress.String() {
					contractBytes, _ := rlp.EncodeToBytes(&tt.callee)
					state.CreateAccount(contractCreator, tt.calleeAddress, contractBytes)

				}
				// contract init
				initFunc := "init"
				function, err := tt.caller.Header.GetFunction(initFunc)
				if err != nil {
					panic(err)
				}
				args, err := abi.EncodeFromString(function.Parameters, []string{tt.calleeAddress.String()})
				if err != nil {
					panic(err)
				}
				_, err = execEngine.Ignite(initFunc, args)
				if err != nil {
					panic(err)
				}

			}

			// exec
			function, err := tt.caller.Header.GetFunction(tt.funcName)
			if err != nil {
				panic(err)
			}
			args, err := abi.EncodeFromString(function.Parameters, tt.args)
			if err != nil {
				panic(err)
			}
			got, err := execEngine.Ignite(tt.funcName, args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Engine.Ignite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Engine.Ignite() = %v, want %v", got, tt.want)
			}
		})
	}
}
