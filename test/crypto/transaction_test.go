package test

import (
	"testing"

	"github.com/QuoineFinancial/vertex/crypto"
)

func TestTransaction(t *testing.T) {
	toAddress := crypto.AddressFromString("LB3Z6N6HTFUPQ573QENJ4OCFFUPENY2EW7ZHQZSSIO4AODT3HHE53N52")
	txSigner := crypto.TxSigner{Nonce: 10}
	tx := &crypto.Tx{To: toAddress, From: txSigner}
	txRecouped := &crypto.Tx{}
	txRecouped.Deserialize(tx.Serialize())
	if tx.String() != txRecouped.String() {
		t.Errorf("Expect deserialization to produce the same value, expected: %s, got %s", tx.String(), txRecouped.String())
	}
}

func TestTxData(t *testing.T) {
	txData := &crypto.TxData{Method: "method"}
	txDataRecouped := &crypto.TxData{}
	txDataRecouped.Deserialize(txData.Serialize())
	if txData.Method != txDataRecouped.Method {
		t.Errorf("Expect deserialization to produce the same value, expected: %s, got %s", txData.Method, txDataRecouped.Method)
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
