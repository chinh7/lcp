package storage

import (
	"github.com/QuoineFinancial/vertex/abi"
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/trie"
	"golang.org/x/crypto/sha3"
)

// Account stores information related to the account
type Account struct {
	Nonce        uint64
	ContractHash []byte
	StorageHash  trie.Hash // merkle root of the storage trie
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

// GetContract retrieves contract code for account state
func (account *Account) GetContract() (*abi.Contract, error) {
	return abi.DecodeContract(account.contract)
}

// SetNonce stores the latest nonce to account state
func (account *Account) SetNonce(nonce uint64) {
	account.Nonce = nonce
}

func (account *Account) setContract(contract []byte) {
	account.contract = contract
	contractHash := sha3.Sum256(contract)
	database.Put(contractHash[:], contract)
	account.ContractHash = contractHash[:]
}

func newAccount(creator crypto.Address, address crypto.Address, contract *[]byte) *Account {
	account := &Account{
		Nonce:    0,
		Creator:  creator,
		address:  address,
		storage:  trie.New(trie.Hash{}, database),
		dirty:    true,
		contract: []byte{},
	}
	if contract != nil {
		account.setContract(*contract)
	}
	return account
}
