package consensus

import (
	"fmt"
	"strconv"

	"github.com/QuoineFinancial/vertex/core"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"

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
}

// NewApp initializes a new app
func NewApp(nodeInfo string) *App {
	return &App{
		nodeInfo: nodeInfo,
	}
}

// BeginBlock begins new block
func (app *App) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	appHash := gethCommon.BytesToHash(req.Header.GetAppHash())
	app.state = storage.GetState(appHash)
	return types.ResponseBeginBlock{}
}

// Info returns application chain info
func (app *App) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: fmt.Sprintf("{\"version\":%s}", app.nodeInfo)}
}

// CheckTx checks if submitted transaction is valid and can be passed to next step
func (app *App) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

//DeliverTx executes the submitted transaction
func (app *App) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	code := CodeTypeOK
	info := "ok"
	tx := &crypto.Tx{}
	tx.Deserialize(req.GetTx())
	applyEvents, err := core.ApplyTx(app.state, tx)
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
		Code:   code,
		Events: events,
		Info:   info,
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
