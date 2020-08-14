package consensus

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/constant"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/abci/types"
)

func newAppTestResource() *TestResource {
	rand.Seed(time.Now().UTC().UnixNano())
	dbDir := "./testdata/db/test_" + strconv.Itoa(rand.Intn(10000)) + "/"
	err := os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	app := NewApp(dbDir, "")
	app.state.LoadState(crypto.GenesisBlock.Header)

	return &TestResource{app, dbDir}
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

	gasContractAddress := "LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ7"

	app := NewApp(dbDir, gasContractAddress)
	assert.NotNil(t, app)
}

func TestApp_BeginBlock(t *testing.T) {
	tr := newAppTestResource()
	defer tr.cleanData()
	app := tr.app

	t.Run("Should load state from genesis block", func(t *testing.T) {
		reqHeight := int64(0)
		previousBlockHash := common.EmptyHash.Bytes()
		stateRootHash, txRootHash := tr.app.state.Commit()

		req := types.RequestBeginBlock{Header: types.Header{Height: reqHeight, AppHash: previousBlockHash}}
		got := app.BeginBlock(req)
		want := types.ResponseBeginBlock{}
		if !cmp.Equal(got, want) {
			t.Errorf("App.BeginBlock() = %v, want %v", got, want)
		}

		// loadState() should be called
		assert.NotNil(t, app.state)
		assert.Equal(t, uint64(reqHeight), app.state.GetBlockHeader().Height)
		assert.Equal(t, stateRootHash, app.state.Hash())
		assert.Equal(t, txRootHash, app.state.Hash())
	})

	t.Run("Should load state", func(t *testing.T) {
		stateRootHash, txRootHash := tr.app.state.Commit()
		reqHeight := int64(1)

		previousBlock := crypto.Block{
			Header:       &crypto.BlockHeader{Height: uint64(reqHeight), Time: time.Now(), Parent: common.EmptyHash, StateRoot: stateRootHash, TransactionRoot: txRootHash},
			Transactions: nil,
		}

		rawBlock, _ := previousBlock.Encode()
		blockHash := previousBlock.Header.Hash()
		app.block.Put(blockHash[:], rawBlock)

		assert.NotNil(t, app.state)
		assert.Equal(t, uint64(0), app.state.GetBlockHeader().Height)

		req := types.RequestBeginBlock{Header: types.Header{Height: reqHeight, AppHash: blockHash.Bytes()}}
		got := app.BeginBlock(req)
		want := types.ResponseBeginBlock{}
		if !cmp.Equal(got, want) {
			t.Errorf("App.BeginBlock() = %v, want %v", got, want)
		}

		// loadState() should be called
		assert.NotNil(t, app.state)
		assert.Equal(t, uint64(reqHeight), app.state.GetBlockHeader().Height)
		assert.Equal(t, stateRootHash, app.state.Hash())
		assert.Equal(t, txRootHash, app.state.Hash())
	})
}

func TestApp_Info(t *testing.T) {
	tr := newAppTestResource()
	defer tr.cleanData()
	app := tr.app

	t.Run("Should return valid response", func(t *testing.T) {
		height := 2
		stateRootHash, txRootHash := tr.app.state.Commit()
		block := crypto.Block{
			Header:       &crypto.BlockHeader{Height: uint64(height), Time: time.Now(), Parent: common.EmptyHash, StateRoot: stateRootHash, TransactionRoot: txRootHash},
			Transactions: nil,
		}
		app.meta.StoreBlockIndexes(&block)

		got := app.Info(types.RequestInfo{})
		// returns correct current state
		want := types.ResponseInfo{
			LastBlockHeight:  int64(height),
			LastBlockAppHash: block.Header.Hash().Bytes(),
		}

		if !cmp.Equal(got, want) {
			t.Errorf("Got app.Info() = %v, want %v", got, want)
		}
	})
}

