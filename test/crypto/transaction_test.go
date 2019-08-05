package test

import (
	"testing"

	"github.com/QuoineFinancial/vertex/crypto"
)

func TestTransaction(t *testing.T) {
	address := crypto.AddressFromString("LB3Z6N6HTFUPQ573QENJ4OCFFUPENY2EW7ZHQZSSIO4AODT3HHE53N52")
	tx := &crypto.Tx{To: address}
	txRecouped := &crypto.Tx{}
	txRecouped.Deserialize(tx.Serialize())
	if txRecouped.To.String() != tx.To.String() {
		t.Error("Expect deserialization to produce the same value, got {}", txRecouped.To.String())
	}

}
