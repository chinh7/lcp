package storage

import (
	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/QuoineFinancial/vertex/trie"
	"golang.org/x/crypto/sha3"
)

// Account stores information related to the account
type Account struct {
	Nonce       uint64
	CodeHash    []byte
	StorageHash trie.Hash // merkle root of the storage trie

	dirty   bool
	address crypto.Address
	storage *trie.Trie
	code    []byte
}

// GetStorage get the value at key of storage
func (account *Account) GetStorage(key []byte) ([]byte, error) {
	return account.storage.Get(key)
}

// SetStorage set the account storage
func (account *Account) SetStorage(key, value []byte) error {
	return account.storage.Update(key, value)
}

// GetAddress returns state address
func (account *Account) GetAddress() crypto.Address {
	return account.address
}

// GetCode retrieves contract code for account state
func (account *Account) GetCode() []byte {
	return account.code
}

// SetNonce stores the latest nonce to account state
func (account *Account) SetNonce(nonce uint64) {
	account.Nonce = nonce
}

func (account *Account) setCode(code []byte) {
	account.code = code
	codeHash := sha3.Sum256(code)
	database.Put(codeHash[:], code)
	account.CodeHash = codeHash[:]
}

func newAccount(address crypto.Address, code *[]byte) *Account {
	account := &Account{
		Nonce:   0,
		address: address,
		storage: trie.New(trie.Hash{}, database),
		dirty:   true,
		code:    []byte{},
	}
	if code != nil {
		account.setCode(*code)
	}
	return account
}
