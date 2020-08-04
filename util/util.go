package util

import (
	"fmt"
	"io/ioutil"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// BuildInvokeTxData build data for invoke transaction
func BuildInvokeTxData(headerPath string, methodName string, params []string) (*crypto.TxPayload, error) {
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

	return &crypto.TxPayload{
		Method: methodName,
		Params: encodedArgs,
	}, nil
}

// BuildDeployTxPayload build data for deploy transaction
func BuildDeployTxPayload(codePath string, headerPath string, initFuncName string, params []string) (*crypto.TxPayload, error) {
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

	payload := crypto.TxPayload{
		Contract: contractCode,
	}

	function, err := header.GetFunction(initFuncName)
	if err == nil {
		encodedArgs, err := abi.EncodeFromString(function.Parameters, params)
		if err != nil {
			return nil, err
		}
		payload.Method = initFuncName
		payload.Params = encodedArgs
	} else if err.Error() != fmt.Sprintf("function %s not found", initFuncName) {
		return nil, err
	}

	return &payload, nil
}