func TestApp_CheckTx(t *testing.T) {
	tr := newAppTestResource()
	defer tr.cleanData()
	app := tr.app

	app.BeginBlock(types.RequestBeginBlock{
		Header: types.Header{
			Height:  1,
			Time:    time.Now(),
			AppHash: []byte{},
		},
	})
	deployTx, _ := tr.getDeployTx(0).Serialize()
	app.DeliverTx(types.RequestDeliverTx{Tx: deployTx})
	app.Commit()

	t.Run("CheckTx with error transactions", func(t *testing.T) {
		type txRequest struct {
			tx                      *crypto.Transaction
			expectedResponseCheckTx types.ResponseCheckTx
		}

		checkTxTestTable := []txRequest{{
			tx:                      tr.getInvokeNilContractTx(1),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invoke nil contract"},
		}, {
			tx:                      tr.getInvalidMaxSizeTx(1),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: fmt.Sprintf("Transaction size exceed %vB", constant.MaxTransactionSize)},
		}, {
			tx:                      tr.getInvalidSignatureTx(1),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invalid signature"},
		}, {
			tx:                      tr.getInvalidNonceTx(2),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invalid nonce. Expected 1, got 2"},
		}, {
			tx:                      tr.getInvalidGasPriceTx(1),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invalid gas price"},
		}, {
			tx:                      tr.getInvokeNonContractTx(1),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invoke a non-contract account"},
		}}

		for i, checkTxTest := range checkTxTestTable {
			rawTx, _ := checkTxTest.tx.Serialize()
			got := app.CheckTx(types.RequestCheckTx{Tx: rawTx})
			want := checkTxTest.expectedResponseCheckTx
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("Case %d: App.CheckTx() is expected to be = %v, got %v", i+1, want, got)
			}
		}
	})
}

func TestApp_DeliverTx(t *testing.T) {
	tr := newAppTestResource()
	defer tr.cleanData()
	app := tr.app

	app.BeginBlock(types.RequestBeginBlock{
		Header: types.Header{
			Height:  1,
			Time:    time.Now(),
			AppHash: []byte{},
		},
	})
	deployTx, _ := tr.getDeployTx(0).Serialize()
	app.DeliverTx(types.RequestDeliverTx{Tx: deployTx})
	app.Commit()

	t.Run("Deserialize tx error", func(t *testing.T) {
		got := app.DeliverTx(types.RequestDeliverTx{Tx: []byte{1, 2, 3}})
		want := types.ResponseDeliverTx{Code: ResponseCodeNotOK}

		if !cmp.Equal(got, want) {
			t.Errorf("App.DeliverTx() = %v, want %v", got, want)
		}
	})

	t.Run("DeliverTx with error transactions", func(t *testing.T) {
		type txRequest struct {
			tx                        *crypto.Transaction
			expectedResponseDeliverTx types.ResponseDeliverTx
		}

		deliverTxTestTable := []txRequest{{
			tx:                        tr.getInvokeNilContractTx(1),
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeNotOK},
		}, {
			tx:                        tr.getInvalidMaxSizeTx(1),
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeNotOK},
		}, {
			tx:                        tr.getInvalidSignatureTx(1),
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeNotOK},
		}, {
			tx:                        tr.getInvalidNonceTx(2),
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeNotOK},
		}, {
			tx:                        tr.getInvalidGasPriceTx(1),
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeNotOK},
		}, {
			tx:                        tr.getInvokeNonContractTx(1),
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeNotOK},
		}}

		for i, deliverTxTest := range deliverTxTestTable {
			rawTx, _ := deliverTxTest.tx.Serialize()
			got := app.DeliverTx(types.RequestDeliverTx{Tx: rawTx})
			want := deliverTxTest.expectedResponseDeliverTx
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("Case %d: App.DeliverTx() is expected to be = %v, got %v", i+1, want, got)
			}
		}
	})

	t.Run("DeliverTx with success transactions", func(t *testing.T) {
		type txRequest struct {
			tx                        *crypto.Transaction
			expectedResponseDeliverTx types.ResponseDeliverTx
		}

		deliverTxTestTable := []txRequest{{
			tx:                        tr.getDeployTx(1),
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeOK},
		}, {
			tx:                        tr.getInvokeTx(2),
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeOK},
		}}

		for i, deliverTxTest := range deliverTxTestTable {
			rawTx, _ := deliverTxTest.tx.Serialize()
			got := app.DeliverTx(types.RequestDeliverTx{Tx: rawTx})
			want := deliverTxTest.expectedResponseDeliverTx
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("Case %d: App.DeliverTx() is expected to be = %v, got %v", i+1, want, got)
			}
		}
	})
}

