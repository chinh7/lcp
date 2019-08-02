package consensus

import (
	"fmt"
	"log"

	"github.com/QuoineFinancial/vertex/core"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"

	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
)

// App basic Tendermint base app
type App struct {
	types.BaseApplication
	nodeInfo string
}

// NewApp initializes a new app
func NewApp(nodeInfo string) *App {
	crypto.RegisterCodec()
	return &App{
		nodeInfo: nodeInfo,
	}
}

// Info returns application chain info
func (app *App) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: fmt.Sprintf("{\"version\":%s}", app.nodeInfo)}
}

// CheckTx checks if submitted transaction is valid and can be passed to next step
func (app *App) CheckTx(tx types.RequestCheckTx) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

//DeliverTx executes the submitted transaction
func (app *App) DeliverTx(deliverTx types.RequestDeliverTx) types.ResponseDeliverTx {
	// timed out after config.TimeoutBroadcastTxCommit
	// time.Sleep(5 * time.Second)
	// log.Println("DeliverTx", hex.EncodeToString(txBytes))
	tx := &crypto.Tx{}
	tx.Deserialize(deliverTx.GetTx())
	log.Println(tx)
	core.ApplyTx(tx)
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: core.GetEvents()}
}

// Commit returns the state root of application storage. Called once all block processing is complete
func (app *App) Commit() types.ResponseCommit {
	// Using a memdb - just return the big endian size of the db
	appHash, _ := storage.GetState().Commit()
	return types.ResponseCommit{Data: appHash[:]}
}
