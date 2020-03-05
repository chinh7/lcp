package util

import (
	"io/ioutil"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// BuildInvokeTxData build data for invoke transaction
func BuildInvokeTxData(headerPath string, methodName string, params []string) ([]byte, error) {
	header, err := abi.LoadHeaderFromFile(headerPath)
	if err != nil {
		return nil, err
	}

	function, err := header.GetFunction(methodName)
	if err != nil {
		return nil, err
	}

	encodedArgs, err := abi.EncodeFromString(function.Parameters, params)
	if err != nil {
		return nil, err
	}

	txData := crypto.TxData{Method: methodName, Params: encodedArgs}
	return txData.Serialize(), nil
}

// BuildDeployTxData build data for deploy transaction
func BuildDeployTxData(codePath string, headerPath string) ([]byte, error) {
	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		return nil, err
	}

	encodedHeader, err := abi.EncodeHeaderToBytes(headerPath)
	if err != nil {
		return nil, err
	}

	header, err := abi.DecodeHeader(encodedHeader)
	if err != nil {
		return nil, err
	}

	data, err := rlp.EncodeToBytes(&abi.Contract{Header: header, Code: code})
	if err != nil {
		return nil, err
	}

	return data, nil
}
