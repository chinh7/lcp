package consensus

import (
	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
)

const (
	// LastBlockHashKey is key where we store last block hash in infoDB
	LastBlockHashKey string = "last_block_hash"
)

func (app *App) loadLastBlock() error {
	var lastBlockHash common.Hash

	// Load last block hash
	if err := rlp.DecodeBytes(app.InfoDB.Get([]byte(LastBlockHashKey)), &lastBlockHash); err != nil {
		return err
	}

	// Load last block
	if err := rlp.DecodeBytes(app.BlockDB.Get(lastBlockHash[:]), &app.block); err != nil {
		return err
	}

	return nil
}
