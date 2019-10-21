package test

import (
	"encoding/binary"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/QuoineFinancial/vertex/trie"

	"github.com/QuoineFinancial/vertex/engine"
)

func TestVM(t *testing.T) {
	data, err := ioutil.ReadFile("../data/token.wasm")
	if err != nil {
		panic(err)
	}
	contractAddress := "LB3Z6N6HTFUPQ573QENJ4OCFFUPENY2EW7ZHQZSSIO4AODT3HHE53N52"
	state := storage.GetState(trie.Hash{})
	accountState := state.CreateAccount(crypto.AddressFromString(contractAddress), &data)
	engine := engine.NewEngine(accountState)
	mintAddress := "LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O"
	toAddress := "LA7OPN4A3JNHLPHPEWM4PJDOYYDYNZOM7ES6YL3O7NC3PRY3V3UX6ANM"
	var mintAmount uint64 = 500
	var transferAmount uint64 = 321

	mintAmountBytes := make([]byte, 8)
	transferAmountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(mintAmountBytes, mintAmount)
	binary.BigEndian.PutUint64(transferAmountBytes, transferAmount)
	engine.Ignite("mint", []byte(mintAddress), mintAmountBytes)
	engine.Ignite("transfer", []byte(mintAddress), []byte(toAddress), transferAmountBytes)
	ret, _, _ := engine.Ignite("get_balance", []byte(toAddress))
	value, ok := ret.(uint64)
	if !ok {
		t.Error("Expect return value to be uint64, got {}", reflect.TypeOf(ret))
	}
	if value != transferAmount {
		t.Errorf("Expect return value to be %d, got %d", transferAmount, value)
	}
	ret, _, _ = engine.Ignite("get_balance", []byte(mintAddress))
	value, ok = ret.(uint64)
	if uint64(value) != mintAmount-transferAmount {
		t.Error("Expect return value to be {}, got {}", mintAmount-transferAmount, value)
	}
}
