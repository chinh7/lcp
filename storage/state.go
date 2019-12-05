package storage

import (
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/db"
	"github.com/QuoineFinancial/vertex/trie"
	"github.com/ethereum/go-ethereum/rlp"
)

// State is the global account state consisting of many address->state mapping
type State struct {
	db         db.Database
	trie       *trie.Trie
	checkpoint trie.Hash
	accounts   map[crypto.Address]*Account
}

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
	if raw == nil {
		return nil, nil
	}
	var account Account
	if rlp.DecodeBytes(raw, &account); err != nil {
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
func (state *State) Commit() (trie.Hash, error) {
	var err error
	for _, account := range state.accounts {
		if account == nil || !account.dirty {
			continue
		}

		// Add source code

		// Update account storage
		if account.StorageHash, err = account.storage.Commit(); err != nil {
			return trie.Hash{}, err
		}

		// Update account
		raw, _ := rlp.EncodeToBytes(account)
		if err = state.trie.Update(account.address[:], raw); err != nil {
			return trie.Hash{}, err
		}
		account.dirty = false
	}
	state.checkpoint, err = state.trie.Commit()
	return state.checkpoint, err
}

// Revert state to last checkpoint
func (state *State) Revert() error {
	t, err := trie.New(state.checkpoint, state.db)
	if err != nil {
		return err
	}
	state.trie = t
	state.accounts = make(map[crypto.Address]*Account)
	return nil
}
