package consensus

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/stretchr/testify/assert"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/QuoineFinancial/liquid-chain/token"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"

	"golang.org/x/crypto/ed25519"
)

const SEED = "0c61093a4983f5ba8cf83939efc6719e0c61093a4983f5ba8cf83939efc6719e"

type TestResource struct {
	app   *App
	dbDir string
}

func NewTestResource() *TestResource {
	rand.Seed(time.Now().UTC().UnixNano())
	dbDir := "./testdata/db/test_" + strconv.Itoa(rand.Intn(10000)) + "/"
	err := os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	app := NewApp("testapp", dbDir, "")
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: common.Hash{}, Time: time.Now()}
	app.loadState(blockInfo)

	// Manually deploy contract
	address, err := crypto.AddressFromString("LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7")
	if err != nil {
		panic(err)
	}
	// This file can be construct using liquid-chain-js
	contract, err := ioutil.ReadFile("./testdata/deploy-contract.dat")
	if err != nil {
		panic(err)
	}
	_, err = app.state.CreateAccount(address, address, contract)
	if err != nil {
		panic(err)
	}
	_ = app.state.Commit()
	app.gasContractAddress = "LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7"

	return &TestResource{app, dbDir}
}

func (tc *TestResource) CleanData() {
	err := os.RemoveAll(tc.dbDir)
	if err != nil {
		panic(err)
	}
}

func loadPrivateKey(SEED string) ed25519.PrivateKey {
	hexSeed, err := hex.DecodeString(SEED)
	if err != nil {
		panic(err)
	}
	return ed25519.NewKeyFromSeed(hexSeed)
}

func TestNewApp(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	dbDir := "./testdata/db/test_" + strconv.Itoa(rand.Intn(10000)) + "/"
	err := os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.RemoveAll(dbDir)
	}()

	blockInfo := &storage.BlockInfo{
		Height:  uint64(1),
		AppHash: common.Hash{},
		Time:    time.Now(),
	}
	bytes, _ := rlp.EncodeToBytes(blockInfo)
	infoDB := db.NewRocksDB(filepath.Join(dbDir, "info.db"))
	infoDB.Put([]byte("lastBlockInfo"), bytes)
	infoDB.Close()

	app := NewApp("testapp", dbDir, "")
	assert.NotNil(t, app)
}

func TestApp_BeginBlock(t *testing.T) {
	t.Run("Should load state", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()
		app := tr.app
		appHash := tr.app.state.Commit()
		reqHeight := int64(0)
		req := types.RequestBeginBlock{Header: types.Header{Height: reqHeight, AppHash: appHash.Bytes()}}
		got := app.BeginBlock(req)
		want := types.ResponseBeginBlock{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("App.BeginBlock() = %v, want %v", got, want)
		}

		// loadState() should be called
		assert.NotNil(t, app.state)
		assert.Equal(t, app.state.BlockInfo.Height, uint64(reqHeight))
		assert.Equal(t, app.state.BlockInfo.AppHash, appHash)
	})
}

func TestApp_Info(t *testing.T) {
	tr := NewTestResource()
	defer tr.CleanData()

	t.Run("Should return valid response", func(t *testing.T) {
		app := tr.app
		appHash := tr.app.state.Commit()
		blockInfo := &storage.BlockInfo{Height: 2, AppHash: appHash, Time: time.Now()}
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
	tr := NewTestResource()
	defer tr.CleanData()
	app := tr.app
	appHash := tr.app.state.Commit()
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: appHash, Time: time.Now()}
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
		pubkey, prvkey, err := ed25519.GenerateKey(cryptoRand.Reader)
		if err != nil {
			panic(err)
		}
		txData := crypto.TxData{}
		tx := &crypto.Tx{From: crypto.TxSigner{PubKey: pubkey}, Data: txData.Serialize(), GasPrice: uint32(18)}
		err = tx.Sign(prvkey)
		if err != nil {
			panic(err)
		}
		got := app.CheckTx(types.RequestCheckTx{Tx: tx.Serialize()})
		want := types.ResponseCheckTx{Code: code.CodeTypeOK}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("App.CheckTx() = %v, want %v", got, want)
		}
	})
}

