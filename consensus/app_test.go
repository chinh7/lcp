package consensus

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/QuoineFinancial/liquid-chain/token"
	"github.com/QuoineFinancial/liquid-chain/trie"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
)

func ensureDir(path string) error {
	newpath := filepath.Join("./testdata/db/", path)
	return os.MkdirAll(newpath, os.ModePerm)
}

func removeDir(foldername string) error {
	newpath := filepath.Join("./testdata/db/", foldername)
	return os.RemoveAll(newpath)
}

type TestConfig struct {
	app       *App
	dbDirname string
}

func NewTestConfig() *TestConfig {
	dbDirname := "test_" + strconv.Itoa(rand.Intn(10000))
	err := ensureDir(dbDirname)
	if err != nil {
		panic(err)
	}
	app := NewApp("testapp", "./testdata/db/"+dbDirname+"/", "LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7")
	return &TestConfig{app, dbDirname}
}

func (tc *TestConfig) CleanData() {
	err := removeDir(tc.dbDirname)
	if err != nil {
		panic(err)
	}
}

func TestNewApp(t *testing.T) {
	type args struct {
		nodeInfo           string
		dbDir              string
		gasContractAddress string
	}
	tests := []struct {
		name string
		args args
		want *App
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewApp(tt.args.nodeInfo, tt.args.dbDir, tt.args.gasContractAddress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_BeginBlock(t *testing.T) {
	t.Run("Should load state", func(t *testing.T) {
		tc := NewTestConfig()
		defer tc.CleanData()
		app := tc.app

		reqHeight := int64(1)
		reqAppHash := trie.Hash{}
		req := types.RequestBeginBlock{Header: types.Header{Height: reqHeight, AppHash: reqAppHash[:]}}
		got := app.BeginBlock(req)
		want := types.ResponseBeginBlock{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("App.BeginBlock() = %v, want %v", got, want)
		}

		// loadState() should be called
		assert.NotNil(t, app.state)
		assert.Equal(t, app.state.BlockInfo.Height, uint64(reqHeight))
		assert.Equal(t, app.state.BlockInfo.AppHash, reqAppHash)
	})
}

func TestApp_Info(t *testing.T) {
	tc := NewTestConfig()
	defer tc.CleanData()

	t.Run("Should return valid response", func(t *testing.T) {
		app := tc.app
		blockInfo := &storage.BlockInfo{Height: 1, AppHash: trie.Hash{}, Time: time.Now()}

		app.loadState(blockInfo)
		got := app.Info(types.RequestInfo{})
		want := types.ResponseInfo{
			Data:             "{\"version\":testapp}",
			LastBlockHeight:  int64(blockInfo.Height),
			LastBlockAppHash: blockInfo.AppHash[:],
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got app.Info() = %v, want %v", got, want)
		}
	})
}

func TestApp_CheckTx(t *testing.T) {
	tc := NewTestConfig()
	defer tc.CleanData()
	app := tc.app
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: trie.Hash{}, Time: time.Now()}
	app.loadState(blockInfo)

	t.Run("Deserialize tx error", func(t *testing.T) {
		invalidTxBytes, err := ioutil.ReadFile("./testdata/invalid_tx.dat")
		if err != nil {
			panic(err)
		}

		got := app.CheckTx(types.RequestCheckTx{Tx: invalidTxBytes})
		want := types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  "rlp: expected input list for crypto.Tx",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("App.CheckTx() = %v, want %v", got, want)
		}
	})

	t.Run("Invalid tx", func(t *testing.T) {
		invalidNonceTxBytes, err := ioutil.ReadFile("./testdata/invalid_nonce_tx.dat")
		if err != nil {
			panic(err)
		}

		got := app.CheckTx(types.RequestCheckTx{Tx: invalidNonceTxBytes})
		want := types.ResponseCheckTx{
			Code: code.CodeTypeBadNonce,
			Log:  "Invalid nonce. Expected 0, got 10",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("App.CheckTx() = %v, want %v", got, want)
		}
	})

	t.Run("Valid tx", func(t *testing.T) {
		txBytes, err := ioutil.ReadFile("./testdata/deploy_contract_tx.dat")
		if err != nil {
			panic(err)
		}

		got := app.CheckTx(types.RequestCheckTx{Tx: txBytes})
		want := types.ResponseCheckTx{Code: code.CodeTypeOK}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("App.CheckTx() = %v, want %v", got, want)
		}
	})
}

