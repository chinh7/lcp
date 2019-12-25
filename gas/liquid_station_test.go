package gas

import (
	"errors"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/storage"
)

type MockToken struct {
	Token
}

func (token *MockToken) GetBalance(addr crypto.Address) (uint64, error) {
	return 100, nil
}

func (token *MockToken) Transfer(caller crypto.Address, addr crypto.Address, amount uint64) ([]event.Event, error) {
	if addr.String() != contractAddressStr {
		panic("Expected collector is gas contract address")
	}
	if amount == 10000 {
		return nil, errors.New("Token transfer failed")
	}
	return []event.Event{}, nil
}

func (token *MockToken) GetContract() *storage.Account {
	return nil
}

type MockApp struct {
	App
}

func (app *MockApp) SetGasStation(station Station) {
	panic("Should not be call")
}

func (app *MockApp) GetGasContractToken() Token {
	return &MockToken{}
}

func TestSwitch(t *testing.T) {
	app := &MockApp{}
	contractAddress, _ := crypto.AddressFromString(contractAddressStr)
	station := NewLiquidStation(app, contractAddress)
	ret := station.Switch()
	if ret {
		t.Error("Expected return false")
	}
}

func TestSufficient(t *testing.T) {
	app := &MockApp{}
	contractAddress, _ := crypto.AddressFromString(contractAddressStr)
	otherAddress, _ := crypto.AddressFromString(otherAddressStr)

	station := NewLiquidStation(app, contractAddress)
	ret := station.Sufficient(otherAddress, 10)

	if !ret {
		t.Error("Expected return true")
	}

	ret = station.Sufficient(otherAddress, 1000)

	if ret {
		t.Error("Expected return false")
	}
}

func TestBurn(t *testing.T) {
	app := &MockApp{}
	contractAddress, _ := crypto.AddressFromString(contractAddressStr)
	otherAddress, _ := crypto.AddressFromString(otherAddressStr)

	station := NewLiquidStation(app, contractAddress)

	station.Burn(otherAddress, 10)

	ret := station.Burn(otherAddress, 0)
	if ret != nil {
		t.Error("Expected return nil")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	station.Burn(otherAddress, 10000)
}
