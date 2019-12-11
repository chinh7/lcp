package test

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/db"
	"github.com/QuoineFinancial/vertex/gas"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/QuoineFinancial/vertex/trie"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/QuoineFinancial/vertex/engine"
)

func TestVM(t *testing.T) {
	contract := loadContract("../fixtures/header-event.json", "../data/token-event.wasm")
	header := contract.Header
	contractBytes, _ := rlp.EncodeToBytes(&contract)
	caller := "LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R"
	contractAddress := "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"
	database := db.NewMemoryDB()
	state, _ := storage.New(trie.Hash{}, database)
	accountState, _ := state.CreateAccount(crypto.AddressFromString(caller), crypto.AddressFromString(contractAddress), contractBytes)
	execEngine := engine.NewEngine(state, accountState, crypto.AddressFromString(caller), &gas.FreePolicy{}, 0)
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
	_, err = execEngine.Ignite("mint", mintArgs)
	if err != nil {
		panic(err)
	}

	transferFunction, err := header.GetFunction("transfer")
	if err != nil {
		panic(err)
	}
	transferArgs, err := abi.EncodeFromString(transferFunction.Parameters, []string{toAddress, transferAmount})
	if err != nil {
		panic(err)
	}
	_, err = execEngine.Ignite("transfer", transferArgs)
	if err != nil {
		panic(err)
	}

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
	ret, err := execEngine.Ignite("get_balance", getBalanceTo)
	if err != nil {
		t.Error(err)
	}
	if int(ret) != transfer {
		t.Errorf("Expect return value to be %v, got %v", transfer, ret)
	}
	ret, err = execEngine.Ignite("get_balance", getBalanceMint)
	if err != nil {
		t.Error(err)
	}
	if int(ret) != mint-transfer {
		t.Errorf("Expect return value to be %v, got %v", mint-transfer, ret)
	}
}

func TestChainedInvoke(t *testing.T) {
	mathContract := loadContract("../fixtures/math-abi.json", "../data/math.wasm")
	mathBytes, _ := rlp.EncodeToBytes(&mathContract)
	caller := crypto.AddressFromString("LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R")
	mathAddress := crypto.AddressFromString("LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53")
	database := db.NewMemoryDB()
	state, _ := storage.New(trie.Hash{}, database)
	state.CreateAccount(caller, mathAddress, mathBytes)

	utilContract := loadContract("../fixtures/util-abi.json", "../data/util.wasm")
	utilBytes, _ := rlp.EncodeToBytes(&utilContract)
	utilAddress := crypto.AddressFromString("LCR57ROUHIQ2AV4D3E3D7ZBTR6YXMKZQWTI4KSHSWCUCRXBKNJKKBCNY")
	utilAccount, _ := state.CreateAccount(caller, utilAddress, utilBytes)
	execEngine := engine.NewEngine(state, utilAccount, caller, &gas.FreePolicy{}, 0)

	funcName := "init"
	function, err := utilContract.Header.GetFunction(funcName)
	if err != nil {
		panic(err)
	}
	args, err := abi.EncodeFromString(function.Parameters, []string{mathAddress.String()})
	if err != nil {
		panic(err)
	}
	_, err = execEngine.Ignite(funcName, args)
	if err != nil {
		panic(err)
	}

	funcName = "variance"
	function, err = utilContract.Header.GetFunction(funcName)
	if err != nil {
		panic(err)
	}
	args, err = abi.EncodeFromString(function.Parameters, []string{"[1,2,3,4,5]"})
	if err != nil {
		panic(err)
	}
	ret, err := execEngine.Ignite(funcName, args)
	if err != nil {
		panic(err)
	}
	if int32(ret) != 2 {
		t.Errorf("Expect return value to be %v, got %v", 2, int32(ret))
	}
}

func TestChainedInvokeOverflow(t *testing.T) {
	caller := crypto.AddressFromString("LDH4MEPOJX3EGN3BLBTLEYXVHYCN3AVA7IOE772F3XGI6VNZHAP6GX5R")
	database := db.NewMemoryDB()
	state, _ := storage.New(trie.Hash{}, database)

	utilContract := loadContract("../fixtures/util-abi.json", "../data/util.wasm")
	utilBytes, _ := rlp.EncodeToBytes(&utilContract)
	utilAddress := crypto.AddressFromString("LCR57ROUHIQ2AV4D3E3D7ZBTR6YXMKZQWTI4KSHSWCUCRXBKNJKKBCNY")
	utilAccount, _ := state.CreateAccount(caller, utilAddress, utilBytes)

	execEngine := engine.NewEngine(state, utilAccount, caller, &gas.FreePolicy{}, 0)

	funcName := "init"
	function, err := utilContract.Header.GetFunction(funcName)
	if err != nil {
		panic(err)
	}
	args, err := abi.EncodeFromString(function.Parameters, []string{utilAddress.String()})
	if err != nil {
		panic(err)
	}
	_, err = execEngine.Ignite(funcName, args)
	if err != nil {
		panic(err)
	}

	funcName = "mean"
	function, err = utilContract.Header.GetFunction(funcName)
	if err != nil {
		panic(err)
	}
	args, err = abi.EncodeFromString(function.Parameters, []string{"[1,2,3,4,5]"})
	if err != nil {
		panic(err)
	}
	_, err = execEngine.Ignite(funcName, args)
	if err == nil || err.Error() != "call depth limit reached" {
		t.Errorf("Unexpected error %v", err)
	}
}

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