func TestApp_validateTx(t *testing.T) {
	tc := NewTestConfig()
	defer tc.CleanData()
	app := tc.app
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: trie.Hash{}, Time: time.Now()}
	app.loadState(blockInfo)

	contractHex, err := ioutil.ReadFile("./testdata/contract_hex.txt")
	if err != nil {
		panic(err)
	}
	address, _ := crypto.AddressFromString("LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7")
	app.SetGasStation(gas.NewLiquidStation(app, address))
	_, err = app.state.CreateAccount(address, address, contractHex)
	if err != nil {
		panic(err)
	}
	app.state.Commit()

	invalidNonceTxBytes, err := ioutil.ReadFile("./testdata/invalid_nonce_tx.dat")
	if err != nil {
		panic(err)
	}
	invalidNonceTx := &crypto.Tx{}
	err = invalidNonceTx.Deserialize(invalidNonceTxBytes)
	if err != nil {
		panic(err)
	}

	invalidSigTxBytes, err := ioutil.ReadFile("./testdata/invalid_sig_tx.dat")
	if err != nil {
		panic(err)
	}
	invalidSigTx := &crypto.Tx{}
	err = invalidSigTx.Deserialize(invalidSigTxBytes)
	// Change tx so that sigHash is not the same as origin
	invalidSigTx.GasLimit = 10
	if err != nil {
		panic(err)
	}

	txBytes, err := ioutil.ReadFile("./testdata/invoke_contract_tx.dat")
	if err != nil {
		panic(err)
	}
	tx := &crypto.Tx{}
	err = tx.Deserialize(txBytes)
	if err != nil {
		panic(err)
	}
	// non-existent contract
	txByte, err := ioutil.ReadFile("./testdata/non_existent_contract.json")
	if err != nil {
		panic(err)
	}
	nonExistentTx := &crypto.Tx{}
	_ = json.Unmarshal([]byte(txByte), nonExistentTx)

	type args struct {
		tx     *crypto.Tx
		txSize int
	}
	tests := []struct {
		name    string
		app     *App
		args    args
		want    uint32
		wantErr error
	}{
		{
			"Invalid size",
			app,
			args{&crypto.Tx{}, 1024*1024 + 1},
			code.CodeTypeUnknownError,
			errors.New("Transaction size exceed 1048576B"),
		}, {
			"Invalid nonce",
			app,
			args{invalidNonceTx, len(invalidNonceTxBytes)},
			code.CodeTypeBadNonce,
			errors.New("Invalid nonce. Expected 0, got 10"),
		}, {
			"Invalid signature",
			app,
			args{invalidSigTx, len(invalidSigTxBytes)},
			code.CodeTypeUnknownError,
			errors.New("Invalid signature"),
			// }, {
			// 	"Insufficient fee",
			// 	app,
			// 	args{tx, len(txBytes)},
			// 	code.CodeTypeUnknownError,
			// 	errors.New("Insufficient fee"),
		}, {
			"non-existent contract",
			app,
			args{nonExistentTx, len(txByte)},
			code.CodeTypeUnknownError,
			errors.New("contract not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tc.app
			got, err := app.validateTx(tt.args.tx, tt.args.txSize)
			if err != nil && tt.wantErr.Error() != err.Error() {
				t.Errorf("App.validateTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("App.validateTx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_DeliverTx(t *testing.T) {
	tc := NewTestConfig()
	defer tc.CleanData()
	app := tc.app
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: trie.Hash{}, Time: time.Now()}
	app.loadState(blockInfo)

	t.Run("Deserialize tx error", func(t *testing.T) {
		invalidTxBytes, err := ioutil.ReadFile("./testdata/invalid_tx.dat")
		if err != nil {
			panic(err)
		}

		got := app.DeliverTx(types.RequestDeliverTx{Tx: invalidTxBytes})
		want := types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  "rlp: expected input list for crypto.Tx",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("App.DeliverTx() = %v, want %v", got, want)
		}
	})

	t.Run("Invalid tx", func(t *testing.T) {
		invalidNonceTxBytes, err := ioutil.ReadFile("./testdata/invalid_nonce_tx.dat")
		if err != nil {
			panic(err)
		}

		got := app.DeliverTx(types.RequestDeliverTx{Tx: invalidNonceTxBytes})
		want := types.ResponseDeliverTx{
			Code: code.CodeTypeBadNonce,
			Log:  "Invalid nonce. Expected 0, got 10",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("App.DeliverTx() = %v, want %v", got, want)
		}
	})

	t.Run("Deploy transaction", func(t *testing.T) {
		txBytes, err := ioutil.ReadFile("./testdata/deploy_contract_tx.dat")
		if err != nil {
			panic(err)
		}

		req := types.RequestDeliverTx{Tx: txBytes}
		tx := &crypto.Tx{}
		err = tx.Deserialize(txBytes)
		if err != nil {
			panic(err)
		}
		detailEvent := event.NewDetailsEvent(1, tx.From.Address(), tx.To, tx.From.Nonce, 0)
		got := app.DeliverTx(req)

		// Detail event must equal
		assert.Equal(t, detailEvent.ToTMEvent(), got.Events[len(got.Events)-1])
		assert.Equal(t, code.CodeTypeOK, got.Code)
		assert.Equal(t, "ok", got.Info)
		assert.Equal(t, int64(tx.GasLimit), got.GasWanted)
		assert.Equal(t, int64(0), got.GasUsed)
	})

	t.Run("Invoke transaction", func(t *testing.T) {
		txBytes, err := ioutil.ReadFile("./testdata/invoke_contract_tx.dat")
		if err != nil {
			panic(err)
		}

		req := types.RequestDeliverTx{Tx: txBytes}
		tx := &crypto.Tx{}
		err = tx.Deserialize(txBytes)
		if err != nil {
			panic(err)
		}
		detailEvent := event.NewDetailsEvent(1, tx.From.Address(), tx.To, tx.From.Nonce, 0)
		got := app.DeliverTx(req)

		// Detail event must equal
		assert.Equal(t, detailEvent.ToTMEvent(), got.Events[len(got.Events)-1])
		assert.Equal(t, code.CodeTypeOK, got.Code)
		assert.Equal(t, "ok", got.Info)
		assert.Equal(t, int64(tx.GasLimit), got.GasWanted)
		assert.Equal(t, int64(0), got.GasUsed)
	})
}

func TestApp_Commit(t *testing.T) {
	tc := NewTestConfig()
	defer tc.CleanData()
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: trie.Hash{}, Time: time.Now()}
	tc.app.loadState(blockInfo)
	got := tc.app.Commit()
	want := types.ResponseCommit{Data: []byte{69, 176, 207, 194, 32, 206, 236, 91, 124, 28, 98, 196, 212, 25, 61, 56, 228, 235, 164, 142, 136, 21, 114, 156, 231, 95, 156, 10, 176, 228, 193, 192}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("App.Commit() = %v, want %v", got, want)
	}
}

func TestApp_GetGasContractToken(t *testing.T) {
	tc := NewTestConfig()
	defer tc.CleanData()
	a := tc.app
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: trie.Hash{}, Time: time.Now()}
	a.loadState(blockInfo)

	tc2 := NewTestConfig()
	defer tc.CleanData()
	contractHex, err := ioutil.ReadFile("./testdata/contract_hex.txt")
	if err != nil {
		panic(err)
	}
	address, _ := crypto.AddressFromString("LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7")
	a2 := tc2.app
	a2.loadState(blockInfo)
	account, err := a2.state.CreateAccount(address, address, contractHex)
	_ = a2.state.Commit()
	if err != nil {
		panic(err)
	}

	token := token.NewToken(a2.state, account)
	tests := []struct {
		name string
		app  *App
		want gas.Token
	}{
		{"gasContractAddress not exist", &App{}, nil},
		{"gasContractAddress exist, contract not exist", a, nil},
		{"gasContractAddress exist, contract exist", a2, token},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tt.app
			if got := app.GetGasContractToken(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.GetGasContractToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
