package abi

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	h, _ := LoadHeaderFromFile("../test/testdata/token-abi.json")
	code, _ := hex.DecodeString("1")
	contract := Contract{
		Header: h,
		Code:   code,
	}
	jsonBytes, _ := contract.MarshalJSON()

	var decodedContract struct {
		Header *Header `json:"header"`
		Code   string  `json:"code"`
	}

	json.Unmarshal(jsonBytes, &decodedContract)

	opts := cmpopts.IgnoreUnexported(Event{})
	if diff := cmp.Diff(decodedContract.Header, contract.Header, opts); diff != "" {
		t.Errorf("Decode contract %v is incorrect, expected: %v, got: %v, diff: %v", contract, contract.Header, decodedContract.Header, diff)
	}
	if diff := cmp.Diff(decodedContract.Code, string(contract.Code), opts); diff != "" {
		t.Errorf("Decode contract %v is incorrect, expected: %v, got: %v, diff: %v", contract, contract.Code, decodedContract.Code, diff)
	}
}
