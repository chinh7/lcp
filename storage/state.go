package storage

import (
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/db"
	"github.com/QuoineFinancial/vertex/trie"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stellar/go/support/log"
)

// AccountState stores information related to the account
type AccountState struct {
	Nonce       uint64
	CodeHash    []byte
	StorageHash trie.Hash // merkle root of the storage trie

	dirty   bool
	address crypto.Address
	// contract execution storage
	storage *trie.Trie
	// contract code
	code []byte
}

// State is the global account state consisting of many address->state mapping
type State struct {
	trie          *trie.Trie
	accountStates map[crypto.Address]*AccountState
}

var state *State
var database = db.NewRocksDB("state.db") // TODO: Make this ENV

// GetState get the singleton state
func GetState() *State {
	if state == nil {
		state = &State{
			accountStates: make(map[crypto.Address]*AccountState),
			trie:          trie.New(trie.Hash{}, database),
		}
	}
	return state
}

// LoadAccountState load the account from disk
// TODO: Incomplete
func (state *State) LoadAccountState(addr crypto.Address) (*AccountState, error) {
	raw, err := state.trie.Get(addr[:])
	if err != nil {
		return nil, err
	}
	var account AccountState
	if rlp.DecodeBytes(raw, &account); err != nil {
		return nil, err
	}
	account.address = addr
	account.storage = trie.New(account.StorageHash, database)
	account.code = database.Get(account.CodeHash)
	return &account, nil
}

// GetAccountState retrieve the account state at addr
func (state *State) GetAccountState(addr crypto.Address) *AccountState {
	if state.accountStates[addr] == nil {
		loadedAccount, err := state.LoadAccountState(addr)
		if err != nil {
			log.Fatal(err)
		}
		state.accountStates[addr] = loadedAccount
	}
	return state.accountStates[addr]
}

// CreateAccountState create a new account state for addr
func (state *State) CreateAccountState(addr crypto.Address) *AccountState {
	accountState := newAccountState(addr)
	accountState.SetNonce(0)
	state.accountStates[addr] = accountState
	return accountState
}

// SetCode store contract code to the account state at addr
func (state *State) SetCode(addr crypto.Address, code []byte) {
	accountState := state.GetAccountState(addr)
	if accountState != nil {
		accountState.SetCode(code)
	}
}

// StorageGet retrieve the data stored at key in addr storage
func (state *State) StorageGet(addr crypto.Address, key [32]byte) []byte {
	accountState := state.GetAccountState(addr)
	if accountState != nil {
		if result, err := accountState.storage.Get(key[:]); err == nil {
			// TODO: Handle err
			return result
		}
	}
	return nil
}

// StorageSet save the data to addr storage
func (state *State) StorageSet(addr crypto.Address, key [32]byte, value []byte) {
	accountState := state.GetAccountState(addr)
	if accountState != nil {
		accountState.storage.Update(key[:], value) // TODO: Handle err
		accountState.dirty = true
	}
}

// SetCode store contract code to the account state
func (state *AccountState) SetCode(code []byte) {
	state.code = code
}

// GetAddress returns state address
func (state *AccountState) GetAddress() crypto.Address {
	return state.address
}

// GetCode retrieves contract code for account state
func (state *AccountState) GetCode() []byte {
	return state.code
}

// SetNonce stores the latest nonce to account state
func (state *AccountState) SetNonce(nonce uint64) {
	state.Nonce = nonce
}

func newAccountState(address crypto.Address) *AccountState {
	return &AccountState{
		address: address,
		storage: trie.New(trie.Hash{}, database),
		dirty:   true,
	}
}

// Commit stores all dirty Accounts to state.trie
func (state *State) Commit() (trie.Hash, error) {
	for _, account := range state.accountStates {
		if !account.dirty {
			continue
		}
		account.StorageHash = account.storage.Hash()
		raw, _ := rlp.EncodeToBytes(account)
		if err := state.trie.Update(account.address[:], raw); err != nil {
			return trie.Hash{}, err
		}
	}
	hash, err := state.trie.Commit()
	if err != nil {
		return trie.Hash{}, err
	}
	return hash, nil
}
