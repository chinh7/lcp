package abi

import (
	"bytes"
	"strings"
	"testing"

	"github.com/QuoineFinancial/vertex/crypto"
	"github.com/google/go-cmp/cmp"
)

func parseParameterFromString(s string) (Parameter, error) {
	var p Parameter

	if s[len(s)-2:] == "[]" {
		p.IsArray = true
		t, err := parsePrimitiveTypeFromString(s[:strings.Index(s, "[")])
		if err != nil {
			return Parameter{}, err
		}
		p.Type = t
	} else {
		p.IsArray = false
		t, err := parsePrimitiveTypeFromString(s)
		if err != nil {
			return Parameter{}, err
		}
		p.Type = t
	}
	return p, nil
}

func TestEncode(t *testing.T) {
	address := crypto.AddressFromString("LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O")
	address2 := crypto.AddressFromString("LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O")
	addresses := []crypto.Address{address, address2}
	var parameters1 []*Parameter
	var parameters2 []*Parameter
	paramsString1 := []string{"address", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64"}
	paramsString2 := []string{"address[]", "uint8[]", "uint16[]", "uint32[]", "uint64[]", "int8[]", "int16[]", "int32[]", "int64[]", "float32[]", "float64[]"}

	for _, p := range paramsString1 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters1 = append(parameters1, &param)
	}
	for _, p := range paramsString2 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters2 = append(parameters2, &param)
	}

	testTables := []struct {
		types  []*Parameter
		values []interface{}
		result []byte
	}{
		{
			types:  parameters1,
			values: []interface{}{address, uint8(88), uint16(43221), uint32(3333324342), uint64(3213214325432656666), int8(88), int16(4321), int32(-34325), int64(-321452), float32(8321.38), float64(-4321452.1188)},
			result: []byte{88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 213, 168, 54, 126, 174, 198, 26, 35, 156, 150, 103, 161, 151, 44, 88, 225, 16, 235, 121, 255, 255, 84, 24, 251, 255, 255, 255, 255, 255, 133, 5, 2, 70, 81, 107, 154, 7, 43, 124, 80, 193},
		},
		{
			types:  parameters2,
			values: []interface{}{addresses, []uint8{uint8(88), uint8(255)}, []uint16{uint16(555), uint16(12333)}, []uint32{uint32(3333324342), uint32(3333324342), uint32(33324342)}, []uint64{uint64(3213214325432656666), uint64(32145467)}, []int8{int8(88), int8(-88)}, []int16{int16(333), int16(-542)}, []int32{int32(43298), int32(-321432)}, []int64{int64(-23425254), int64(10875498375)}, []float32{float32(-1341.233), float32(50492.235)}, []float64{float64(-132341.233), float64(50454392.235)}},
			result: []byte{2, 0, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 2, 0, 88, 255, 2, 0, 43, 2, 45, 48, 3, 0, 54, 126, 174, 198, 54, 126, 174, 198, 54, 125, 252, 1, 2, 0, 26, 35, 156, 150, 103, 161, 151, 44, 59, 128, 234, 1, 0, 0, 0, 0, 2, 0, 88, 168, 2, 0, 77, 1, 226, 253, 2, 0, 34, 169, 0, 0, 104, 24, 251, 255, 2, 0, 26, 143, 154, 254, 255, 255, 255, 255, 135, 239, 58, 136, 2, 0, 0, 0, 2, 0, 117, 167, 167, 196, 60, 60, 69, 71, 2, 0, 160, 26, 47, 221, 169, 39, 0, 193, 174, 71, 225, 193, 251, 14, 136, 65},
		},
	}

	for _, table := range testTables {
		encoded, err := Encode(table.types, table.values)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		if !bytes.Equal(encoded, table.result) {
			t.Errorf("Encoding of %v is incorrect, expected: %v, got: %v.", table.values, table.result, encoded)
		}
	}
}

