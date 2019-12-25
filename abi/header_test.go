package abi

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEncodeHeaderFromFile(t *testing.T) {
	encoded, err := EncodeHeaderToBytes("../test/testdata/token-abi.json")
	if err != nil {
		t.Errorf("error: %s", err)
	}
	result := []byte{248, 119, 1, 244, 210, 139, 103, 101, 116, 95, 98, 97, 108, 97, 110, 99, 101, 197, 196, 128, 128, 10, 128, 203, 132, 109, 105, 110, 116, 197, 196, 128, 128, 3, 128, 212, 136, 116, 114, 97, 110, 115, 102, 101, 114, 202, 196, 128, 128, 10, 128, 196, 128, 128, 3, 128, 248, 63, 216, 132, 77, 105, 110, 116, 210, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128, 229, 136, 84, 114, 97, 110, 115, 102, 101, 114, 219, 200, 132, 102, 114, 111, 109, 128, 10, 128, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128}
	if !bytes.Equal(encoded, result) {
		t.Errorf("Encoding is incorrect,\nexpected:\t%v\nreality:\t%v.", result, encoded)
	}
}

func TestDecodeHeader(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/token-abi.json")
	bytes := []byte{248, 119, 1, 244, 210, 139, 103, 101, 116, 95, 98, 97, 108, 97, 110, 99, 101, 197, 196, 128, 128, 10, 128, 203, 132, 109, 105, 110, 116, 197, 196, 128, 128, 3, 128, 212, 136, 116, 114, 97, 110, 115, 102, 101, 114, 202, 196, 128, 128, 10, 128, 196, 128, 128, 3, 128, 248, 63, 216, 132, 77, 105, 110, 116, 210, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128, 229, 136, 84, 114, 97, 110, 115, 102, 101, 114, 219, 200, 132, 102, 114, 111, 109, 128, 10, 128, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128}
	decoded, err := DecodeHeader(bytes)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(decoded, h, cmpopts.IgnoreUnexported(Event{})); diff != "" {
		t.Errorf("Decoding of %v is incorrect, expected: %v, got: %v, diff: %v", bytes, h, decoded, diff)
	}
}

func TestGetEvent(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/token-abi.json")
	event, err := h.GetEvent("Transfer")
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{})
	if diff := cmp.Diff(event, h.Events["Transfer"], opts); diff != "" {
		t.Errorf("GetEvent of %v is incorrect, expected: %v, got: %v, diff: %v", h, h.Events["Transfer"], event, diff)
	}

	notFoundEvent, err := h.GetEvent("nil")
	if err == nil {
		t.Error("expecting error is nil for getting not found event")
	}
	if notFoundEvent != nil || err.Error() != "event nil not found" {
		t.Errorf("Error of GetEvent of %v is incorrect, expected: %v, got: %v", h, "event nil not found", err.Error())
	}
}

func TestGetEventByIndex(t *testing.T) {
	bytes := []byte{248, 119, 1, 244, 210, 139, 103, 101, 116, 95, 98, 97, 108, 97, 110, 99, 101, 197, 196, 128, 128, 10, 128, 203, 132, 109, 105, 110, 116, 197, 196, 128, 128, 3, 128, 212, 136, 116, 114, 97, 110, 115, 102, 101, 114, 202, 196, 128, 128, 10, 128, 196, 128, 128, 3, 128, 248, 63, 216, 132, 77, 105, 110, 116, 210, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128, 229, 136, 84, 114, 97, 110, 115, 102, 101, 114, 219, 200, 132, 102, 114, 111, 109, 128, 10, 128, 198, 130, 116, 111, 128, 10, 128, 202, 134, 97, 109, 111, 117, 110, 116, 128, 3, 128}
	h, err := DecodeHeader(bytes)
	if err != nil {
		t.Error(err)
	}
	event, err := h.GetEventByIndex(1)
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{})
	if diff := cmp.Diff(event, h.Events["Transfer"], opts); diff != "" {
		t.Errorf("GetEventByIndex of %v is incorrect, expected: %v, got: %v, diff: %v", h, h.Events["Transfer"], event, diff)
	}

	notFoundEvent, err := h.GetEventByIndex(100)
	if err == nil {
		t.Error("expecting error is nil for getting not found event")
	}
	if notFoundEvent != nil || err.Error() != "Event not found" {
		t.Errorf("Error of GetEvent of %v is incorrect, expected: %v, got: %v", h, "Event not found", err.Error())
	}
}

func TestGetFunction(t *testing.T) {
	h, _ := LoadHeaderFromFile("../test/testdata/token-abi.json")
	event, err := h.GetFunction("transfer")
	if err != nil {
		t.Error(err)
	}
	opts := cmpopts.IgnoreUnexported(Event{})
	if diff := cmp.Diff(event, h.Functions["transfer"], opts); diff != "" {
		t.Errorf("GetFunction of %v is incorrect, expected: %v, got: %v, diff: %v", h, h.Functions["transfer"], event, diff)
	}

	notFoundFunction, err := h.GetFunction("nil")
	if err == nil {
		t.Error("expecting error is nil for getting not found function")
	}
	if notFoundFunction != nil || err.Error() != "function nil not found" {
		t.Errorf("Error of GetFunction of %v is incorrect, expected: %v, got: %v", h, "function nil not found", err.Error())
	}
}

func TestEventGetIndex(t *testing.T) {
	index := uint32(5)
	event := Event{index: index}
	result := event.GetIndex()
	if result != index {
		t.Errorf("Error of TestEventGetIndex expected: %v, got: %v", event.index, result)
	}
}

func TestEventGetIndexByte(t *testing.T) {
	index := uint32(5)
	bytes := []byte{5, 0, 0, 0}
	event := Event{index: index}
	result := event.GetIndexByte()
	if diff := cmp.Diff(result, bytes); diff != "" {
		t.Errorf("GetIndexByte of %v is incorrect, expected: %v, got: %v, diff: %v", event, bytes, result, diff)
	}
}
