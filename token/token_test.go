package token

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/QuoineFinancial/liquid-chain/trie"

	"github.com/ethereum/go-ethereum/rlp"
)

const contractAddress = "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"
const ownerAddress = "LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"
const otherAddress = "LCR57ROUHIQ2AV4D3E3D7ZBTR6YXMKZQWTI4KSHSWCUCRXBKNJKKBCNY"
const nonExistentAddress = "LANXBHFABEPW5NDSIZUEIENR2LNQHYJ6464NYFVPLE6XKHTMCEZDCLM5"
const contractBalance = uint64(4319)
const ownerBalance = uint64(1000000000 - 10000 - 4319)
const otherBalance = uint64(10000)

func setup() *Token {
	db := db.NewMemoryDB()
	state, err := storage.New(trie.Hash{}, db)
	if err != nil {
		panic(err)
	}

	header, err := abi.LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile("../test/testdata/liquid-token.wasm")
	if err != nil {
		panic(err)
	}
	contract := &abi.Contract{
		Header: header,
		Code:   data,
	}
	contractBytes, err := rlp.EncodeToBytes(&contract)
	if err != nil {
		panic(err)
	}
	_, err = state.CreateAccount(crypto.AddressFromString(ownerAddress), crypto.AddressFromString(contractAddress), contractBytes)
	if err != nil {
		panic(err)
	}
	contractAccount, err := state.GetAccount(crypto.AddressFromString(contractAddress))
	if err != nil {
		panic(err)
	}
	token := NewToken(state, contractAccount)
	_, _, err = token.invokeContract(crypto.AddressFromString(ownerAddress), "mint", []string{strconv.FormatUint(1000000000, 10)})
	if err != nil {
		panic(err)
	}
	_, err = token.Transfer(crypto.AddressFromString(ownerAddress), crypto.AddressFromString(otherAddress), 10000)
	if err != nil {
		panic(err)
	}
	_, err = token.Transfer(crypto.AddressFromString(ownerAddress), crypto.AddressFromString(contractAddress), 4319)
	if err != nil {
		panic(err)
	}
	return token
}

func TestGetBalance(t *testing.T) {
	token := setup()
	ret, err := token.GetBalance(crypto.AddressFromString(contractAddress))
	if err != nil {
		panic(err)
	}
	if ret != contractBalance {
		t.Errorf("Expect contract balance to be %v, got %v", contractBalance, ret)
	}
	ret, err = token.GetBalance(crypto.AddressFromString(ownerAddress))
	if ret != ownerBalance {
		t.Errorf("Expect owner balance to be %v, got %v", ownerBalance, ret)
	}
	ret, err = token.GetBalance(crypto.AddressFromString(otherAddress))
	if ret != otherBalance {
		t.Errorf("Expect other balance to be %v, got %v", otherBalance, ret)
	}
	ret, err = token.GetBalance(crypto.AddressFromString(nonExistentAddress))
	if ret != 0 {
		t.Errorf("Expect non-existent balance to be %v, got %v", 0, ret)
	}
}

func TestTransferOK(t *testing.T) {
	token := setup()
	caller := crypto.AddressFromString(otherAddress)
	collector := crypto.AddressFromString(contractAddress)
	amount := uint64(100)

	events, err := token.Transfer(caller, collector, amount)
	if err != nil {
		panic(err)
	}
	if len(events) != 1 || events[0].Name != "Transfer" {
		t.Errorf("Expect %v transfer event, got %v", 0, len(events))
	}
	ret, err := token.GetBalance(caller)
	if ret != otherBalance-amount {
		t.Errorf("Expect caller balance to be %v, got %v", otherBalance-amount, ret)
	}
	ret, err = token.GetBalance(collector)
	if ret != contractBalance+amount {
		t.Errorf("Expect collector balance to be %v, got %v", contractBalance+amount, ret)
	}
}

func TestTransferFail(t *testing.T) {
	token := setup()
	caller := crypto.AddressFromString(nonExistentAddress)
	collector := crypto.AddressFromString(contractAddress)

	_, err := token.Transfer(caller, collector, 100)
	if err == nil || err.Error() != "Token transfer failed" {
		t.Errorf("Expect token transfer failed")
	}
}