func TestApp_validateTx(t *testing.T) {
	// invalid tx size
	t.Run("invalid size", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()
		tx := &crypto.Tx{}
		txSize := 1024*1024 + 1 // set to large number to produce error
		wantErr := errors.New("Transaction size exceed 1048576B")
		wantErrCode := CodeTypeExceedTransactionSize

		got, err := tr.app.validateTx(tx, txSize)
		if err == nil || wantErr.Error() != err.Error() {
			t.Errorf("App.validateTx() error = %v, wantErr %v", err, wantErr)
			return
		}
		if got != wantErrCode {
			t.Errorf("App.validateTx() code = %v, want %v", got, wantErrCode)
		}
	})

	// invalid tx nonce
	t.Run("invalid nonce", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()
		signer := crypto.TxSigner{Nonce: uint64(10)} // Set nonce to any number but not 0 to trigger error
		tx := &crypto.Tx{Data: nil, From: signer, GasLimit: 1, GasPrice: 1}
		wantErr := errors.New("Invalid nonce. Expected 0, got 10")
		wantErrCode := CodeTypeBadNonce

		got, err := tr.app.validateTx(tx, len(tx.Serialize()))
		if err == nil || wantErr.Error() != err.Error() {
			t.Errorf("App.validateTx() error = %v, wantErr %v", err, wantErr)
			return
		}
		if got != wantErrCode {
			t.Errorf("App.validateTx() code = %v, want %v", got, wantErrCode)
		}
	})

	// invalid public key
	t.Run("invalid public key", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()
		signer := crypto.TxSigner{Nonce: uint64(0)} // add signer without signature to trigger error
		tx := &crypto.Tx{Data: nil, From: signer, GasLimit: 1, GasPrice: 1}
		wantErr := errors.New("Invalid public key. Expected size of 32B, got 0B")
		wantErrCode := CodeTypeInvalidPubKey

		got, err := tr.app.validateTx(tx, len(tx.Serialize()))
		if err == nil || wantErr.Error() != err.Error() {
			t.Errorf("App.validateTx() error = %v, wantErr %v", err, wantErr)
			return
		}
		if got != wantErrCode {
			t.Errorf("App.validateTx() code = %v, want %v", got, wantErrCode)
		}
	})

	// invalid signature
	t.Run("invalid signature", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()
		signer := crypto.TxSigner{Nonce: uint64(0)}          // add signer without signature to trigger error
		_, priv, _ := ed25519.GenerateKey(cryptoRand.Reader) // generate random privkey to sign
		tx := &crypto.Tx{Data: nil, From: signer, GasLimit: 1, GasPrice: 1}
		// This step mainly to populate valid publickey
		if err := tx.Sign(priv); err != nil {
			panic(err)
		}
		tx.From.Signature = []byte{} // Remove valid private key to trigger error
		wantErr := errors.New("Invalid signature")
		wantErrCode := CodeTypeInvalidSignature

		got, err := tr.app.validateTx(tx, len(tx.Serialize()))
		if err == nil || wantErr.Error() != err.Error() {
			t.Errorf("App.validateTx() error = %v, wantErr %v", err, wantErr)
			return
		}
		if got != wantErrCode {
			t.Errorf("App.validateTx() code = %v, want %v", got, wantErrCode)
		}
	})

	// insufficient fee
	t.Run("insufficient fee", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()

		tr.app.SetGasStation(gas.NewDummyStation(tr.app)) // Set to dummy station to trigger insufficient fee check
		signer := crypto.TxSigner{Nonce: uint64(0)}
		tx := &crypto.Tx{Data: nil, From: signer, GasLimit: 0, GasPrice: 0}
		privKey := loadPrivateKey(SEED)
		if err := tx.Sign(privKey); err != nil {
			panic(err)
		}
		wantErr := errors.New("Insufficient fee")
		wantErrCode := CodeTypeInsufficientFee

		got, err := tr.app.validateTx(tx, len(tx.Serialize()))
		if got != wantErrCode {
			t.Errorf("App.validateTx() code = %v, want %v", got, wantErrCode)
		}
		if err == nil || wantErr.Error() != err.Error() {
			t.Errorf("App.validateTx() error = %+v, wantErr %v", err, wantErr)
			return
		}
	})

	// invalid gas fee
	t.Run("invalid gas price", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()

		tr.app.SetGasStation(gas.NewDummyStation(tr.app)) // Set to dummy station to trigger check gas price check
		signer := crypto.TxSigner{Nonce: uint64(0)}
		tx := &crypto.Tx{Data: nil, From: signer, GasLimit: 1, GasPrice: 1}
		privKey := loadPrivateKey(SEED)
		if err := tx.Sign(privKey); err != nil {
			panic(err)
		}
		wantErr := errors.New("Invalid gas price")
		wantErrCode := CodeTypeInvalidGasPrice

		got, err := tr.app.validateTx(tx, len(tx.Serialize()))
		if got != wantErrCode {
			t.Errorf("App.validateTx() code = %v, want %v", got, wantErrCode)
		}
		if err == nil || wantErr.Error() != err.Error() {
			t.Errorf("App.validateTx() error = %+v, wantErr %v", err, wantErr)
			return
		}
	})

	// invalid data
	t.Run("invalid data", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()
		signer := crypto.TxSigner{Nonce: uint64(0)}
		// Provide un-deserialize-able data to trigger error
		tx := &crypto.Tx{Data: []byte{0}, From: signer, GasLimit: 1000, GasPrice: 1000}
		privKey := loadPrivateKey(SEED)
		if err := tx.Sign(privKey); err != nil {
			panic(err)
		}
		wantErrCode := CodeTypeInvalidData

		got, _ := tr.app.validateTx(tx, len(tx.Serialize()))
		if got != wantErrCode {
			t.Errorf("App.validateTx() code = %v, want %v", got, wantErrCode)
		}
	})

	// non-existent contract account
	t.Run("non-existent contract account", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()
		// Create a random address
		pub, _, _ := ed25519.GenerateKey(cryptoRand.Reader)
		toAddr := crypto.AddressFromPubKey(pub)

		// Create a sign tx
		signer := crypto.TxSigner{Nonce: uint64(0)}
		// Set To to non-exist account address to trigger error
		tx := &crypto.Tx{Data: nil, From: signer, GasLimit: 1, GasPrice: 1, To: toAddr}
		privKey := loadPrivateKey(SEED)
		err := tx.Sign(privKey)
		if err != nil {
			panic(err)
		}
		wantErr := errors.New("contract account not exist")
		wantErrCode := CodeTypeAccountNotExist

		got, err := tr.app.validateTx(tx, len(tx.Serialize()))
		if err == nil || wantErr.Error() != err.Error() {
			t.Errorf("App.validateTx() error = %v, wantErr %v", err, wantErr)
			return
		}
		if got != wantErrCode {
			t.Errorf("App.validateTx() code = %v, want %v", got, wantErrCode)
		}
	})

	// account exist but contains empty contract
	t.Run("empty contract account", func(t *testing.T) {
		tr := NewTestResource()
		defer tr.CleanData()
		// Create a random address
		pub, _, _ := ed25519.GenerateKey(cryptoRand.Reader)
		toAddr := crypto.AddressFromPubKey(pub)
		_, err := tr.app.state.CreateAccount(toAddr, toAddr, []byte{})
		if err != nil {
			panic(err)
		}

		// Create a sign tx
		signer := crypto.TxSigner{Nonce: uint64(0)} // add signer without signature to trigger error
		// Set To to exist account but not contains a contract to trigger error
		tx := &crypto.Tx{Data: nil, From: signer, GasLimit: 1, GasPrice: 1, To: toAddr}
		privKey := loadPrivateKey(SEED)
		err = tx.Sign(privKey)
		if err != nil {
			panic(err)
		}
		wantErr := errors.New("Invoke a non-contract account")
		wantErrCode := CodeTypeNonContractAccount

		got, err := tr.app.validateTx(tx, len(tx.Serialize()))
		if err == nil || wantErr.Error() != err.Error() {
			t.Errorf("App.validateTx() error = %v, wantErr %v", err, wantErr)
			return
		}
		if got != wantErrCode {
			t.Errorf("App.validateTx() code = %v, want %v", got, wantErrCode)
		}
	})
}

