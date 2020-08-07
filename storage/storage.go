package storage

import (
	"log"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/trie"
)

// State is the global account state consisting of many address->state mapping
type State struct {
	db                db.Database
	blockHeader       *crypto.BlockHeader
	txTrie            *trie.Trie
	stateTrie         *trie.Trie
	accounts          map[crypto.Address]*Account
	accountCheckpoint common.Hash
}

// NewState returns a state database
func NewState(blockHeader *crypto.BlockHeader, db db.Database) (*State, error) {
	stateTrie, err := trie.New(blockHeader.StateRoot, db)
	if err != nil {
		return nil, err
	}

	txTrie, err := trie.New(blockHeader.TransactionRoot, db)
	if err != nil {
		return nil, err
	}

	return &State{
		db:                db,
		blockHeader:       blockHeader,
		txTrie:            txTrie,
		accounts:          make(map[crypto.Address]*Account),
		stateTrie:         stateTrie,
		accountCheckpoint: blockHeader.StateRoot,
	}, nil
}

// GetBlockHeader return header of block that inits current state
func (state *State) GetBlockHeader() *crypto.BlockHeader {
	return state.blockHeader
}

// Hash retrive hash of entire state
func (state *State) Hash() common.Hash {
	var err error
	for _, account := range state.accounts {
		if account == nil || !account.dirty {
			continue
		}

		// Update account storage
		account.StorageHash = account.storage.Hash()

		// Update account
		raw, _ := rlp.EncodeToBytes(account)
		if err = state.stateTrie.Update(account.address[:], raw); err != nil {
			panic(err)
		}
	}
	return state.stateTrie.Hash()
}

// Commit stores all dirty Accounts to storage.trie
func (state *State) Commit() (common.Hash, common.Hash) {
	var err error
	for _, account := range state.accounts {
		if account == nil || !account.dirty {
			continue
		}

		if account.IsContract() {
			// Update contract
			state.db.Put(account.ContractHash[:], account.contract)
		}

		// Update account storage
		if account.StorageHash, err = account.storage.Commit(); err != nil {
			panic(err)
		}

		// Update account
		raw, err := rlp.EncodeToBytes(account)
		if err != nil {
			panic(err)
		}

		if err := state.stateTrie.Update(account.address[:], raw); err != nil {
			panic(err)
		}

		account.dirty = false
	}

	stateRootHash, err := state.stateTrie.Commit()
	if err != nil {
		log.Fatal(err)
	}

	txRootHash, err := state.txTrie.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return stateRootHash, txRootHash
}

// Revert state to last checkpoint
func (state *State) Revert() {
	t, err := trie.New(state.accountCheckpoint, state.db)
	if err != nil {
		panic(err)
	}
	state.stateTrie = t
	state.accounts = make(map[crypto.Address]*Account)
}
