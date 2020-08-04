package storage

import "github.com/QuoineFinancial/liquid-chain/crypto"

// AddTransaction add new tx to txTrie
func (state *State) AddTransaction(tx *crypto.Transaction) error {
	rawTx, err := tx.Serialize()
	if err != nil {
		return err
	}
	return state.txTrie.Update(tx.Hash().Bytes(), rawTx)
}
