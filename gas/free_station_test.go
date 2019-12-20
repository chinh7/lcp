package gas

import (
	"testing"

	"github.com/QuoineFinancial/liquid-chain/crypto"
)

const contractAddress = "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"
const otherAddress = "LCR57ROUHIQ2AV4D3E3D7ZBTR6YXMKZQWTI4KSHSWCUCRXBKNJKKBCNY"

type MockFreeApp struct {
	App
}

func (app *MockFreeApp) SetGasStation(station Station) {
	panic("Should not be call")
}

func (app *MockFreeApp) GetGasContractToken() Token {
	return &MockToken{}
}

func TestFreeSufficient(t *testing.T) {
	app := &MockFreeApp{}
	station := NewFreeStation(app)
	ret := station.Sufficient(crypto.AddressFromString(otherAddress), 10)

	if !ret {
		t.Error("Expected return true")
	}
}

func TestFreeBurn(t *testing.T) {
	app := &MockFreeApp{}
	station := NewFreeStation(app)

	station.Burn(crypto.AddressFromString(otherAddress), 10)

	ret := station.Burn(crypto.AddressFromString(otherAddress), 0)
	if ret != nil {
		t.Error("Expected return nil")
	}
}
