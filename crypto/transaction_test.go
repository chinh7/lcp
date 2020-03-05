package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/ed25519"
)

func TestTxSerialization(t *testing.T) {
	toAddress, err := AddressFromString("LB3Z6N6HTFUPQ573QENJ4OCFFUPENY2EW7ZHQZSSIO4AODT3HHE53N52")
	if err != nil {
		panic(err)
	}
	txSigner := TxSigner{Nonce: 10}
	tx := &Tx{To: toAddress, From: txSigner, GasLimit: 100}
	txRecouped := &Tx{}
	txRecouped.Deserialize(tx.Serialize())
	if tx.String() != txRecouped.String() {
		t.Errorf("Expect deserialization to produce the same value, expected: %s, got %s", tx.String(), txRecouped.String())
	}
}

func TestTxDataSerialization(t *testing.T) {
	params := []byte{0, 0, 1, 1, 0, 1, 1}
	var txDataRecouped TxData

	txData := TxData{Method: "method", Params: params, ContractCode: []byte{}}
	txDataRecouped.Deserialize(txData.Serialize())
	if diff := cmp.Diff(txData, txDataRecouped); diff != "" {
		t.Errorf("Decoding of %v is incorrect, expected: %v, got: %v, diff: %v", txData, txData, txDataRecouped, diff)
	}
}

func TestTxSignerSerialization(t *testing.T) {
	txSigner := &TxSigner{Nonce: 10}
	txSignerRecouped := &TxSigner{}
	txSignerRecouped.Deserialize(txSigner.Serialize())
	if txSigner.String() != txSignerRecouped.String() {
		t.Errorf("Expect deserialization to produce the same value, expected: %s, got %s", txSigner.String(), txSignerRecouped.String())
	}
}

func TestTxSignature(t *testing.T) {
	pubkey, prvkey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	_, invalidPrv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name   string
		want   bool
		prvkey ed25519.PrivateKey
	}{
		{
			name:   "valid private key",
			prvkey: prvkey,
			want:   true,
		},
		{
			name:   "invalid private key",
			prvkey: invalidPrv,
			want:   false,
		},
	}

	if err != nil {
		panic(err)
	}

	tx := &Tx{From: TxSigner{}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tx.Sign(tt.prvkey); err != nil {
				panic(err)
			}
			tx.From.PubKey = pubkey
			if tx.SigVerified() != tt.want {
				t.Errorf("Tx.SigVerified() = %v, want %v", tx.SigVerified(), tt.want)
			}
		})
	}

}

func TestCreateAddress(t *testing.T) {
	prvkey := ed25519.NewKeyFromSeed([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31})
	pubkey := prvkey[32:]
	txSigner := TxSigner{Nonce: 7, PubKey: pubkey}
	generatedAddress := txSigner.CreateAddress()
	expected := "LD6QC6YFHIZ6DCN452PJULI5WHVKYRCNIJHKRQUOF6QOKQ6WPFGT5CYW"
	if generatedAddress.String() != expected {
		t.Errorf("CreateAddress = %v, want %v", generatedAddress.String(), expected)
	}

}

func TestTxGetFee(t *testing.T) {
	type fields struct {
		GasLimit uint32
		GasPrice uint32
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint64
		wantErr string
	}{
		{
			name:   "zero gas",
			fields: fields{GasLimit: 0, GasPrice: 0},
			want:   0,
		},
		{
			name:   "valid fee",
			fields: fields{GasLimit: 250, GasPrice: 5},
			want:   250 * 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Tx{
				GasLimit: tt.fields.GasLimit,
				GasPrice: tt.fields.GasPrice,
			}
			got, err := tx.GetFee()
			if err != nil && err.Error() != tt.wantErr {
				t.Errorf("Tx.GetFee() error = %v, wantErr %v", err.Error(), tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Tx.GetFee() = %v, want %v", got, tt.want)
			}
		})
	}
}
