package test

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/db"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/QuoineFinancial/vertex/trie"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/QuoineFinancial/vertex/engine"
)

func TestVM(t *testing.T) {
	header, err := abi.LoadHeaderFromFile("../fixtures/header-event.json")
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("../data/token-event.wasm")
	if err != nil {
		panic(err)
	}
	contract := abi.Contract{
		Header: header,
		Code:   data,
	}
	contractBytes, _ := rlp.EncodeToBytes(&contract)
	caller := "LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"
	contractAddress := "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"
	database := db.NewMemoryDB()
	state, _ := storage.New(trie.Hash{}, database)
	accountState, _ := state.CreateAccount(crypto.AddressFromString(caller), crypto.AddressFromString(contractAddress), contractBytes)
	execEngine := engine.NewEngine(accountState, crypto.AddressFromString(caller))
	toAddress := "LB3EICIUKOUYCY4D7T2O6RKL7ISEPISNKUXNILDTJ76V2PDZVT5ZDP3U"
	var mint = 10000000
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
	ret, _ := execEngine.Ignite("get_balance", getBalanceTo)
	if int(*ret) != transfer {
		t.Errorf("Expect return value to be %v, got %v", transfer, *ret)
	}
	ret, _ = execEngine.Ignite("get_balance", getBalanceMint)
	if int(*ret) != mint-transfer {
		t.Errorf("Expect return value to be %v, got %v", mint-transfer, *ret)
	}
}
