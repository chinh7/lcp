package api

import (
	"net/http"
	"time"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

// BlockArgs is params of BlockService
type BlockArgs struct {
	Height *int64
}

// BlockReply is response of BlockService
type BlockReply struct {
	Hash            string     `json:"hash"`
	Time            time.Time  `json:"time"`
	Height          int64      `json:"height"`
	TotalTxs        int64      `json:"total_txs"`
	ChainID         string     `json:"chain_id"`
	LastCommitHash  string     `json:"last_commit_hash"`  // commit from validators from the last block
	DataHash        string     `json:"data_hash"`         // transactions
	ConsensusHash   string     `json:"consensus_hash"`    // consensus params for current block
	AppHash         string     `json:"app_hash"`          // state after txs from the previous block
	LastResultsHash string     `json:"last_results_hash"` // root hash of all results from the txs from the previous block
	Txs             []types.Tx `json:"txs"`
}

// BlockService is first service
type BlockService struct {
	client *rpcclient.Client
}

// NewBlockService returns new instance of BlockService
func (api *API) NewBlockService() *BlockService {
	if api.Client == nil {
		panic("api.NewBlockService call without api.Client")
	}
	return &BlockService{api.Client}
}

// Get is handler of BlockService
func (service *BlockService) Get(r *http.Request, args *BlockArgs, reply *BlockReply) error {
	client := *service.client
	block, err := client.Block(args.Height)
	if err != nil {
		return err
	}
	reply.Hash = block.BlockMeta.Header.Hash().String()
	reply.Height = block.BlockMeta.Header.Height
	reply.ChainID = block.BlockMeta.Header.ChainID
	reply.Time = block.BlockMeta.Header.Time
	reply.TotalTxs = block.BlockMeta.Header.TotalTxs
	reply.ConsensusHash = block.BlockMeta.Header.ConsensusHash.String()
	reply.AppHash = block.BlockMeta.Header.AppHash.String()
	reply.LastResultsHash = block.BlockMeta.Header.LastResultsHash.String()
	reply.Txs = block.Block.Data.Txs
	return nil
}
