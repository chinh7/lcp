package test

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strconv"
	"testing"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/QuoineFinancial/vertex/trie"

	"github.com/QuoineFinancial/vertex/engine"
)

func TestVM(t *testing.T) {
	var header abi.Header
	headerFile, err := ioutil.ReadFile("../fixtures/header2.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(headerFile, &header)

	data, err := ioutil.ReadFile("../data/token.wasm")
	if err != nil {
		panic(err)
	}
	encodedHeader, _ := header.Encode()
	contract := append(encodedHeader, data...)
	caller := "LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"
	contractAddress := "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"
	state := storage.GetState(trie.Hash{})
	accountState := state.CreateAccount(crypto.AddressFromString(caller), crypto.AddressFromString(contractAddress), &contract)
	execEngine := engine.NewEngine(accountState, crypto.AddressFromString(caller))
	toAddress := "LB3EICIUKOUYCY4D7T2O6RKL7ISEPISNKUXNILDTJ76V2PDZVT5ZDP3U"
	var mint = 100
	mintAmount := strconv.Itoa(mint)
	var transfer = 30
	transferAmount := strconv.Itoa(transfer)

	mintFunction, err := header.GetFunction("mint")
	if err != nil {
		panic(err)
	}
	mintArgs, err := abi.EncodeFromString(mintFunction.Parameters, []string{mintAmount})
	if err != nil {
		panic(err)
	}
	execEngine.Ignite("mint", mintArgs)

	transferFunction, err := header.GetFunction("transfer")
	if err != nil {
		panic(err)
	}
	transferArgs, err := abi.EncodeFromString(transferFunction.Parameters, []string{toAddress, transferAmount})
	if err != nil {
		panic(err)
	}
	execEngine.Ignite("transfer", transferArgs)

	getBalanceFunction, err := header.GetFunction("get_balance")
	if err != nil {
		panic(err)
	}
	getBalanceMint, _ := abi.EncodeFromString(getBalanceFunction.Parameters, []string{caller})
	if err != nil {
		panic(err)
	}

	getBalanceTo, err := abi.EncodeFromString(getBalanceFunction.Parameters, []string{toAddress})
	if err != nil {
		panic(err)
	}
	ret, _, _ := execEngine.Ignite("get_balance", getBalanceTo)
	value, ok := ret.(uint64)
	if !ok {
		t.Errorf("Expect return value to be uint32, got %s", reflect.TypeOf(ret))
	}
	if int(value) != transfer {
		t.Errorf("Expect return value to be %v, got %v", transfer, value)
	}
	ret, _, _ = execEngine.Ignite("get_balance", getBalanceMint)

	value, ok = ret.(uint64)
	if int(value) != mint-transfer {
		t.Errorf("Expect return value to be %v, got %v", mint-transfer, value)
	}
}
