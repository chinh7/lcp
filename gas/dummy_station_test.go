package gas

import (
	cryptoRand "crypto/rand"
	"reflect"
	"testing"

	"golang.org/x/crypto/ed25519"

	"github.com/QuoineFinancial/liquid-chain/crypto"
)

func TestDummyStation_Sufficient(t *testing.T) {
	t.Run("Sufficient", func(t *testing.T) {
		station := &DummyStation{app: nil, policy: &FreePolicy{}}
		pub, _, _ := ed25519.GenerateKey(cryptoRand.Reader)
		toAddr := crypto.AddressFromPubKey(pub)
		want := false

		if got := station.Sufficient(toAddr, uint64(0)); got != want {
			t.Errorf("DummyStation.Sufficient() = %v, want %v", got, want)
		}
	})

	t.Run("Insufficient", func(t *testing.T) {
		pub, _, _ := ed25519.GenerateKey(cryptoRand.Reader)
		station := &DummyStation{app: nil, policy: &FreePolicy{}}
		toAddr := crypto.AddressFromPubKey(pub)
		want := true

		if got := station.Sufficient(toAddr, uint64(1)); got != want {
			t.Errorf("DummyStation.Sufficient() = %v, want %v", got, want)
		}
	})
}

func TestDummyStation_Burn(t *testing.T) {
	station := &DummyStation{
		app:    nil,
		policy: &FreePolicy{},
	}
	pub, _, _ := ed25519.GenerateKey(cryptoRand.Reader)
	addr := crypto.AddressFromPubKey(pub)
	var want []*crypto.TxEvent
	if got := station.Burn(addr, uint64(0)); !reflect.DeepEqual(got, want) {
		t.Errorf("DummyStation.Burn() = %v, want %v", got, want)
	}
}

func TestDummyStation_Switch(t *testing.T) {
	station := &DummyStation{
		app:    nil,
		policy: &FreePolicy{},
	}
	want := false
	if got := station.Switch(); got != want {
		t.Errorf("DummyStation.Switch() = %v, want %v", got, want)
	}
}

func TestDummyStation_GetPolicy(t *testing.T) {
	station := &DummyStation{
		app:    nil,
		policy: &FreePolicy{},
	}
	want := &FreePolicy{}
	if got := station.GetPolicy(); !reflect.DeepEqual(got, want) {
		t.Errorf("DummyStation.GetPolicy() = %v, want %v", got, want)
	}
}

func TestDummyStation_CheckGasPrice(t *testing.T) {
	station := &DummyStation{
		app:    nil,
		policy: &FreePolicy{},
	}
	want := false
	if got := station.CheckGasPrice(uint32(0)); got != want {
		t.Errorf("DummyStation.CheckGasPrice() = %v, want %v", got, want)
	}
}

type DummyApp struct{}

func (app *DummyApp) SetGasStation(station Station) {}

func (app *DummyApp) GetGasContractToken() Token {
	return nil
}

func TestNewDummyStation(t *testing.T) {
	app := &DummyApp{}
	want := &DummyStation{
		app:    app,
		policy: &FreePolicy{},
	}

	if got := NewDummyStation(app); !reflect.DeepEqual(got, want) {
		t.Errorf("NewTestStation() = %v, want %v", got, want)
	}
}
