package consensus

import (
	"fmt"
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
	state              *storage.State
	nodeInfo           string
	Database           db.Database
	gasStation         gas.Station
	gasContractAddress string
}

// NewApp initializes a new app
func NewApp(nodeInfo string, dbPath string, gasContractAddress string) *App {
	database := db.NewRocksDB(dbPath)
	app := &App{
		nodeInfo:           nodeInfo,
		Database:           database,
		gasContractAddress: gasContractAddress,
	}
	app.SetGasStation(gas.NewFreeGasStation(app))
	return app
}

// BeginBlock begins new block
func (app *App) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	var err error
	appHash := gethCommon.BytesToHash(req.Header.GetAppHash())
	if app.state, err = storage.New(appHash, app.Database); err != nil {
		panic(err)
	}

	// Keep moving forward
	for app.gasStation.Switch() {
	}
	return types.ResponseBeginBlock{}
}

// Info returns application chain info
func (app *App) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: fmt.Sprintf("{\"version\":%s}", app.nodeInfo)}
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
		GasUsed:   gasUsed,
	}
}

// Commit returns the state root of application storage. Called once all block processing is complete
func (app *App) Commit() types.ResponseCommit {
	appHash, _ := app.state.Commit()
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
			return token.NewToken(contract)
		}
	}
	return nil
}
