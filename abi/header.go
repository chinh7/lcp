package abi

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Parameter is model for function signature
type Parameter struct {
	IsArray bool          `json:"is_array"`
	Type    PrimitiveType `json:"type"`
}

// Function is model for function signature
type Function struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"`
}

// Header is model for function signature
type Header struct {
	Version   uint16              `json:"version"`
	Functions map[string]Function `json:"functions"`
}

const (
	// HeaderVersionByteLength is number of bytes preservered for version number
	HeaderVersionByteLength = 2
	// HeaderFunctionCountByteLength is number of bytes preservered for number of functions in header
	HeaderFunctionCountByteLength = 1
	// FunctionNameByteLength is number of bytes preservered for function name
	FunctionNameByteLength = 64
	// FunctionParameterCountByteLength is number of bytes preservered for number of Parameters in a function
	FunctionParameterCountByteLength = 1
	// ParameterByteLength is number of bytes preservered for a parameter
	ParameterByteLength = 2
)

// GetFunction returns function of a header from the func name
func (h Header) GetFunction(funcName string) (Function, error) {
	if f, found := h.Functions[funcName]; found {
		return f, nil
	}
	return Function{}, fmt.Errorf("function %s not found", funcName)
}

func decodeParam(b []byte) (Parameter, error) {
	var param Parameter
	switch b[0] {
	case 0:
		param.IsArray = false
	case 1:
		param.IsArray = true
	default:
		return Parameter{}, fmt.Errorf("not valid IsArray byte for parameter decoding: %v", b[0])
	}
	param.Type = PrimitiveType(b[1])
	return param, nil
}

func decodeFunction(b []byte) (Function, error) {
	var fc Function
	var offset int
	fc.Name = string(bytes.Trim(b[offset:offset+FunctionNameByteLength], "\x00"))
	offset += FunctionNameByteLength
	paramsLength := int(b[offset])
	offset++
	for index := 0; index < paramsLength; index++ {
		param, err := decodeParam(b[offset : offset+ParameterByteLength])
		if err != nil {
			return Function{}, err
		}
		offset += ParameterByteLength
		fc.Parameters = append(fc.Parameters, param)
	}
	return fc, nil
}

// DecodeHeader decode byte array of header into header
func DecodeHeader(b []byte) (Header, int, error) {
	var header Header
	header.Functions = make(map[string]Function)
	var offset int

	header.Version = binary.LittleEndian.Uint16(b[offset : offset+HeaderVersionByteLength])
	offset += HeaderVersionByteLength

	funcsLength := int(b[offset : offset+1][0])
	offset++
	for index := 0; index < funcsLength; index++ {
		var funcByteLength int
		funcByteLength += FunctionNameByteLength
		paramsLength := int(b[offset+funcByteLength])
		funcByteLength += FunctionParameterCountByteLength + ParameterByteLength*paramsLength
		fc, err := decodeFunction(b[offset : offset+funcByteLength])
		if err != nil {
			return Header{}, 0, err
		}
		offset += funcByteLength
		header.Functions[fc.Name] = fc
	}
	return header, offset, nil
}

// Encode encode a Parameter struct into byte array
// encoding schema: is array(1 byte)|type (1 byte)
func (param Parameter) Encode() ([]byte, error) {
	encodedBytes := make([]byte, ParameterByteLength)
	if param.IsArray {
		encodedBytes[0] = 1
	}
	encodedBytes[1] = byte(param.Type)
	return encodedBytes, nil
}

// Encode encode a Function struct into byte array
// encoding schema: name of function(64 bytes)|number of parameters(1 byte)|parameter1(2 bytes)|parameter2(2 bytes)|...
func (f Function) Encode() ([]byte, error) {
	var encodedBytes []byte
	if len(f.Name) > FunctionNameByteLength {
		return []byte{0}, fmt.Errorf("function name too long, got: %v, expected less or equal: %v", len(f.Name), FunctionNameByteLength)
	}
	nameBytes := make([]byte, FunctionNameByteLength)
	copy(nameBytes[:], f.Name)
	encodedBytes = append(encodedBytes, nameBytes...)

	encodedBytes = append(encodedBytes, byte(len(f.Parameters)))
	for _, param := range f.Parameters {
		paramBytes, err := param.Encode()
		if err != nil {
			return []byte{0}, err
		}
		encodedBytes = append(encodedBytes, paramBytes...)
	}
	return encodedBytes, nil
}

// Encode encode a header struct into byte array
// encoding schema: version(2 bytes)|number of functions(1 byte)|function1|function2|...
func (h Header) Encode() ([]byte, error) {
	var encodedBytes []byte

	versionBytes := make([]byte, HeaderVersionByteLength)
	binary.LittleEndian.PutUint16(versionBytes, h.Version)
	encodedBytes = append(encodedBytes, versionBytes...)

	encodedBytes = append(encodedBytes, byte(len(h.Functions)))
	for _, f := range h.Functions {
		functionBytes, err := f.Encode()
		if err != nil {
			return []byte{0}, err
		}
		encodedBytes = append(encodedBytes, functionBytes...)
	}

	return encodedBytes, nil
}
