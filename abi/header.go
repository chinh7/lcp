package abi

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/rlp"
)

// Event is emitting from engine
type Event struct {
	Name       string
	Parameters []*Parameter
}

// Parameter is model for function signature
type Parameter struct {
	Name    string
	IsArray bool
	Type    PrimitiveType
	Size    uint
}

// Function is model for function signature
type Function struct {
	Name       string
	Parameters []*Parameter
}

// Header is model for function signature
type Header struct {
	Version   string
	Functions map[string]*Function
	Events    map[string]*Event
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

// GetEvent return the event
func (h Header) GetEvent(name string) (*Event, error) {
	if event, ok := h.Events[name]; ok {
		return event, nil
	}
	return nil, fmt.Errorf("event %s not found", name)
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
		Version   string
		Functions []*Function
		Events    []*Event
	}
	rlp.DecodeBytes(b, &header)

	functions := make(map[string]*Function)
	for _, function := range header.Functions {
		functions[function.Name] = function
	}

	events := make(map[string]*Event)
	for _, event := range header.Events {
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
	for _, event := range h.Events {
		events = append(events, event)
	}
	return events
}

func (h *Header) getFunctions() []*Function {
	functions := []*Function{}
	for _, function := range h.Functions {
		functions = append(functions, function)
	}
	return functions
}

// EncodeRLP encodes a header to RLP format
func (h *Header) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, struct {
		Version   string
		Functions []*Function
		Events    []*Event
	}{
		Version:   h.Version,
		Functions: h.getFunctions(),
		Events:    h.getEvents(),
	})
}
