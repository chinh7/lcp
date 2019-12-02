package abi

import (
	"fmt"
	"math"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// PrimitiveType PrimitiveType
type PrimitiveType uint

// enum for types
const (
	Uint8   PrimitiveType = 0x0
	Uint16  PrimitiveType = 0x1
	Uint32  PrimitiveType = 0x2
	Uint64  PrimitiveType = 0x3
	Int8    PrimitiveType = 0x4
	Int16   PrimitiveType = 0x5
	Int32   PrimitiveType = 0x6
	Int64   PrimitiveType = 0x7
	Float32 PrimitiveType = 0x8
	Float64 PrimitiveType = 0x9
	Address PrimitiveType = 0xa
)

// IsPointer return whether p is pointer or not
func (p PrimitiveType) IsPointer() bool {
	switch p {
	case Address:
		return true
	default:
		return false
	}
}

func (p PrimitiveType) String() string {
	return []string{"uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "address"}[p]
}

// GetMemorySize returns memory size for a primitive type
func (p PrimitiveType) GetMemorySize() (int, error) {
	switch p {
	case Address:
		return crypto.AddressLength, nil
	case Uint8, Int8:
		return 1, nil
	case Uint16, Int16:
		return 2, nil
	case Uint32, Int32, Float32:
		return 4, nil
	case Uint64, Int64, Float64:
		return 8, nil
	default:
		return 0, fmt.Errorf("Not supported type")
	}
}

// newPrimitiveArg parse type string and value into PrimitiveArg
func newPrimitiveArg(t PrimitiveType, value interface{}) (interface{}, error) {
	var parsedValue interface{}
	switch t {
	case Address, Uint8, Uint16, Uint32, Uint64:
		parsedValue = value
	case Int8:
		parsedValue = uint8(value.(int8))
	case Int16:
		parsedValue = uint16(value.(int16))
	case Int32:
		parsedValue = uint32(value.(int32))
	case Int64:
		parsedValue = uint64(value.(int64))
	case Float32:
		parsedValue = math.Float32bits(value.(float32))
	case Float64:
		parsedValue = math.Float64bits(value.(float64))
	default:
		return nil, fmt.Errorf("not supported type: %s", t)
	}
	return parsedValue, nil
}

// parseArrayArg parse type and values into ArrayArg
func parseArrayArg(t PrimitiveType, value interface{}) ([]interface{}, error) {
	var parsedArgs []interface{}

	switch t {
	case Address:
		parsed, ok := value.([]crypto.Address)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Uint8:
		parsed, ok := value.([]uint8)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Uint16:
		parsed, ok := value.([]uint16)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Uint32:
		parsed, ok := value.([]uint32)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Uint64:
		parsed, ok := value.([]uint64)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Int8:
		parsed, ok := value.([]int8)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Int16:
		parsed, ok := value.([]int16)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Int32:
		parsed, ok := value.([]int32)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Int64:
		parsed, ok := value.([]int64)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Float32:
		parsed, ok := value.([]float32)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	case Float64:
		parsed, ok := value.([]float64)
		if !ok {
			return nil, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			arg, err := newPrimitiveArg(t, p)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, arg)
		}
	default:
		return nil, fmt.Errorf("not supported type: %s", t)
	}

	return parsedArgs, nil
}

func reverseByte(input []byte) []byte {
	reversed := []byte{}
	for index := len(input) - 1; index >= 0; index-- {
		reversed = append(reversed, input[index])
	}
	return reversed
}

func convertToLittleEndian(t PrimitiveType, bytes []byte) []byte {
	var buffer []byte
	switch t.String() {
	case "address":
		buffer = make([]byte, 35)
		copy(buffer, bytes)
	case "uint8", "int8":
		buffer = make([]byte, 1)
		copy(buffer, reverseByte(bytes))
	case "uint16", "int16":
		buffer = make([]byte, 2)
		copy(buffer, reverseByte(bytes))
	case "uint32", "int32", "float32":
		buffer = make([]byte, 4)
		copy(buffer, reverseByte(bytes))
	case "uint64", "int64", "float64":
		buffer = make([]byte, 8)
		copy(buffer, reverseByte(bytes))
	}
	return buffer
}

// Encode return []byte from an inputted params and values pair
func Encode(params []*Parameter, values []interface{}) ([]byte, error) {
	if len(params) != len(values) {
		return []byte{0}, fmt.Errorf("Parameter count mismatch, expecting: %d, got: %d", len(params), len(values))
	}
	result := []byte{}

	var rlpCompatibleArgs []interface{}

	for index, param := range params {
		if param.IsArray {
			arrayArg, err := parseArrayArg(param.Type, values[index])
			if err != nil {
				return nil, err
			}
			rlpCompatibleArgs = append(rlpCompatibleArgs, arrayArg)
		} else {
			arrayArg, err := newPrimitiveArg(param.Type, values[index])
			if err != nil {
				return nil, err
			}
			rlpCompatibleArgs = append(rlpCompatibleArgs, arrayArg)
		}
	}
	result, err := rlp.EncodeToBytes(rlpCompatibleArgs)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DecodeToBytes returns uint64 array compatible with VM
func DecodeToBytes(params []*Parameter, bytes []byte) ([][]byte, error) {
	var decoded []interface{}
	err := rlp.DecodeBytes(bytes, &decoded)
	if err != nil {
		return nil, err
	}

	var result [][]byte
	for i, in := range decoded {
		var buffer []byte
		if params[i].IsArray {
			arrArgs := in.([]interface{})
			for _, arg := range arrArgs {
				argByte := convertToLittleEndian(params[i].Type, arg.([]byte))
				buffer = append(buffer, argByte...)
			}
		} else {
			buffer = convertToLittleEndian(params[i].Type, in.([]byte))
		}
		result = append(result, buffer)
	}

	return result, nil
}