func TestEncodeFromString(t *testing.T) {
	var parameters1 []*Parameter
	var parameters2 []*Parameter
	paramsString1 := []string{"address", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64"}
	paramsString2 := []string{"address[]", "uint8[]", "uint16[]", "uint32[]", "uint64[]", "int8[]", "int16[]", "int32[]", "int64[]", "float32[]", "float64[]"}

	for _, p := range paramsString1 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters1 = append(parameters1, &param)
	}
	for _, p := range paramsString2 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters2 = append(parameters2, &param)
	}

	testTables := []struct {
		types  []*Parameter
		values []string
		result []byte
	}{
		{
			types:  parameters1,
			values: []string{"LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O", "88", "43221", "3333324342", "3213214325432656666", "88", "4321", "-34325", "-321452", "8321.38", "-4321452.1188"},
			result: []byte{88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 213, 168, 54, 126, 174, 198, 26, 35, 156, 150, 103, 161, 151, 44, 88, 225, 16, 235, 121, 255, 255, 84, 24, 251, 255, 255, 255, 255, 255, 133, 5, 2, 70, 81, 107, 154, 7, 43, 124, 80, 193},
		},
		{
			types:  parameters2,
			values: []string{"[LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O, LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O]", "[88,255]", "[555,12333]", "[3333324342,3333324342,33324342]", "[3213214325432656666,32145467]", "[88,-88]", "[333,-542]", "[43298,-321432]", "[-23425254,10875498375]", "[-1341.233,50492.235]", "[-132341.233,50454392.235]"},
			result: []byte{2, 0, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 2, 0, 88, 255, 2, 0, 43, 2, 45, 48, 3, 0, 54, 126, 174, 198, 54, 126, 174, 198, 54, 125, 252, 1, 2, 0, 26, 35, 156, 150, 103, 161, 151, 44, 59, 128, 234, 1, 0, 0, 0, 0, 2, 0, 88, 168, 2, 0, 77, 1, 226, 253, 2, 0, 34, 169, 0, 0, 104, 24, 251, 255, 2, 0, 26, 143, 154, 254, 255, 255, 255, 255, 135, 239, 58, 136, 2, 0, 0, 0, 2, 0, 117, 167, 167, 196, 60, 60, 69, 71, 2, 0, 160, 26, 47, 221, 169, 39, 0, 193, 174, 71, 225, 193, 251, 14, 136, 65},
		},
	}

	for _, table := range testTables {
		encoded, err := EncodeFromString(table.types, table.values)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		if !bytes.Equal(encoded, table.result) {
			t.Errorf("Encoding of %v is incorrect, expected: %v, got: %v.", table.values, table.result, encoded)
		}
	}
}

func TestDecode(t *testing.T) {
	address := crypto.AddressFromString("LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O")
	address2 := crypto.AddressFromString("LCHILMXMODD5DMDMPKVSD5MUODDQMBRU5GZVLGXEFBPG36HV4CLSYM7O")
	addresses := []interface{}{address, address2}
	var parameters1 []*Parameter
	var parameters2 []*Parameter
	paramsString1 := []string{"address", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64"}
	paramsString2 := []string{"address[]", "uint8[]", "uint16[]", "uint32[]", "uint64[]", "int8[]", "int16[]", "int32[]", "int64[]", "float32[]", "float64[]"}

	for _, p := range paramsString1 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters1 = append(parameters1, &param)
	}
	for _, p := range paramsString2 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters2 = append(parameters2, &param)
	}

	testTables := []struct {
		types  []*Parameter
		values []byte
		result []interface{}
	}{
		{
			types:  parameters1,
			values: []byte{88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 213, 168, 54, 126, 174, 198, 26, 35, 156, 150, 103, 161, 151, 44, 88, 225, 16, 235, 121, 255, 255, 84, 24, 251, 255, 255, 255, 255, 255, 133, 5, 2, 70, 81, 107, 154, 7, 43, 124, 80, 193},
			result: []interface{}{address, uint8(88), uint16(43221), uint32(3333324342), uint64(3213214325432656666), int8(88), int16(4321), int32(-34325), int64(-321452), float32(8321.38), float64(-4321452.1188)},
		},
		{
			types:  parameters2,
			values: []byte{2, 0, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 2, 0, 88, 255, 2, 0, 43, 2, 45, 48, 3, 0, 54, 126, 174, 198, 54, 126, 174, 198, 54, 125, 252, 1, 2, 0, 26, 35, 156, 150, 103, 161, 151, 44, 59, 128, 234, 1, 0, 0, 0, 0, 2, 0, 88, 168, 2, 0, 77, 1, 226, 253, 2, 0, 34, 169, 0, 0, 104, 24, 251, 255, 2, 0, 26, 143, 154, 254, 255, 255, 255, 255, 135, 239, 58, 136, 2, 0, 0, 0, 2, 0, 117, 167, 167, 196, 60, 60, 69, 71, 2, 0, 160, 26, 47, 221, 169, 39, 0, 193, 174, 71, 225, 193, 251, 14, 136, 65},
			result: []interface{}{addresses, []interface{}{uint8(88), uint8(255)}, []interface{}{uint16(555), uint16(12333)}, []interface{}{uint32(3333324342), uint32(3333324342), uint32(33324342)}, []interface{}{uint64(3213214325432656666), uint64(32145467)}, []interface{}{int8(88), int8(-88)}, []interface{}{int16(333), int16(-542)}, []interface{}{int32(43298), int32(-321432)}, []interface{}{int64(-23425254), int64(10875498375)}, []interface{}{float32(-1341.233), float32(50492.235)}, []interface{}{float64(-132341.233), float64(50454392.235)}},
		},
	}

	for _, table := range testTables {
		decoded, err := Decode(table.types, table.values)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		if diff := cmp.Diff(decoded, table.result); diff != "" {
			t.Errorf("Decoding of %v is incorrect, expected: %v, got: %v, diff: %v", table.values, table.result, decoded, diff)
		}
	}
}

