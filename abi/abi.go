package abi

import (
	"errors"
	"fmt"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
)

// Encode return []byte from an inputted params and values pair
func Encode(params []*Parameter, values []interface{}) ([]byte, error) {
	if len(params) != len(values) {
		return []byte{0}, fmt.Errorf("Parameter count mismatch, expecting: %d, got: %d", len(params), len(values))
	}

	var rlpCompatibleArgs []interface{}

	for index, param := range params {
		if param.IsArray {
			arrayArg, err := param.Type.NewArrayArgument(values[index])
			if err != nil {
				return nil, err
			}
			rlpCompatibleArgs = append(rlpCompatibleArgs, arrayArg)
		} else {
			argument, err := param.Type.NewArgument(values[index])
			if err != nil {
				return nil, err
			}
			rlpCompatibleArgs = append(rlpCompatibleArgs, argument)
		}
	}
	result, err := rlp.EncodeToBytes(rlpCompatibleArgs)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeBytes returns uint64 array compatible with VM
func DecodeBytes(params []*Parameter, bytes []byte) ([][]byte, error) {
	var decoded []interface{}
	if err := rlp.DecodeBytes(bytes, &decoded); err != nil {
		return nil, err
	}

	var result [][]byte
	for i, in := range decoded {
		var buffer []byte
		if params[i].IsArray {
			arrArgs := in.([]interface{})
			for _, arg := range arrArgs {
				buffer = append(buffer, arg.([]byte)...)
			}
		} else {
			buffer = append(buffer, in.([]byte)...)
		}
		result = append(result, buffer)
	}

	return result, nil
}

// EncodeToBytes encodes arguments in byte format - an inverse of DecodeToBytes
func EncodeToBytes(params []*Parameter, bytes [][]byte) ([]byte, error) {
	var rlpCompatibleArgs []interface{}

	for index, param := range params {
		if param.IsArray {
			elementSize := param.Type.GetMemorySize()
			arrayBytes := bytes[index]
			if len(arrayBytes)%elementSize != 0 {
				return nil, errors.New("misaligned array byte size")
			}
			argsCount := len(arrayBytes) / elementSize
			var arrayArgs [][]byte
			var offset int

			for i := 0; i < argsCount; i++ {
				arg := arrayBytes[offset : offset+elementSize]
				offset += elementSize
				arrayArgs = append(arrayArgs, arg)
			}

			rlpCompatibleArgs = append(rlpCompatibleArgs, arrayArgs)
		} else {
			rlpCompatibleArgs = append(rlpCompatibleArgs, bytes[index])
		}
	}
	return rlp.EncodeToBytes(rlpCompatibleArgs)
}
