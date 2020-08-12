package storage

import (
	"log"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/trie"
)

// StateStorage is the global account state consisting of many address->state mapping
type StateStorage struct {
	db.Database
	blockHeader       *crypto.BlockHeader
	txTrie            *trie.Trie
	stateTrie         *trie.Trie
	accounts          map[crypto.Address]*Account
	accountCheckpoint common.Hash
}

// AddTransaction add new tx to txTrie
func (state *StateStorage) AddTransaction(tx *crypto.Transaction) error {
	rawTx, err := tx.Serialize()
	if err != nil {
		return err
	}
	return state.txTrie.Update(tx.Hash().Bytes(), rawTx)
}

// NewStateStorage returns a state storage
func NewStateStorage(db db.Database) *StateStorage {
	return &StateStorage{Database: db}
}

// MustLoadState do LoadState, but panic if error
func (state *StateStorage) MustLoadState(blockHeader *crypto.BlockHeader) {
	if err := state.LoadState(blockHeader); err != nil {
		panic(err)
	}
}

// LoadState load state root of blockHeader into trie
func (state *StateStorage) LoadState(blockHeader *crypto.BlockHeader) error {
	stateTrie, err := trie.New(blockHeader.StateRoot, state.Database)
	if err != nil {
		return err
	}

	txTrie, err := trie.New(blockHeader.TransactionRoot, state.Database)
	if err != nil {
		return err
	}

	state.blockHeader = blockHeader
	state.txTrie = txTrie
	state.stateTrie = stateTrie
	state.accountCheckpoint = blockHeader.StateRoot
	state.accounts = make(map[crypto.Address]*Account)

	return nil
}

// GetBlockHeader return header of block that inits current state
func (state *StateStorage) GetBlockHeader() *crypto.BlockHeader {
	return state.blockHeader
}

// Hash retrive hash of entire state
func (state *StateStorage) Hash() common.Hash {
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
func (state *StateStorage) Commit() (common.Hash, common.Hash) {
	var err error
	for _, account := range state.accounts {
		if account == nil || !account.dirty {
			continue
		}

		if account.IsContract() {
			// Update contract
			state.Put(account.ContractHash[:], account.contract)
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

	log.Println("StateStorage.Commit successfully for block", state.blockHeader.Height)
	return stateRootHash, txRootHash
}

// Revert state to last checkpoint
func (state *StateStorage) Revert() {
	t, err := trie.New(state.accountCheckpoint, state.Database)
	if err != nil {
		panic(err)
	}
	state.stateTrie = t
	state.accounts = make(map[crypto.Address]*Account)
}
