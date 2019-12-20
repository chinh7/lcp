package consensus

import (
	"encoding/binary"
	"fmt"
	"path/filepath"

	"github.com/QuoineFinancial/liquid-chain/constant"
	"github.com/QuoineFinancial/liquid-chain/core"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/event"
	"github.com/QuoineFinancial/liquid-chain/gas"
	"github.com/QuoineFinancial/liquid-chain/storage"
	"github.com/QuoineFinancial/liquid-chain/token"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
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

	app.SetGasStation(gas.NewFreeStation(app))

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
	tx := &crypto.Tx{}
	if err := tx.Deserialize(req.GetTx()); err != nil {
		return types.ResponseCheckTx{
			Code: CodeTypeEncodingError,
			Log:  err.Error(),
		}
	}

	if code, err := app.validateTx(tx, len(req.GetTx())); err != nil {
		return types.ResponseCheckTx{
			Code: code,
			Log:  err.Error(),
		}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *App) validateTx(tx *crypto.Tx, txSize int) (uint32, error) {
	// Validate tx size
	if txSize > constant.MaxTransactionSize {
		err := fmt.Errorf("Transaction size exceed %dB", constant.MaxTransactionSize)
		return code.CodeTypeUnknownError, err
	}

	nonce := uint64(0)
	address := tx.From.Address()
	account, _ := app.state.GetAccount(address)
	if account != nil {
		nonce = account.Nonce
	}

	// Validate tx nonce
	if tx.From.Nonce != nonce {
		err := fmt.Errorf("Invalid nonce. Expected %v, got %v", nonce, tx.From.Nonce)
		return code.CodeTypeBadNonce, err
	}

	// Validate tx signature
	if !tx.SigVerified() {
		return code.CodeTypeUnknownError, fmt.Errorf("Invalid signature")
	}

	// Validate gas limit
	fee, err := tx.GetFee()
	if err != nil {
		return code.CodeTypeUnknownError, err
	}
	if !app.gasStation.Sufficient(address, fee) {
		return code.CodeTypeBadNonce, fmt.Errorf("Insufficient fee")
	}

	// Validate tx data
	txData := &crypto.TxData{}
	err = txData.Deserialize(tx.Data)
	if err != nil {
		return code.CodeTypeUnknownError, err
	}

	return code.CodeTypeOK, nil
}

//DeliverTx executes the submitted transaction
func (app *App) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	tx := &crypto.Tx{}
	if err := tx.Deserialize(req.GetTx()); err != nil {
		return types.ResponseDeliverTx{
			Code: CodeTypeEncodingError,
			Log:  err.Error(),
		}
	}
	if code, err := app.validateTx(tx, len(req.GetTx())); err != nil {
		return types.ResponseDeliverTx{
			Code: code,
			Log:  err.Error(),
		}
	}

	info := "ok"
	codeType := CodeTypeOK
	result, applyEvents, gasUsed, err := core.ApplyTx(app.state, tx, app.gasStation)
	if err != nil {
		codeType = CodeTypeUnknownError
		info = err.Error()
	}
	fromAddress := tx.From.Address()
	detailEvent := event.NewDetailsEvent(fromAddress, tx.To, tx.From.Nonce, result)
	events := append(applyEvents, detailEvent)
	tmEvents := make([]types.Event, len(events))
	for index := range events {
		tmEvents[index] = events[index].ToTMEvent()
	}
	return types.ResponseDeliverTx{
		Code:      codeType,
		Events:    tmEvents,
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
