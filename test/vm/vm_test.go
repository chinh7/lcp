package test

import (
	"io/ioutil"
	"reflect"
	"strconv"
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
	var mintAmount int64 = 500
	var transferAmount int64 = 321
	mintAmountStr := strconv.FormatInt(mintAmount, 10)
	transferAmountStr := strconv.FormatInt(transferAmount, 10)

	vm.Call(accountState, "mint", []byte(mintAddress), []byte(mintAmountStr))
	vm.Call(accountState, "transfer", []byte(mintAddress), []byte(toAddress), []byte(transferAmountStr))
	ret, _ := vm.Call(accountState, "get_balance", []byte(toAddress))
	value, ok := ret.(uint32)
	if !ok {
		t.Error("Expect return value to be uint32, got {}", reflect.TypeOf(ret))
	}
	if int64(value) != transferAmount {
		t.Error("Expect return value to be {}, got {}", mintAmount, value)
	}
	ret, _ = vm.Call(accountState, "get_balance", []byte(mintAddress))
	value, ok = ret.(uint32)
	if int64(value) != mintAmount-transferAmount {
		t.Error("Expect return value to be {}, got {}", mintAmount-transferAmount, value)
	}
}
