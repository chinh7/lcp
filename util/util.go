package util

import (
	"fmt"
	"io/ioutil"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
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
func BuildDeployTxData(codePath string, headerPath string, initFuncName string, params []string) ([]byte, error) {
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

	contractCode, err := rlp.EncodeToBytes(&abi.Contract{Header: header, Code: code})
	if err != nil {
		return nil, err
	}

	txData := crypto.TxData{ContractCode: contractCode}

	function, err := header.GetFunction(initFuncName)
	if err == nil {
		encodedArgs, err := abi.EncodeFromString(function.Parameters, params)
		if err != nil {
			panic(err)
		}
		txData.Method = initFuncName
		txData.Params = encodedArgs
	} else if err.Error() != fmt.Sprintf("function %s not found", initFuncName) {
		panic(err)
	}

	return txData.Serialize(), nil
}
