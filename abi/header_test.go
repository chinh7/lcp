package abi

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEncodeHeaderFromFile(t *testing.T) {
	encoded, err := EncodeHeaderToBytes("../test/testdata/liquid-token-abi.json")
	if err != nil {
		t.Errorf("error: %s", err)
	}
	result := []byte{248, 157, 1, 248, 86, 208, 132, 105, 110, 105, 116, 202, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3, 218, 136, 116, 114, 97, 110, 115, 102, 101, 114, 208, 197, 130, 116, 111, 128, 10, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3, 208, 132, 109, 105, 110, 116, 202, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3, 216, 139, 103, 101, 116, 95, 98, 97, 108, 97, 110, 99, 101, 203, 202, 135, 97, 100, 100, 114, 101, 115, 115, 128, 10, 248, 66, 234, 136, 84, 114, 97, 110, 115, 102, 101, 114, 224, 199, 132, 102, 114, 111, 109, 128, 10, 197, 130, 116, 111, 128, 10, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3, 199, 132, 109, 101, 109, 111, 128, 3, 214, 132, 77, 105, 110, 116, 208, 197, 130, 116, 111, 128, 10, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3}
	if !bytes.Equal(encoded, result) {
		t.Errorf("Encoding is incorrect,\nexpected:\t%v\nreality:\t%v.", result, encoded)
	}
}

func TestDecodeHeader(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	bytes := []byte{248, 157, 1, 248, 86, 208, 132, 105, 110, 105, 116, 202, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3, 218, 136, 116, 114, 97, 110, 115, 102, 101, 114, 208, 197, 130, 116, 111, 128, 10, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3, 208, 132, 109, 105, 110, 116, 202, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3, 216, 139, 103, 101, 116, 95, 98, 97, 108, 97, 110, 99, 101, 203, 202, 135, 97, 100, 100, 114, 101, 115, 115, 128, 10, 248, 66, 234, 136, 84, 114, 97, 110, 115, 102, 101, 114, 224, 199, 132, 102, 114, 111, 109, 128, 10, 197, 130, 116, 111, 128, 10, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3, 199, 132, 109, 101, 109, 111, 128, 3, 214, 132, 77, 105, 110, 116, 208, 197, 130, 116, 111, 128, 10, 201, 134, 97, 109, 111, 117, 110, 116, 128, 3}
	decoded, err := DecodeHeader(bytes)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(decoded, h, cmpopts.IgnoreUnexported(Event{}, Function{})); diff != "" {
		t.Errorf("Decoding of %v is incorrect, expected: %v, got: %v, diff: %v", bytes, h, decoded, diff)
	}
}

func TestGetEvent(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	event, err := h.GetEvent("Transfer")
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{}, Function{})
	if diff := cmp.Diff(event, h.Events[crypto.GetMethodID("Transfer")], opts); diff != "" {
		t.Errorf("GetEvent of %v is incorrect, expected: %v, got: %v, diff: %v", h, h.Events[crypto.GetMethodID("Transfer")], event, diff)
	}

	notFoundEvent, err := h.GetEvent("nil")
	if err == nil {
		t.Error("expecting error is nil for getting not found event")
	}
	if notFoundEvent != nil || err.Error() != "event nil not found" {
		t.Errorf("Error of GetEvent of %v is incorrect, expected: %v, got: %v", h, "event nil not found", err.Error())
	}
}

func TestGetFunction(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	event, err := h.GetFunction("transfer")
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{}, Function{})
	if diff := cmp.Diff(event, h.Functions[crypto.GetMethodID("transfer")], opts); diff != "" {
		t.Errorf("GetFunction of %v is incorrect, expected: %v, got: %v, diff: %v", h, h.Functions[crypto.GetMethodID("transfer")], event, diff)
	}

	notFoundFunction, err := h.GetFunction("nil")
	if err == nil {
		t.Error("expecting error is nil for getting not found function")
	}
	if notFoundFunction != nil || err.Error() != "function nil not found" {
		t.Errorf("Error of GetFunction of %v is incorrect, expected: %v, got: %v", h, "function nil not found", err.Error())
	}
}

func TestGetFunctionByMethodID(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/liquid-token-abi.json")
	function, err := h.GetFunctionByMethodID(crypto.GetMethodID("transfer"))
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{}, Function{})
	if diff := cmp.Diff(function, h.Functions[crypto.GetMethodID("transfer")], opts); diff != "" {
		t.Errorf("GetFunction of %v is incorrect, expected: %v, got: %v, diff: %v", h, h.Functions[crypto.GetMethodID("transfer")], function, diff)
	}

	notFoundFunction, err := h.GetFunctionByMethodID(crypto.MethodID{})
	if err == nil {
		t.Error("expecting error is nil for getting not found function")
	}
	expectedErr := fmt.Sprintf("function with methodID %v not found", crypto.MethodID{})
	if notFoundFunction != nil || err.Error() != expectedErr {
		t.Errorf("Error of GetFunction of %v is incorrect, expected: %v, got: %v", h, expectedErr, err.Error())
	}
}
