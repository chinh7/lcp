package consensus

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/constant"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/QuoineFinancial/liquid-chain/token"

	abciTypes "github.com/tendermint/tendermint/abci/types"
)

// App basic Tendermint base app
type App struct {
	abciTypes.BaseApplication

	meta  *storage.MetaStorage
	state *storage.StateStorage
	block *storage.BlockStorage

	gasStation         gas.Station
	gasContractAddress string
}

func blockHashToAppHash(blockHash common.Hash) []byte {
	if blockHash == common.EmptyHash {
		return []byte{}
	}
	return blockHash.Bytes()
}

func appHashToBlockHash(appHash []byte) common.Hash {
	if len(appHash) == 0 {
		return common.EmptyHash
	}
	return common.BytesToHash(appHash)
}

// NewApp initializes a new app
func NewApp(dbDir string, gasContractAddress string) *App {
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.Mkdir(dbDir, os.ModePerm)
	}
	app := &App{
		meta:               storage.NewMetaStorage(db.NewRocksDB(filepath.Join(dbDir, "index.db"))),
		state:              storage.NewStateStorage(db.NewRocksDB(filepath.Join(dbDir, "state.db"))),
		block:              storage.NewBlockStorage(db.NewRocksDB(filepath.Join(dbDir, "block.db"))),
		gasContractAddress: gasContractAddress,
	}
	app.SetGasStation(gas.NewFreeStation(app))
	return app
}

// BeginBlock begins new block
func (app *App) BeginBlock(req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	lastBlockHash := appHashToBlockHash(req.Header.AppHash)
	previousBlock := app.block.MustGetBlock(lastBlockHash)
	app.state.MustLoadState(previousBlock.Header)
	app.block.ComposeBlock(previousBlock, req.Header.Time)
	for app.gasStation.Switch() {
	}
	return abciTypes.ResponseBeginBlock{}
}

// Info returns application chain info
func (app *App) Info(req abciTypes.RequestInfo) (resInfo abciTypes.ResponseInfo) {
	lastBlockHeight := app.meta.LatestBlockHeight()
	lastBlockHash := app.meta.HeightToBlockHash(lastBlockHeight)
	return abciTypes.ResponseInfo{
		LastBlockHeight:  int64(lastBlockHeight),
		LastBlockAppHash: blockHashToAppHash(lastBlockHash),
	}
}

// CheckTx checks if submitted transaction is valid and can be passed to next step
func (app *App) CheckTx(req abciTypes.RequestCheckTx) abciTypes.ResponseCheckTx {
	if len(req.Tx) > constant.MaxTransactionSize {
		return abciTypes.ResponseCheckTx{
			Code: CodeTypeExceedTransactionSize,
			Log:  fmt.Sprintf("Transaction size exceed %dB", constant.MaxTransactionSize),
		}
	}

	var tx crypto.Transaction
	if err := tx.Deserialize(req.GetTx()); err != nil {
		return abciTypes.ResponseCheckTx{
			Code: CodeTypeEncodingError,
			Log:  err.Error(),
		}
	}

	if code, err := app.validateTx(&tx); err != nil {
		return abciTypes.ResponseCheckTx{
			Code: code,
			Log:  err.Error(),
		}
	}

	return abciTypes.ResponseCheckTx{Code: CodeTypeOK}
}

//DeliverTx executes the submitted transaction
func (app *App) DeliverTx(req abciTypes.RequestDeliverTx) abciTypes.ResponseDeliverTx {
	var tx crypto.Transaction
	if err := tx.Deserialize(req.GetTx()); err != nil {
		return abciTypes.ResponseDeliverTx{}
	}

	if _, err := app.validateTx(&tx); err != nil {
		return abciTypes.ResponseDeliverTx{}
	}

	if receipt, err := app.applyTransaction(&tx); err != nil {
		panic(err)
	} else {
		tx.Receipt = receipt
	}

	if err := app.state.AddTransaction(&tx); err != nil {
		log.Fatal(err)
	}
	app.block.AddTransaction(&tx)

	return abciTypes.ResponseDeliverTx{}
}

// Commit returns the state root of application storage. Called once all block processing is complete
func (app *App) Commit() abciTypes.ResponseCommit {
	stateRootHash, txRootHash := app.state.Commit()
	app.block.FinalizeBlock(stateRootHash, txRootHash)
	blockHash, err := app.block.Commit()
	if err != nil {
		panic(err)
	}
	app.meta.StoreBlockIndexes(app.block.MustGetBlock(blockHash))
	return abciTypes.ResponseCommit{Data: blockHashToAppHash(blockHash)}
}

// SetGasStation active the gas station
func (app *App) SetGasStation(gasStation gas.Station) {
	app.gasStation = gasStation
}

// GetGasContractToken designated
func (app *App) GetGasContractToken() gas.Token {
	if len(app.gasContractAddress) > 0 {
		address, err := crypto.AddressFromString(app.gasContractAddress)
		if err != nil {
			panic(err)
		}
		contract, err := app.state.GetAccount(address)
		if err != nil {
			panic(err)
		}
		return token.NewToken(app.state, contract)
	}
	return nil
}
