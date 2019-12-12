package test

import (
	"testing"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/google/go-cmp/cmp"
)

func TestTransaction(t *testing.T) {
	toAddress := crypto.AddressFromString("LB3Z6N6HTFUPQ573QENJ4OCFFUPENY2EW7ZHQZSSIO4AODT3HHE53N52")
	txSigner := crypto.TxSigner{Nonce: 10}
	tx := &crypto.Tx{To: toAddress, From: txSigner, GasLimit: 100}
	txRecouped := &crypto.Tx{}
	txRecouped.Deserialize(tx.Serialize())
	if tx.String() != txRecouped.String() {
		t.Errorf("Expect deserialization to produce the same value, expected: %s, got %s", tx.String(), txRecouped.String())
	}
}

func TestTxData(t *testing.T) {
	params := []byte{0, 0, 1, 1, 0, 1, 1}
	var txDataRecouped crypto.TxData

	txData := crypto.TxData{Method: "method", Params: params}
	txDataRecouped.Deserialize(txData.Serialize())
	if diff := cmp.Diff(txData, txDataRecouped); diff != "" {
		t.Errorf("Decoding of %v is incorrect, expected: %v, got: %v, diff: %v", txData, txData, txDataRecouped, diff)
	}
}

func TestTxSigner(t *testing.T) {
	txSigner := &crypto.TxSigner{Nonce: 10}
	txSignerRecouped := &crypto.TxSigner{}
	txSignerRecouped.Deserialize(txSigner.Serialize())
	if txSigner.String() != txSignerRecouped.String() {
		t.Errorf("Expect deserialization to produce the same value, expected: %s, got %s", txSigner.String(), txSignerRecouped.String())
	}
}