func TestApp_DeliverTx(t *testing.T) {
	tr := NewTestResource()
	defer tr.CleanData()
	app := tr.app
	appHash := tr.app.state.Commit()
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: appHash, Time: time.Now()}
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
		pubkey, prvkey, err := ed25519.GenerateKey(cryptoRand.Reader)
		if err != nil {
			panic(err)
		}
		txData := crypto.TxData{}
		gasPrice := uint32(18)
		tx := &crypto.Tx{From: crypto.TxSigner{PubKey: pubkey}, Data: txData.Serialize(), GasPrice: gasPrice}
		err = tx.Sign(prvkey)
		if err != nil {
			panic(err)
		}
		req := types.RequestDeliverTx{Tx: tx.Serialize()}
		detailEvent := event.NewDetailsEvent(1, tx.From.Address(), tx.To, tx.From.Nonce, 0, gasPrice)
		got := app.DeliverTx(req)

		// Detail event must equal
		assert.Equal(t, detailEvent.ToTMEvent(), got.Events[len(got.Events)-1])
		assert.Equal(t, code.CodeTypeOK, got.Code)
		assert.Equal(t, "ok", got.Info)
		assert.Equal(t, int64(tx.GasLimit), got.GasWanted)
		assert.Equal(t, int64(0), got.GasUsed)
	})

	t.Run("Invoke transaction", func(t *testing.T) {
		pubkey, prvkey, err := ed25519.GenerateKey(cryptoRand.Reader)
		if err != nil {
			panic(err)
		}
		txData := crypto.TxData{}
		gasPrice := uint32(18)
		tx := &crypto.Tx{From: crypto.TxSigner{PubKey: pubkey}, Data: txData.Serialize(), GasPrice: gasPrice}
		err = tx.Sign(prvkey)
		if err != nil {
			panic(err)
		}
		req := types.RequestDeliverTx{Tx: tx.Serialize()}
		detailEvent := event.NewDetailsEvent(1, tx.From.Address(), tx.To, tx.From.Nonce, 0, gasPrice)
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
	tr := NewTestResource()
	defer tr.CleanData()
	appHash := tr.app.state.Commit()
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: appHash, Time: time.Now()}
	tr.app.loadState(blockInfo)

	got := tr.app.Commit()
	want := types.ResponseCommit{Data: appHash.Bytes()}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("App.Commit() = %v, want %v", got, want)
	}
}

