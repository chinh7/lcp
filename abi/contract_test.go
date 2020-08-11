package abi

import (
	"encoding/hex"
	"testing"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/nsf/jsondiff"
	"github.com/tendermint/tendermint/libs/os"
)

func TestDecodeContract(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/token-abi.json")
	contract := Contract{
		Header: h,
		Code:   []byte{1},
	}
	encodedContract, err := rlp.EncodeToBytes(&contract)
	if err != nil {
		t.Error(err)
	}

	decodedContract, err := DecodeContract(encodedContract)
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{})
	if diff := cmp.Diff(*decodedContract, contract, opts); diff != "" {
		t.Errorf("Decode contract %v is incorrect, expected: %v, got: %v, diff: %v", contract, contract, decodedContract, diff)
	}
}

func TestMarshalJSON(t *testing.T) {
	abiFile := "../test/testdata/token-abi.json"
	h, _ := LoadHeaderFromFile(abiFile)
	code, _ := hex.DecodeString("1")
	contract := Contract{
		Header: h,
		Code:   code,
	}
	jsonBytes, _ := contract.Header.MarshalJSON()
	expectedJSONBytes, _ := os.ReadFile(abiFile)
	if diff, result := jsondiff.Compare(jsonBytes, expectedJSONBytes, &jsondiff.Options{}); diff != jsondiff.FullMatch {
		t.Log(result)
		t.Error("JSON not matched")
	}
}
