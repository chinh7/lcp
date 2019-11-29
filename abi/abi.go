package abi

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/QuoineFinancial/vertex/crypto"
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

const (
	// ArrayLengthByte is number of bytes preservered to indicate an array's length
	ArrayLengthByte = 2
)

// PrimitiveArg is model for a primitive type
type PrimitiveArg struct {
	Type  PrimitiveType
	Value interface{}
}

// ArrayArg is model for dynamic array of primitive type
type ArrayArg []PrimitiveArg

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
func (p PrimitiveType) GetMemorySize() int {
	switch p {
	case Address:
		return crypto.AddressLength
	case Uint8, Int8:
		return 1
	case Uint16, Int16:
		return 2
	case Uint32, Int32, Float32:
		return 4
	case Uint64, Int64, Float64:
		return 8
	default:
		panic("primitive type not found")
	}
}

// newPrimitiveArg parse type string and value into PrimitiveArg
func newPrimitiveArg(t PrimitiveType, value interface{}) PrimitiveArg {
	var res PrimitiveArg
	res.Type = t
	res.Value = value
	return res
}

func ArrayFromUint64(t PrimitiveType, values []uint64) interface{} {
	switch t {
	case Uint8:
		var v []uint8
		for _, value := range values {
			v = append(v, uint8(value))
		}
		return v
	case Uint16:
		var v []uint16
		for _, value := range values {
			v = append(v, uint16(value))
		}
		return v
	case Uint32:
		var v []uint32
		for _, value := range values {
			v = append(v, uint32(value))
		}
		return v
	case Uint64:
		return values
	case Int8:
		var v []int8
		for _, value := range values {
			v = append(v, int8(value))
		}
		return v
	case Int16:
		var v []int16
		for _, value := range values {
			v = append(v, int16(value))
		}
		return v
	case Int32:
		var v []int32
		for _, value := range values {
			v = append(v, int32(value))
		}
		return v
	case Int64:
		var v []int64
		for _, value := range values {
			v = append(v, int64(value))
		}
		return v
	case Float32:
		var v []uint8
		for _, value := range values {
			v = append(v, uint8(value))
		}
		return v
	case Float64:
		var v []float64
		for _, value := range values {
			v = append(v, math.Float64frombits(value))
		}
		return v
	}
	return nil
}

func PrimitiveFromUint64(t PrimitiveType, value uint64) interface{} {
	var v interface{}
	switch t {
	case Uint8:
		v = uint8(value)
	case Uint16:
		v = uint16(value)
	case Uint32:
		v = uint32(value)
	case Uint64:
		v = value
	case Int8:
		v = int8(value)
	case Int16:
		v = int16(value)
	case Int32:
		v = int32(value)
	case Int64:
		v = int64(value)
	case Float32:
		v = math.Float32frombits(uint32(value))
	case Float64:
		v = math.Float64frombits(value)
	}
	return v
}

// parseArrayArg parse type and values into ArrayArg
func parseArrayArg(t PrimitiveType, value interface{}) (ArrayArg, error) {
	var elements ArrayArg
	switch t {
	case Address:
		parsed, ok := value.([]crypto.Address)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Uint8:
		parsed, ok := value.([]uint8)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Uint16:
		parsed, ok := value.([]uint16)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Uint32:
		parsed, ok := value.([]uint32)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Uint64:
		parsed, ok := value.([]uint64)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Int8:
		parsed, ok := value.([]int8)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Int16:
		parsed, ok := value.([]int16)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Int32:
		parsed, ok := value.([]int32)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Int64:
		parsed, ok := value.([]int64)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Float32:
		parsed, ok := value.([]float32)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	case Float64:
		parsed, ok := value.([]float64)
		if !ok {
			return ArrayArg{}, fmt.Errorf("unable to convert array element into %s", t)
		}
		for _, p := range parsed {
			e := newPrimitiveArg(t, p)
			elements = append(elements, e)
		}
	default:
		return ArrayArg{}, fmt.Errorf("not supported type: %s", t)
	}
	return elements, nil
}

// Encode is common interface method for encoding a PrimitiveArg into byte array
func (p PrimitiveArg) encode() ([]byte, error) {
	memorySize := p.Type.GetMemorySize()
	buf := make([]byte, memorySize)
	switch p.Type {
	case Address:
		address := p.Value.(crypto.Address)
		copy(buf, address[:])
	case Uint8:
		buf[0] = byte(p.Value.(uint8))
	case Uint16:
		binary.LittleEndian.PutUint16(buf, p.Value.(uint16))
	case Uint32:
		binary.LittleEndian.PutUint32(buf, p.Value.(uint32))
	case Uint64:
		binary.LittleEndian.PutUint64(buf, p.Value.(uint64))
	case Int8:
		buf[0] = byte(p.Value.(int8))
	case Int16:
		binary.LittleEndian.PutUint16(buf, uint16(p.Value.(int16)))
	case Int32:
		binary.LittleEndian.PutUint32(buf, uint32(p.Value.(int32)))
	case Int64:
		binary.LittleEndian.PutUint64(buf, uint64(p.Value.(int64)))
	case Float32:
		binary.LittleEndian.PutUint32(buf, math.Float32bits(p.Value.(float32)))
	case Float64:
		binary.LittleEndian.PutUint64(buf, math.Float64bits(p.Value.(float64)))
	default:
		return nil, fmt.Errorf("not supported type: %s", p.Type)
	}
	return buf, nil
}

