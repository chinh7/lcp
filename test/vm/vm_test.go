package test

import (
	"encoding/binary"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"

	"github.com/QuoineFinancial/vertex/vm"
)

func TestVM(t *testing.T) {
	data, err := ioutil.ReadFile("../data/token.wasm")
	if err != nil {
		panic(err)
	}
	contractAddress := "LB3Z6N6HTFUPQ573QENJ4OCFFUPENY2EW7ZHQZSSIO4AODT3HHE53N52"
	state := storage.GetState()
	accountState := state.CreateAccountState(crypto.AddressFromString(contractAddress))
	accountState.SetCode(data)
	mintAddress := "LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O"
	toAddress := "LA7OPN4A3JNHLPHPEWM4PJDOYYDYNZOM7ES6YL3O7NC3PRY3V3UX6ANM"
	var mintAmount uint64 = 500
	var transferAmount uint64 = 321

	mintAmountBytes := make([]byte, 8)
	transferAmountBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(mintAmountBytes, mintAmount)
	binary.BigEndian.PutUint64(transferAmountBytes, transferAmount)

	vm.Call(accountState, "mint", []byte(mintAddress), mintAmountBytes)
	vm.Call(accountState, "transfer", []byte(mintAddress), []byte(toAddress), transferAmountBytes)
	ret, _ := vm.Call(accountState, "get_balance", []byte(toAddress))
	value, ok := ret.(uint32)
	if !ok {
		t.Error("Expect return value to be uint32, got {}", reflect.TypeOf(ret))
	}
	if uint64(value) != transferAmount {
		t.Error("Expect return value to be {}, got {}", mintAmount, value)
	}
	ret, _ = vm.Call(accountState, "get_balance", []byte(mintAddress))
	value, ok = ret.(uint32)
	if uint64(value) != mintAmount-transferAmount {
		t.Error("Expect return value to be {}, got {}", mintAmount-transferAmount, value)
	}
}
