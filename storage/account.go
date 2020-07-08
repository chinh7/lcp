package storage

import (
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/common"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/trie"
	"golang.org/x/crypto/blake2b"
)

// Account stores information related to the account
type Account struct {
	Nonce        uint64
	ContractHash []byte
	StorageHash  common.Hash // merkle root of the storage trie
	Creator      crypto.Address

	dirty    bool
	address  crypto.Address
	storage  *trie.Trie
	contract []byte
}

// GetStorage get the value at key of storage
func (account *Account) GetStorage(key []byte) ([]byte, error) {
	return account.storage.Get(key)
}

// SetStorage set the account storage
func (account *Account) SetStorage(key, value []byte) error {
	account.dirty = true
	return account.storage.Update(key, value)
}

// GetAddress returns state address
func (account *Account) GetAddress() crypto.Address {
	return account.address
}

// IsContract check whether this is an contract account or a normal account
func (account *Account) IsContract() bool {
	return len(account.ContractHash) > 0
}

// GetContract retrieves contract code for account state
func (account *Account) GetContract() (*abi.Contract, error) {
	return abi.DecodeContract(account.contract)
}

// SetNonce stores the latest nonce to account state
func (account *Account) SetNonce(nonce uint64) {
	account.dirty = true
	account.Nonce = nonce
}

// GetCreator contract creator
func (account *Account) GetCreator() crypto.Address {
	return account.Creator
}

func (account *Account) setContract(contract []byte) {
	account.dirty = true
	account.contract = contract
	if len(account.contract) > 0 {
		contractHash := blake2b.Sum256(contract)
		account.ContractHash = contractHash[:]
	}
}
