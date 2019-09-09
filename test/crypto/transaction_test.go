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
	var params []interface{}
	params = append(params, "LACWIGXH6CZCRRHFSK2F4BINXGUGUS2FSX5GSYG3RMP5T55EV72DHAJ1")
	params = append(params, "100")
	var txDataRecouped crypto.TxData

	txData := &crypto.TxData{Method: "method", Params: params}
	txDataRecouped.Deserialize(txData.Serialize())
	if txData.Method != txDataRecouped.Method {
		t.Errorf("Expect deserialization to produce the same value, expected: %s, got %s", txData.Method, txDataRecouped.Method)
	}
	for i, v := range txDataRecouped.Params {
		if string(v.([]byte)) != txData.Params[i] {
			t.Errorf("Expect deserialization to produce the same value, expected: %s, got %s", txData.Params[0], string(v.([]byte)))
		}
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
