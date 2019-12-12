package consensus

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/QuoineFinancial/vertex/core"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/db"
	"github.com/QuoineFinancial/vertex/gas"
	"github.com/QuoineFinancial/vertex/storage"
	"github.com/QuoineFinancial/vertex/token"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

// App basic Tendermint base app
type App struct {
	types.BaseApplication
	state    *storage.State
	nodeInfo string

	InfoDB  db.Database
	StateDB db.Database

	gasStation         gas.Station
	gasContractAddress string

	lastBlockHeight  int64
	lastBlockAppHash []byte
}

// NewApp initializes a new app
func NewApp(nodeInfo string, dbDir string, gasContractAddress string) *App {
	infoDB := db.NewRocksDB(filepath.Join(dbDir, "info.db"))
	stateDB := db.NewRocksDB(filepath.Join(dbDir, "storage.db"))

	app := &App{
		nodeInfo:           nodeInfo,
		StateDB:            stateDB,
		InfoDB:             infoDB,
		gasContractAddress: gasContractAddress,
		lastBlockHeight:    0,
	}

	app.SetGasStation(gas.NewFreeGasStation(app))

	// Load last proccessed block height
	b := app.InfoDB.Get([]byte("lastBlockHeight"))
	lastBlockAppHash := app.InfoDB.Get([]byte("lastBlockAppHash"))

	if len(b) > 0 && len(lastBlockAppHash) > 0 {
		app.loadState(int64(binary.LittleEndian.Uint64(b)), lastBlockAppHash)
	}

	return app
}

func (app *App) loadState(height int64, hash []byte) {
	var err error
	if app.state, err = storage.New(gethCommon.BytesToHash(hash), app.StateDB); err != nil {
		panic(err)
	}
	app.lastBlockHeight = height
	app.lastBlockAppHash = hash
	// Keep moving forward
	for app.gasStation.Switch() {
	}
}

// BeginBlock begins new block
func (app *App) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	app.loadState(req.Header.Height, req.Header.AppHash)
	return types.ResponseBeginBlock{}
}

// Info returns application chain info
func (app *App) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{
		Data:             fmt.Sprintf("{\"version\":%s}", app.nodeInfo),
		LastBlockHeight:  app.lastBlockHeight,
		LastBlockAppHash: app.lastBlockAppHash,
	}
}

// CheckTx checks if submitted transaction is valid and can be passed to next step
func (app *App) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	// Check sig
	// Check nonce
	// Check gas wanted (limit)
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

//DeliverTx executes the submitted transaction
func (app *App) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	code := CodeTypeOK
	info := "ok"
	tx := &crypto.Tx{}
	tx.Deserialize(req.GetTx())
	applyEvents, gasUsed, err := core.ApplyTx(app.state, tx, app.gasStation)
	if err != nil {
		code = CodeTypeUnknownError
		info = err.Error()
	}
	fromAddress := tx.From.Address()
	events := append(applyEvents, types.Event{
		Type: "detail",
		Attributes: []common.KVPair{
			common.KVPair{
				Key: []byte("from"), Value: []byte(fromAddress.String()),
			},
			common.KVPair{
				Key: []byte("to"), Value: []byte(tx.To.String()),
			},
			common.KVPair{
				Key: []byte("nonce"), Value: []byte(strconv.FormatUint(tx.From.Nonce, 10)),
			},
		},
	})
	return types.ResponseDeliverTx{
		Code:      code,
		Events:    events,
		Info:      info,
		GasWanted: int64(tx.GasLimit),
		GasUsed:   int64(gasUsed),
	}
}

// Commit returns the state root of application storage. Called once all block processing is complete
func (app *App) Commit() types.ResponseCommit {
	appHash := app.state.Commit()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(app.lastBlockHeight))
	app.InfoDB.Put([]byte("lastBlockHeight"), b)
	app.InfoDB.Put([]byte("lastBlockAppHash"), app.lastBlockAppHash)
	return types.ResponseCommit{Data: appHash[:]}
}

func (app *App) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	return types.ResponseQuery{Log: "hello"}
}

// SetGasStation active the gas station
func (app *App) SetGasStation(gasStation gas.Station) {
	app.gasStation = gasStation
}

// GetGasContractToken designated
func (app *App) GetGasContractToken() gas.Token {
	if len(app.gasContractAddress) > 0 {
		address := crypto.AddressFromString(app.gasContractAddress)
		contract, err := app.state.GetAccount(address)
		if err != nil {
			panic(err)
		}
		if contract != nil {
			return token.NewToken(app.state, contract)
		}
	}
	return nil
}
