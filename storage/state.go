package storage

import (
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/db"
	"github.com/QuoineFinancial/vertex/trie"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stellar/go/support/log"
)

// State is the global account state consisting of many address->state mapping
type State struct {
	trie     *trie.Trie
	accounts map[crypto.Address]*Account
}

var database = db.NewRocksDB("state.db") // TODO: Make this ENV

// GetState get the singleton state
func GetState(rootHash trie.Hash) *State {
	return &State{
		accounts: make(map[crypto.Address]*Account),
		trie:     trie.New(rootHash, database),
	}
}

// LoadAccount load the account from disk
func (state *State) LoadAccount(addr crypto.Address) (*Account, error) {
	raw, err := state.trie.Get(addr[:])
	if err != nil {
		return nil, err
	}
	var account Account
	if rlp.DecodeBytes(raw, &account); err != nil {
		return nil, err
	}
	account.address = addr
	account.storage = trie.New(account.StorageHash, database)
	account.contract = database.Get(account.ContractHash)
	return &account, nil
}

// GetAccount retrieve the account state at addr
func (state *State) GetAccount(addr crypto.Address) *Account {
	if state.accounts[addr] == nil {
		loadedAccount, err := state.LoadAccount(addr)
		if err != nil {
			log.Fatal(err)
		}
		state.accounts[addr] = loadedAccount
	}
	return state.accounts[addr]
}

// CreateAccount create a new account state for addr
func (state *State) CreateAccount(creator crypto.Address, addr crypto.Address, contract *[]byte) *Account {
	Account := newAccount(creator, addr, contract)
	state.accounts[addr] = Account
	return Account
}

// Commit stores all dirty Accounts to state.trie
func (state *State) Commit() (trie.Hash, error) {
	var err error
	for _, account := range state.accounts {
		if !account.dirty {
			continue
		}
		if account.StorageHash, err = account.storage.Commit(); err != nil {
			return trie.Hash{}, err
		}
		raw, _ := rlp.EncodeToBytes(account)
		if err = state.trie.Update(account.address[:], raw); err != nil {
			return trie.Hash{}, err
		}
		account.dirty = false
	}
	return state.trie.Commit()
}
