package storage

import (
	"errors"
	"time"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/db"
	"github.com/QuoineFinancial/liquid-chain/trie"
	"github.com/ethereum/go-ethereum/rlp"
)

// BlockInfo contains essential block information
type BlockInfo struct {
	Height  uint64
	AppHash trie.Hash
	Time    time.Time
}

// State is the global account state consisting of many address->state mapping
type State struct {
	BlockInfo  *BlockInfo
	db         db.Database
	trie       *trie.Trie
	checkpoint trie.Hash
	accounts   map[crypto.Address]*Account
}

// ErrAccountNotExist returns when loadAccount returns nil
var ErrAccountNotExist = errors.New("contract account not exist")

// New returns a state database
func New(root trie.Hash, db db.Database) (*State, error) {
	t, err := trie.New(root, db)
	if err != nil {
		return nil, err
	}
	return &State{
		db:         db,
		trie:       t,
		checkpoint: root,
		accounts:   make(map[crypto.Address]*Account),
	}, nil
}

// LoadAccount load the account from disk
func (state *State) LoadAccount(address crypto.Address) (*Account, error) {
	raw, err := state.trie.Get(address[:])
	if err != nil {
		return nil, err
	}
	var account Account
	if len(raw) <= 0 {
		return nil, nil
	}
	if err := rlp.DecodeBytes(raw, &account); err != nil {
		return nil, err
	}
	account.address = address
	account.contract = state.db.Get(account.ContractHash)
	if account.storage, err = trie.New(account.StorageHash, state.db); err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccount retrieve the account state at addr
func (state *State) GetAccount(address crypto.Address) (*Account, error) {
	if state.accounts[address] == nil {
		loadedAccount, err := state.LoadAccount(address)
		if err != nil {
			return nil, err
		}
		if loadedAccount == nil {
			return nil, ErrAccountNotExist
		}
		state.accounts[address] = loadedAccount
	}
	return state.accounts[address], nil
}

// CreateAccount create a new account state for addr
func (state *State) CreateAccount(creator crypto.Address, address crypto.Address, contract []byte) (*Account, error) {
	storage, err := trie.New(trie.Hash{}, state.db)
	if err != nil {
		return nil, err
	}
	account := &Account{
		Nonce:    0,
		Creator:  creator,
		address:  address,
		storage:  storage,
		contract: contract,
		dirty:    true,
	}

	account.setContract(contract)
	state.db.Put(account.ContractHash[:], account.contract)

	state.accounts[address] = account
	return account, nil
}

// Commit stores all dirty Accounts to state.trie
func (state *State) Commit() trie.Hash {
	var err error
	for _, account := range state.accounts {
		if account == nil || !account.dirty {
			continue
		}

		// Add source code

		// Update account storage
		if account.StorageHash, err = account.storage.Commit(); err != nil {
			panic(err)
		}

		// Update account
		raw, _ := rlp.EncodeToBytes(account)
		if err = state.trie.Update(account.address[:], raw); err != nil {
			panic(err)
		}
		account.dirty = false
	}
	state.checkpoint, err = state.trie.Commit()
	if err != nil {
		panic(err)
	}
	return state.checkpoint
}

// Revert state to last checkpoint
func (state *State) Revert() {
	t, err := trie.New(state.checkpoint, state.db)
	if err != nil {
		panic(err)
	}
	state.trie = t
	state.accounts = make(map[crypto.Address]*Account)
}
