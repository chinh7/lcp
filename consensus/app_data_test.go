package consensus

import (
	"crypto/ed25519"
	"math/rand"

	"github.com/QuoineFinancial/liquid-chain/constant"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/util"
)

type testResource struct{}

func (testResource) getDeployTx() *crypto.Transaction {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(0),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	data, err := util.BuildDeployTxPayload("./execution_testdata/contract.wasm", "./execution_testdata/contract-abi.json", "", []string{})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 1,
		Receipt:  &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (testResource) getInvokeTx() *crypto.Transaction {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(1),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxData("./execution_testdata/contract-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 0),
		GasLimit: 0,
		GasPrice: 1,
		Receipt:  &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (testResource) getInvalidMaxSizeTx() *crypto.Transaction {
	seed := make([]byte, 32)
	rand.Read(seed)

	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(0),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	type maxSizeContart [constant.MaxTransactionSize]byte
	var contract maxSizeContart
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  &crypto.TxPayload{Contract: contract[:]},
		Receiver: crypto.EmptyAddress,
		GasLimit: 0,
		GasPrice: 1,
		Receipt:  &crypto.TxReceipt{},
	}
	return tx
}

func (testResource) getInvaliSignatureTx() *crypto.Transaction {
	seed := make([]byte, 32)
	rand.Read(seed)

	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(0),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxData("./execution_testdata/contract-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 0),
		GasLimit: 0,
		GasPrice: 1,
		Receipt:  &crypto.TxReceipt{},
	}
	tx.Signature = []byte{1, 2, 3}
	return tx
}

func (testResource) getInvalidNonceTx() *crypto.Transaction {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(123),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxData("./execution_testdata/contract-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 0),
		GasLimit: 0,
		GasPrice: 1,
		Receipt:  &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (testResource) getInvalidGasPriceTx() *crypto.Transaction {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(2),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxData("./execution_testdata/contract-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 0),
		GasLimit: 0,
		GasPrice: 0,
		Receipt:  &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (testResource) getInvokeNilContractTx() *crypto.Transaction {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(2),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxData("./execution_testdata/contract-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: crypto.NewDeploymentAddress(senderAddress, 123),
		GasLimit: 0,
		GasPrice: 0,
		Receipt:  &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}

func (testResource) getInvokeNonContractTx() *crypto.Transaction {
	seed := make([]byte, 32)
	privateKey := ed25519.NewKeyFromSeed(seed)
	sender := crypto.TxSender{
		Nonce:     uint64(2),
		PublicKey: privateKey.Public().(ed25519.PublicKey),
	}
	senderAddress := crypto.AddressFromPubKey(sender.PublicKey)
	data, err := util.BuildInvokeTxData("./execution_testdata/contract-abi.json", "mint", []string{"1000"})
	if err != nil {
		panic(err)
	}
	tx := &crypto.Transaction{
		Version:  1,
		Sender:   &sender,
		Payload:  data,
		Receiver: senderAddress,
		GasLimit: 0,
		GasPrice: 0,
		Receipt:  &crypto.TxReceipt{},
	}
	dataToSign := crypto.GetSigHash(tx)
	tx.Signature = crypto.Sign(privateKey, dataToSign.Bytes())
	return tx
}
