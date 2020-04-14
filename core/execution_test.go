package core

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/QuoineFinancial/liquid-chain/trie"
	"github.com/QuoineFinancial/liquid-chain/util"
)

type TestResource struct {
	state *storage.State
	dbDir string
}

func NewTestResource() *TestResource {
	dbDir := "./testdata/db/test_" + strconv.Itoa(rand.Intn(10000))
	err := os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	stateDB := db.NewRocksDB(filepath.Join(dbDir, "storage.db"))
	state, err := storage.New(trie.Hash{}, stateDB)
	if err != nil {
		panic(err)
	}

	return &TestResource{state, dbDir}
}

func (tr *TestResource) cleanData() {
	err := os.RemoveAll(tr.dbDir)
	if err != nil {
		panic(err)
	}
}

func TestApplyTx(t *testing.T) {
	tr := NewTestResource()
	state := tr.state
	defer tr.cleanData()

	// Setup deploy contract transaction
	signer := crypto.TxSigner{Nonce: uint64(0)}
	data, err := util.BuildDeployTxData("./testdata/contract.wasm", "./testdata/contract-abi.json")
	if err != nil {
		panic(err)
	}
	deployTx := &crypto.Tx{From: signer, Data: data, To: crypto.Address{}, GasLimit: 0, GasPrice: 0}

	// Setup result events after deploy contract
	contractAddress := deployTx.From.CreateAddress()
	deployContractEvents := []event.Event{event.NewDeploymentEvent(contractAddress)}

	// Setup invoke contract transaction
	signer2 := crypto.TxSigner{Nonce: uint64(1)}
	invokeData, err := util.BuildInvokeTxData("./testdata/contract-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	invokeTx := &crypto.Tx{
		From:     signer2,
		To:       contractAddress,
		Data:     invokeData,
		GasLimit: 0,
		GasPrice: 0,
	}

	// Setup falsy tx to trigger reverse
	signer3 := crypto.TxSigner{Nonce: uint64(2)}
	invalidInvokeData, err := util.BuildInvokeTxData("./testdata/contract-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	invalidInvokeTx := &crypto.Tx{
		From:     signer3,
		To:       signer2.Address(), // Any valid address that is not a contract address
		Data:     invalidInvokeData,
		GasLimit: 0,
		GasPrice: 0,
	}

	type args struct {
		state      *storage.State
		tx         *crypto.Tx
		gasStation gas.Station
	}
	tests := []struct {
		name       string
		args       args
		result     uint64
		events     []event.Event
		gasUsed    uint64
		wantErr    bool
		wantErrObj error
	}{
		{
			name:       "invalid deploy tx, out of gas",
			args:       args{state, deployTx, gas.NewLiquidStation(nil, crypto.Address{})},
			result:     0,
			events:     nil,
			gasUsed:    0,
			wantErr:    true,
			wantErrObj: errors.New("out of gas"),
		},
		{
			name:       "valid deploy tx",
			args:       args{state, deployTx, gas.NewFreeStation(nil)},
			result:     0,
			events:     deployContractEvents,
			gasUsed:    0,
			wantErr:    false,
			wantErrObj: nil,
		},
		{
			name:       "valid invoke tx",
			args:       args{state, invokeTx, gas.NewFreeStation(nil)},
			result:     0,
			events:     make([]event.Event, 0),
			gasUsed:    0,
			wantErr:    false,
			wantErrObj: nil,
		},
		{
			name:       "invalid invoke tx, reverse",
			args:       args{state, invalidInvokeTx, gas.NewFreeStation(nil)},
			result:     0,
			events:     make([]event.Event, 0),
			gasUsed:    0,
			wantErr:    true,
			wantErrObj: errors.New("abi: cannot decode empty contract"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, events, gasUsed, err := ApplyTx(tt.args.state, tt.args.tx, tt.args.gasStation)
			if tt.wantErr && (err == nil) {
				t.Errorf("%s: applyTx() error = %v, wantErr %v", tt.name, err, tt.wantErrObj.Error())
			}
			if tt.wantErr && (err != nil) {
				if tt.wantErrObj.Error() != err.Error() {
					t.Errorf("%s: applyTx() error = %v, wantErr %v", tt.name, err, tt.wantErrObj.Error())
				}
			}
			if result != tt.result {
				t.Errorf("%s: applyTx() result = %v, want %v", tt.name, result, tt.result)
			}
			if !reflect.DeepEqual(events, tt.events) {
				t.Errorf("%s: applyTx() events = %v, want %v", tt.name, events, tt.events)
			}
			if gasUsed != tt.gasUsed {
				t.Errorf("%s: applyTx() gasUsed = %v, want %v", tt.name, gasUsed, tt.gasUsed)
			}
		})
	}
}
