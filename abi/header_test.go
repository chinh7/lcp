package abi

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEncodeHeaderFromFile(t *testing.T) {
	encoded, err := EncodeHeaderToBytes("../test/fixtures/header-event.json")
	if err != nil {
		t.Errorf("error: %s", err)
	}
	result := []byte{248, 122, 131, 49, 46, 48, 244, 210, 139, 103, 101, 116, 95, 98, 97, 108, 97, 110, 99, 101, 197, 196, 128, 128, 10, 128, 203, 132, 109, 105, 110, 116, 197, 196, 128, 128, 3, 128, 212, 136, 116, 114, 97, 110, 115, 102, 101, 114, 202, 196, 128, 128, 10, 128, 196, 128, 128, 3, 128, 248, 63, 216, 132, 77, 105, 110, 116, 210, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128, 229, 136, 84, 114, 97, 110, 115, 102, 101, 114, 219, 200, 132, 102, 114, 111, 109, 128, 10, 128, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128}
	if !bytes.Equal(encoded, result) {
		t.Errorf("Encoding is incorrect,\nexpected:\t%v\nreality:\t%v.", result, encoded)
	}
}

func TestDecodeHeader(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/fixtures/header-event.json")
	bytes := []byte{248, 122, 131, 49, 46, 48, 244, 210, 139, 103, 101, 116, 95, 98, 97, 108, 97, 110, 99, 101, 197, 196, 128, 128, 10, 128, 203, 132, 109, 105, 110, 116, 197, 196, 128, 128, 3, 128, 212, 136, 116, 114, 97, 110, 115, 102, 101, 114, 202, 196, 128, 128, 10, 128, 196, 128, 128, 3, 128, 248, 63, 216, 132, 77, 105, 110, 116, 210, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128, 229, 136, 84, 114, 97, 110, 115, 102, 101, 114, 219, 200, 132, 102, 114, 111, 109, 128, 10, 128, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128}
	decoded, err := DecodeHeader(bytes)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(decoded, h); diff != "" {
		t.Errorf("Decoding of %v is incorrect, expected: %v, got: %v, diff: %v", bytes, h, decoded, diff)
	}
}