// Encode is common interface method for encoding a ArrayArg into byte array
func (arrayArg ArrayArg) encode() ([]byte, error) {
	result := []byte{}
	encodedLength := make([]byte, ArrayLengthByte)
	binary.LittleEndian.PutUint16(encodedLength, uint16(len(arrayArg)))
	result = append(result, encodedLength...)

	for _, e := range arrayArg {
		encodedBytes, err := e.encode()
		if err != nil {
			return []byte{0}, err
		}
		result = append(result, encodedBytes...)
	}

	return result, nil
}

// singleDecode decode byte array into interface based on its type
func singleDecode(t PrimitiveType, buf []byte) (interface{}, error) {
	var result interface{}
	switch t {
	case Address:
		var address crypto.Address
		copy(address[:], buf)
		result = address
	case Uint8:
		result = uint8(buf[0])
	case Uint16:
		result = binary.LittleEndian.Uint16(buf)
	case Uint32:
		result = binary.LittleEndian.Uint32(buf)
	case Uint64:
		result = binary.LittleEndian.Uint64(buf)
	case Int8:
		result = int8(buf[0])
	case Int16:
		result = int16(binary.LittleEndian.Uint16(buf))
	case Int32:
		result = int32(binary.LittleEndian.Uint32(buf))
	case Int64:
		result = int64(binary.LittleEndian.Uint64(buf))
	case Float32:
		bits := binary.LittleEndian.Uint32(buf)
		result = math.Float32frombits(bits)
	case Float64:
		bits := binary.LittleEndian.Uint64(buf)
		result = math.Float64frombits(bits)
	default:
		return nil, fmt.Errorf("not supported type: %s", t)
	}
	return result, nil
}

// arrayDecode decode encoded dynamic array based on type and length
func arrayDecode(t PrimitiveType, length int, buf []byte) (interface{}, error) {
	results := []interface{}{}
	var offset int
	for index := 0; index < length; index++ {
		memorySize := t.GetMemorySize()
		result, err := singleDecode(t, buf[offset:offset+memorySize])
		if err != nil {
			return []interface{}{}, err
		}
		offset += memorySize
		results = append(results, result)
	}
	return results, nil
}

// Encode return []byte from an inputted params and values pair
func Encode(params []*Parameter, values []interface{}) ([]byte, error) {
	if len(params) != len(values) {
		return []byte{0}, fmt.Errorf("Parameter count mismatch, expecting: %d, got: %d", len(params), len(values))
	}
	result := []byte{}

	for index, param := range params {
		var encodedBytes []byte
		if param.IsArray {
			arrayArg, err := parseArrayArg(param.Type, values[index])
			if err != nil {
				return nil, err
			}
			bytes, err := arrayArg.encode()
			if err != nil {
				return nil, err
			}
			encodedBytes = append(encodedBytes, bytes...)
		} else {
			arg := newPrimitiveArg(param.Type, values[index])
			bytes, err := arg.encode()
			if err != nil {
				return nil, err
			}
			encodedBytes = append(encodedBytes, bytes...)
		}
		result = append(result, encodedBytes...)
	}
	return result, nil
}

// Decode return []interface from an inputted params and []byte
func Decode(params []*Parameter, bytes []byte) ([]interface{}, error) {
	var results []interface{}
	var offset int
	for _, param := range params {
		if param.IsArray {
			length := int(binary.LittleEndian.Uint16(bytes[offset : offset+ArrayLengthByte]))
			offset += ArrayLengthByte

			elementSize := param.Type.GetMemorySize()
			memorySize := elementSize * length
			result, err := arrayDecode(param.Type, length, bytes[offset:offset+memorySize])
			if err != nil {
				return []interface{}{}, err
			}
			offset += memorySize
			results = append(results, result)
		} else {
			memorySize := param.Type.GetMemorySize()
			result, err := singleDecode(param.Type, bytes[offset:offset+memorySize])
			if err != nil {
				return []interface{}{}, err
			}
			offset += memorySize
			results = append(results, result)
		}
	}
	return results, nil
}

// DecodeToBytes returns uint64 array compatible with VM
func DecodeToBytes(params []*Parameter, bytes []byte) ([][]byte, error) {
	var decoded [][]byte
	var offset int
	for _, param := range params {
		var arg []byte
		var memorySize int
		if param.IsArray {
			length := int(binary.LittleEndian.Uint16(bytes[offset : offset+ArrayLengthByte]))
			offset += ArrayLengthByte
			elementSize := param.Type.GetMemorySize()
			memorySize = elementSize * length
		} else {
			memorySize = param.Type.GetMemorySize()
		}
		arg = bytes[offset : offset+memorySize]
		offset += memorySize
		decoded = append(decoded, arg)
	}
	return decoded, nil
}

// EncodeFromBytes encodes arguments in byte format - an inverse of DecodeToBytes
func EncodeFromBytes(params []*Parameter, bytes [][]byte) ([]byte, error) {
	var encoded []byte
	for i, param := range params {
		var memorySize int
		if param.IsArray {
			elementSize := param.Type.GetMemorySize()
			length := len(bytes[i])
			if length%elementSize != 0 {
				return nil, errors.New("misaligned array byte size")
			}
			length = length / elementSize
			lengthBytes := make([]byte, ArrayLengthByte)
			binary.LittleEndian.PutUint16(lengthBytes, uint16(length))
			encoded = append(encoded, lengthBytes...)

		} else {
			memorySize = param.Type.GetMemorySize()
			if memorySize != len(bytes[i]) {
				return nil, errors.New("mismatched primitive byte size")
			}
		}
		encoded = append(encoded, bytes[i]...)
	}
	return encoded, nil
}
