package consensus

import (
	"fmt"
	"log"
	"strconv"

	"github.com/QuoineFinancial/vertex/core"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/storage"

	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

// App basic Tendermint base app
type App struct {
	types.BaseApplication
	nodeInfo string
}

// NewApp initializes a new app
func NewApp(nodeInfo string) *App {
	return &App{
		nodeInfo: nodeInfo,
	}
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
	// timed out after config.TimeoutBroadcastTxCommit
	// time.Sleep(5 * time.Second)
	// log.Println("DeliverTx", hex.EncodeToString(txBytes))
	tx := &crypto.Tx{}
	tx.Deserialize(req.GetTx())
	log.Println(tx)
	core.ApplyTx(tx)
	events := core.GetEvents()
	fromAddress := tx.From.Address()
	attributes := []common.KVPair{
		common.KVPair{
			Key: []byte("tx.from"), Value: []byte(fromAddress.String()),
		},
		common.KVPair{
			Key: []byte("tx.to"), Value: []byte(tx.To.String()),
		},
		common.KVPair{
			Key: []byte("tx.nonce"), Value: []byte(strconv.FormatUint(tx.From.Nonce, 10)),
		},
	}
	events = append(events, types.Event{
		Attributes: attributes,
	})
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: events}
}

// Commit returns the state root of application storage. Called once all block processing is complete
func (app *App) Commit() types.ResponseCommit {
	// Using a memdb - just return the big endian size of the db
	appHash, _ := storage.GetState().Commit()
	return types.ResponseCommit{Data: appHash[:]}
}

func (app *App) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	return types.ResponseQuery{Log: "hello"}
}
