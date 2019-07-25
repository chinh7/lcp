package storage

import (
	"github.com/vertexdlt/vertex/crypto"
)

// Account Nonce + Merkel hash
type Account struct {
	Nonce    uint64
	CodeHash []byte
	// Root common.Hash // merkle root of the storage trie
}

// AccountState stores information related to the account
type AccountState struct {
	Address crypto.Address
	// contract execution storage
	storage map[[32]byte][]byte
	// account information
	account Account
	// contract code
	code []byte
}

// State is the global account state consisting of many address->state mapping
type State struct {
	accountStates map[crypto.Address]*AccountState
}

var state *State

// GetState get the singleton state
func GetState() *State {
	if state == nil {
		state = &State{accountStates: make(map[crypto.Address]*AccountState)}
	}
	return state
}

// GetAccountState retrieve the account state at addr
func (state *State) GetAccountState(addr crypto.Address) *AccountState {
	return state.accountStates[addr]
}

// CreateAccountState create a new account state for addr
func (state *State) CreateAccountState(addr crypto.Address) *AccountState {
	accountState := newAccountState(addr, Account{})
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
		return accountState.storage[key]
	}
	return nil
}

// StorageSet save the data to addr storage
func (state *State) StorageSet(addr crypto.Address, key [32]byte, value []byte) {
	accountState := state.GetAccountState(addr)
	if accountState != nil {
		accountState.storage[key] = value
	}
}

// SetCode store contract code to the account state
func (state *AccountState) SetCode(code []byte) {
	state.code = code
}

// GetCode retrieves contract code for account state
func (state *AccountState) GetCode() []byte {
	return state.code
}

// SetNonce stores the latest nonce to account state
func (state *AccountState) SetNonce(nonce uint64) {
	state.account.Nonce = nonce
}

func newAccountState(address crypto.Address, account Account) *AccountState {
	return &AccountState{
		Address: address,
		account: account,
		storage: make(map[[32]byte][]byte),
	}
}