func TestDecodeToBytes(t *testing.T) {
	var parameters1 []*Parameter
	var parameters2 []*Parameter
	paramsString1 := []string{"address", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64"}
	paramsString2 := []string{"address[]", "uint8[]", "uint16[]", "uint32[]", "uint64[]", "int8[]", "int16[]", "int32[]", "int64[]", "float32[]", "float64[]"}

	for _, p := range paramsString1 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters1 = append(parameters1, &param)
	}
	for _, p := range paramsString2 {
		param, err := parseParameterFromString(p)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		parameters2 = append(parameters2, &param)
	}

	testTables := []struct {
		types  []*Parameter
		values []byte
		result [][]byte
	}{
		{
			types:  parameters1,
			values: []byte{88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 213, 168, 54, 126, 174, 198, 26, 35, 156, 150, 103, 161, 151, 44, 88, 225, 16, 235, 121, 255, 255, 84, 24, 251, 255, 255, 255, 255, 255, 133, 5, 2, 70, 81, 107, 154, 7, 43, 124, 80, 193},
			result: [][]byte{[]byte{88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238}, []byte{88}, []byte{213, 168}, []byte{54, 126, 174, 198}, []byte{26, 35, 156, 150, 103, 161, 151, 44}, []byte{88}, []byte{225, 16}, []byte{235, 121, 255, 255}, []byte{84, 24, 251, 255, 255, 255, 255, 255}, []byte{133, 5, 2, 70}, []byte{81, 107, 154, 7, 43, 124, 80, 193}},
		},
		{
			types:  parameters2,
			values: []byte{2, 0, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 2, 0, 88, 255, 2, 0, 43, 2, 45, 48, 3, 0, 54, 126, 174, 198, 54, 126, 174, 198, 54, 125, 252, 1, 2, 0, 26, 35, 156, 150, 103, 161, 151, 44, 59, 128, 234, 1, 0, 0, 0, 0, 2, 0, 88, 168, 2, 0, 77, 1, 226, 253, 2, 0, 34, 169, 0, 0, 104, 24, 251, 255, 2, 0, 26, 143, 154, 254, 255, 255, 255, 255, 135, 239, 58, 136, 2, 0, 0, 0, 2, 0, 117, 167, 167, 196, 60, 60, 69, 71, 2, 0, 160, 26, 47, 221, 169, 39, 0, 193, 174, 71, 225, 193, 251, 14, 136, 65},
			result: [][]byte{[]byte{88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238, 88, 142, 133, 178, 236, 112, 199, 209, 176, 108, 122, 171, 33, 245, 148, 112, 199, 6, 6, 52, 233, 179, 85, 154, 228, 40, 94, 109, 248, 245, 224, 151, 44, 51, 238}, []byte{88, 255}, []byte{43, 2, 45, 48}, []byte{54, 126, 174, 198, 54, 126, 174, 198, 54, 125, 252, 1}, []byte{26, 35, 156, 150, 103, 161, 151, 44, 59, 128, 234, 1, 0, 0, 0, 0}, []byte{88, 168}, []byte{77, 1, 226, 253}, []byte{34, 169, 0, 0, 104, 24, 251, 255}, []byte{26, 143, 154, 254, 255, 255, 255, 255, 135, 239, 58, 136, 2, 0, 0, 0}, []byte{117, 167, 167, 196, 60, 60, 69, 71}, []byte{160, 26, 47, 221, 169, 39, 0, 193, 174, 71, 225, 193, 251, 14, 136, 65}},
		},
	}

	for _, table := range testTables {
		decoded, err := DecodeToBytes(table.types, table.values)
		if err != nil {
			t.Errorf("error: %s", err)
		}
		if diff := cmp.Diff(decoded, table.result); diff != "" {
			t.Errorf("Decoding of %v is incorrect, expected: %v, got: %v, diff: %v", table.values, table.result, decoded, diff)
		}
	}
}