func TestBlockHashAndAppHashConversion(t *testing.T) {
	tests := []struct {
		name      string
		appHash   []byte
		blockHash common.Hash
	}{{
		name:      "Empty",
		appHash:   []byte{},
		blockHash: common.EmptyHash,
	}, {
		name:      "Normal",
		appHash:   []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 34},
		blockHash: common.Hash{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 34},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := blockHashToAppHash(tt.blockHash); !cmp.Equal(got, tt.appHash) {
				t.Errorf("blockHashToAppHash() = %v, want %v", got, tt.appHash)
			}

			if got := appHashToBlockHash(tt.appHash); !cmp.Equal(got, tt.blockHash) {
				t.Errorf("appHashToBlockHash() = %v, want %v", got, tt.blockHash)
			}
		})
	}
}

func TestFullAppFlow(t *testing.T) {
	tr := newAppTestResource()
	defer tr.cleanData()
	app := tr.app

	type txRequest struct {
		tx                        *crypto.Transaction
		expectedResponseCheckTx   types.ResponseCheckTx
		expectedResponseDeliverTx types.ResponseDeliverTx
	}

	type round struct {
		height int64
		time   time.Time

		txRequests []txRequest
	}

	rounds := []round{{
		height: 1,
		time:   time.Unix(0, 1),
		txRequests: []txRequest{{
			tx:                        tr.getDeployTx(0),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeOK},
		}},
	}, {
		height: 2,
		time:   time.Unix(0, 2),
		txRequests: []txRequest{{
			tx:                        tr.getInvokeTx(1),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeOK},
		}},
	}, {
		height: 3,
		time:   time.Unix(0, 3),
		txRequests: []txRequest{{
			tx:                      tr.getInvokeNilContractTx(2),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invoke nil contract"},
		}, {
			tx:                      tr.getInvalidMaxSizeTx(2),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: fmt.Sprintf("Transaction size exceed %vB", constant.MaxTransactionSize)},
		}, {
			tx:                      tr.getInvalidSignatureTx(2),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invalid signature"},
		}, {
			tx:                      tr.getInvalidNonceTx(123),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invalid nonce. Expected 2, got 123"},
		}, {
			tx:                      tr.getInvalidGasPriceTx(2),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invalid gas price"},
		}, {
			tx:                      tr.getInvokeNonContractTx(2),
			expectedResponseCheckTx: types.ResponseCheckTx{Code: ResponseCodeNotOK, Log: "Invoke a non-contract account"},
		}},
	}}

	appHash := []byte{}
	for _, round := range rounds {
		app.BeginBlock(types.RequestBeginBlock{
			Header: types.Header{
				Height:  round.height,
				Time:    round.time,
				AppHash: appHash,
			},
		})

		for _, txRequest := range round.txRequests {
			rawTx, _ := txRequest.tx.Serialize()
			responseCheckTx := app.CheckTx(types.RequestCheckTx{Tx: rawTx})
			if !cmp.Equal(responseCheckTx, txRequest.expectedResponseCheckTx) {
				t.Errorf("app.CheckTx error, got %v, want %v", responseCheckTx, txRequest.expectedResponseCheckTx)
			}

			if responseCheckTx.Code == ResponseCodeOK {
				responseDeliverTx := app.DeliverTx(types.RequestDeliverTx{Tx: rawTx})
				if !cmp.Equal(responseDeliverTx, txRequest.expectedResponseDeliverTx) {
					t.Errorf("app.CheckTx error, got %v, want %v", responseDeliverTx, txRequest.expectedResponseDeliverTx)
				}
			}
		}

		responseCommit := app.Commit()
		appHash = responseCommit.Data
		info := app.Info(types.RequestInfo{})
		if !bytes.Equal(info.LastBlockAppHash, appHash) {
			t.Errorf("Commit app hash = %v, is different from info app hash = %v", appHash, info.LastBlockAppHash)
		}
	}
}
