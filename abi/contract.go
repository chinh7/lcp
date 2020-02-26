package abi

import (
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/ethereum/go-ethereum/rlp"
)

// Contract contains header and wasm code
type Contract struct {
	Header *Header
	Code   []byte
}

// DecodeContract decode []byte into contract
func DecodeContract(b []byte) (*Contract, error) {
	var contract struct {
		Header []byte
		Code   []byte
	}
	if err := rlp.DecodeBytes(b, &contract); err != nil {
		return nil, err
	}

	header, err := DecodeHeader(contract.Header)
	if err != nil {
		return nil, err
	}
	return &Contract{header, contract.Code}, nil
}

// EncodeRLP encodes a contract to RLP format
func (c *Contract) EncodeRLP(w io.Writer) error {
	contractHeader, _ := rlp.EncodeToBytes(c.Header)
	return rlp.Encode(w, struct {
		Header []byte
		Code   []byte
	}{
		Header: contractHeader,
		Code:   c.Code,
	})
}

// MarshalJSON returns json string of contract
func (c *Contract) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Header *Header `json:"header"`
		Code   string  `json:"code"`
	}{
		Header: c.Header,
		Code:   hex.EncodeToString(c.Code),
	})
}
