package abi

import (
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/go-cmp/cmp"
)

func TestDecodeContract(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/fixtures/header-event.json")
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
	if diff := cmp.Diff(*decodedContract, contract); diff != "" {
		t.Errorf("Decode contract %v is incorrect, expected: %v, got: %v, diff: %v", contract, contract, decodedContract, diff)
	}
}
