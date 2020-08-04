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

	state *storage.State
	block *crypto.Block

	InfoDB  db.Database
	StateDB db.Database
	BlockDB db.Database

	gasStation         gas.Station
	gasContractAddress string
}

// NewApp initializes a new app
func NewApp(dbDir string, gasContractAddress string) *App {
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.Mkdir(dbDir, os.ModePerm)
	}
	app := &App{
		BlockDB:            db.NewRocksDB(filepath.Join(dbDir, "block.db")),
		StateDB:            db.NewRocksDB(filepath.Join(dbDir, "state.db")),
		InfoDB:             db.NewRocksDB(filepath.Join(dbDir, "info.db")),
		gasContractAddress: gasContractAddress,
	}
	app.SetGasStation(gas.NewFreeStation(app))
	if err := app.loadLastBlock(); err == nil {
		app.LoadState(app.block.Header)
	}
	return app
}

// LoadState fetch app state from block header
func (app *App) LoadState(blockHeader *crypto.BlockHeader) {
	var err error
	if app.state, err = storage.NewState(blockHeader, app.StateDB); err != nil {
		panic(err)
	}

	// Keep switching until a desire gasStation is meet
	for app.gasStation.Switch() {
	}
}

// BeginBlock begins new block
func (app *App) BeginBlock(req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	var previousBlockHash common.Hash
	copy(previousBlockHash[:], req.Header.AppHash)

	if previousBlockHash == common.EmptyHash {
		app.LoadState(&crypto.GenesisBlock)
	} else {
		rawBlock := app.BlockDB.Get(previousBlockHash[:])
		previousBlock := crypto.MustDecodeBlock(rawBlock)
		app.LoadState(previousBlock.Header)
	}
	app.block = crypto.NewEmptyBlock(previousBlockHash, uint64(req.Header.GetHeight()), req.Header.GetTime())
	return abciTypes.ResponseBeginBlock{}
}

// Info returns application chain info
func (app *App) Info(req abciTypes.RequestInfo) (resInfo abciTypes.ResponseInfo) {
	var lastBlockHeight int64
	var lastBlockAppHash []byte

	if app.state != nil {
		lastBlockHeight = int64(app.state.GetBlockHeader().Height)
		lastBlockAppHash = app.state.GetBlockHeader().Hash().Bytes()
	}

	return abciTypes.ResponseInfo{
		LastBlockHeight:  lastBlockHeight,
		LastBlockAppHash: lastBlockAppHash,
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
		return abciTypes.ResponseDeliverTx{
			Code: CodeTypeEncodingError,
			Log:  err.Error(),
		}
	}

	if code, err := app.validateTx(&tx); err != nil {
		return abciTypes.ResponseDeliverTx{
			Code: code,
			Log:  err.Error(),
		}
	}

	if receipt, err := app.applyTransaction(&tx); err != nil {
		panic(err)
	} else {
		tx.Receipt = receipt
	}

	if err := app.state.AddTransaction(&tx); err != nil {
		log.Fatal(err)
	}
	app.block.Transactions = append(app.block.Transactions, &tx)

	return abciTypes.ResponseDeliverTx{}
}

// Commit returns the state root of application storage. Called once all block processing is complete
func (app *App) Commit() abciTypes.ResponseCommit {
	app.block.Header.StateRoot = app.state.Commit()
	app.block.Header.TransactionRoot = app.state.Commit()
	rawBlock, err := app.block.Encode()
	if err != nil {
		log.Fatal(err)
	}
	blockHash := app.block.Header.Hash()
	app.BlockDB.Put(blockHash[:], rawBlock)
	app.InfoDB.Put([]byte(LastBlockHashKey), blockHash[:])
	return abciTypes.ResponseCommit{Data: blockHash[:]}
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
			switch err {
			case storage.ErrAccountNotExist:
				return nil
			default:
				panic(err)
			}
		}
		return token.NewToken(app.state, contract)
	}
	return nil
}