func TestApp_GetGasContractToken(t *testing.T) {
	// test config for gasContractAddress not exist & gasContractAddress exist, contract not exist
	rand.Seed(time.Now().UTC().UnixNano())
	dbDir := "./testdata/db/test_" + strconv.Itoa(rand.Intn(10000)) + "/"
	err := os.MkdirAll(dbDir, os.ModePerm)
	defer os.RemoveAll(dbDir)
	if err != nil {
		panic(err)
	}

	// init app and loadState without switching gas station
	app := NewApp("testapp", dbDir, "")
	blockInfo := &storage.BlockInfo{Height: 1, AppHash: common.Hash{}, Time: time.Now()}
	app.loadState(blockInfo)

	// Set gasContractAddress to trigger panic
	app.gasContractAddress = "LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7"

	// test config for gasContractAddress exist, contract exist
	tr2 := NewTestResource()
	defer tr2.CleanData()
	address, _ := crypto.AddressFromString("LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7")
	// Get gasContractAccount to create a dummy token
	account, err := tr2.app.state.GetAccount(address)
	if err != nil {
		panic(err)
	}
	token := token.NewToken(tr2.app.state, account)
	tests := []struct {
		name    string
		app     *App
		want    gas.Token
		wantErr error
	}{
		{"gasContractAddress not exist", &App{}, nil, nil},
		{"gasContractAddress exist, contract not exist", app, nil, storage.ErrAccountNotExist},
		{"gasContractAddress exist, contract exist", tr2.app, token, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					err := r.(error)
					if err != tt.wantErr {
						t.Errorf("App.GetGasContractToken() got error = %v, want error %v", err, tt.wantErr)
					}
				}
			}()

			app := tt.app
			if got := app.GetGasContractToken(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.GetGasContractToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
