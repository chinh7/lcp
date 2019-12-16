package abi

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/ethereum/go-ethereum/rlp"
)

// Event is emitting from engine
type Event struct {
	Name       string       `json:"name"`
	Parameters []*Parameter `json:"parameters"`
	index      uint32
}

// Parameter is model for function signature
type Parameter struct {
	Name    string        `json:"name"`
	IsArray bool          `json:"is_array"`
	Type    PrimitiveType `json:"type"`
	Size    uint          `json:"size"`
}

// Function is model for function signature
type Function struct {
	Name       string       `json:"name"`
	Parameters []*Parameter `json:"parameters"`
}

// Header is model for function signature
type Header struct {
	Version   uint16               `json:"version"`
	Functions map[string]*Function `json:"functions"`
	Events    map[string]*Event    `json:"events"`
}

// GetEvent return the event
func (h Header) GetEvent(name string) (*Event, error) {
	if event, ok := h.Events[name]; ok {
		return event, nil
	}
	return nil, fmt.Errorf("event %s not found", name)
}

// GetEventByIndex use index to retrieve event
func (h Header) GetEventByIndex(index uint32) (*Event, error) {
	for _, event := range h.Events {
		if event.GetIndex() == index {
			return event, nil
		}
	}
	return nil, errors.New("Event not found")
}

// GetFunction returns function of a header from the func name
func (h Header) GetFunction(funcName string) (*Function, error) {
	if f, found := h.Functions[funcName]; found {
		return f, nil
	}
	return nil, fmt.Errorf("function %s not found", funcName)
}

// DecodeHeader decode byte array of header into header
func DecodeHeader(b []byte) (*Header, error) {
	var header struct {
		Version   uint16
		Functions []*Function
		Events    []*Event
	}
	if err := rlp.DecodeBytes(b, &header); err != nil {
		return nil, err
	}

	functions := make(map[string]*Function)
	for _, function := range header.Functions {
		functions[function.Name] = function
	}

	events := make(map[string]*Event)
	for index, event := range header.Events {
		event.index = uint32(index)
		events[event.Name] = event
	}

	return &Header{header.Version, functions, events}, nil
}

// Encode encode a header struct into byte array
// encoding schema: version(2 bytes)|number of functions(1 byte)|function1|function2|...
func (h *Header) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(h)
}

func (h *Header) getEvents() []*Event {
	events := []*Event{}
	var keys []string
	for key := range h.Events {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		events = append(events, h.Events[key])
	}
	return events
}

func (h *Header) getFunctions() []*Function {
	functions := []*Function{}
	var keys []string
	for key := range h.Functions {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		functions = append(functions, h.Functions[key])
	}
	return functions
}

// EncodeRLP encodes a header to RLP format
func (h *Header) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, struct {
		Version   uint16
		Functions []*Function
		Events    []*Event
	}{
		Version:   h.Version,
		Functions: h.getFunctions(),
		Events:    h.getEvents(),
	})
}

// GetIndex return index of event
func (e *Event) GetIndex() uint32 {
	return e.index
}

// GetIndexByte return []byte representation of event
func (e *Event) GetIndexByte() []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, e.index)
	return b
}
