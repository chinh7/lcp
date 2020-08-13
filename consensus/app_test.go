package consensus

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/util"
	"github.com/google/uuid"
	"github.com/tendermint/tendermint/abci/types"
)

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
			if got := blockHashToAppHash(tt.blockHash); !reflect.DeepEqual(got, tt.appHash) {
				t.Errorf("blockHashToAppHash() = %v, want %v", got, tt.appHash)
			}

			if got := appHashToBlockHash(tt.appHash); !reflect.DeepEqual(got, tt.blockHash) {
				t.Errorf("appHashToBlockHash() = %v, want %v", got, tt.blockHash)
			}
		})
	}
}

func getDeployTx() *crypto.Transaction {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(0),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	data, err := util.BuildDeployTxPayload("./execution_testdata/contract.wasm", "./execution_testdata/contract-abi.json", "", []string{})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 1,
		Receipt:  &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func getInvokeTx() *crypto.Transaction {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(1),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxData("./execution_testdata/contract-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 0),
		GasLimit: 0,
		GasPrice: 1,
		Receipt:  &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func TestFulLAppFlow(t *testing.T) {
	id, _ := uuid.NewUUID()
	path := fmt.Sprintf("./data-" + id.String())
	app := NewApp(path, "")
	defer func() {
		os.RemoveAll(path)
	}()

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
			tx:                        getDeployTx(),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeOK},
		}},
	}, {
		height: 2,
		time:   time.Unix(0, 2),
		txRequests: []txRequest{{
			tx:                        getInvokeTx(),
			expectedResponseCheckTx:   types.ResponseCheckTx{Code: ResponseCodeOK},
			expectedResponseDeliverTx: types.ResponseDeliverTx{Code: ResponseCodeOK},
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
			if !reflect.DeepEqual(responseCheckTx, txRequest.expectedResponseCheckTx) {
				t.Errorf("app.CheckTx error, got %v, want %v", responseCheckTx, txRequest.expectedResponseCheckTx)
			}

			if responseCheckTx.Code == ResponseCodeOK {
				responseDeliverTx := app.DeliverTx(types.RequestDeliverTx{Tx: rawTx})
				if !reflect.DeepEqual(responseDeliverTx, txRequest.expectedResponseDeliverTx) {
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
