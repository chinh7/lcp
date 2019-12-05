package crypto

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"

	"github.com/pkg/errors"
	"github.com/stellar/go/crc16"
	"golang.org/x/crypto/ed25519"
)

const (
	versionByteAccountID byte = 11 << 3 // Base32-encodes to 'L...'
	// AddressLength size of a crypto address
	AddressLength = 35
)

// Address crypto address
type Address [AddressLength]byte

func (address *Address) setBytes(b []byte) {
	if len(b) > len(address) {
		b = b[len(b)-AddressLength:]
	}
	copy(address[AddressLength-len(b):], b)
}

// String Address string presentation
func (address *Address) String() string {
	return base32.StdEncoding.EncodeToString(address[:])
}

// AddressFromPubKey create an address from public key
func AddressFromPubKey(src []byte) Address {
	version := versionByteAccountID
	var raw bytes.Buffer

	// write version byte
	if err := binary.Write(&raw, binary.LittleEndian, version); err != nil {
		return [AddressLength]byte{}
	}

	// write payload
	if _, err := raw.Write(src); err != nil {
		return [AddressLength]byte{}
	}

	// calculate and write checksum
	checksum := crc16.Checksum(raw.Bytes())
	if _, err := raw.Write(checksum); err != nil {
		return [AddressLength]byte{}
	}
	var address Address
	address.setBytes(raw.Bytes())
	return address
}

// AddressFromString parse an address string to Address
func AddressFromString(address string) Address {
	pubkeyString, err := decodeAddress(versionByteAccountID, address)
	if err != nil {
		panic(err)
	}
	pubkey := ed25519.PublicKey(pubkeyString)
	return AddressFromPubKey(pubkey)
}

func decodeString(src string) ([]byte, error) {
	raw, err := base32.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, errors.Wrap(err, "base32 decode failed")
	}

	if len(raw) < 3 {
		return nil, errors.Errorf("encoded value is %d bytes; minimum valid length is 3", len(raw))
	}

	return raw, nil
}

func decodeAddress(expected byte, src string) ([]byte, error) {
	raw, err := decodeString(src)
	if err != nil {
		return nil, err
	}

	// decode into components
	version := byte(raw[0])
	vp := raw[0 : len(raw)-2]
	payload := raw[1 : len(raw)-2]
	checksum := raw[len(raw)-2:]

	if version != expected {
		panic("Unexpected version")
	}

	// ensure checksum is valid
	if err := crc16.Validate(vp, checksum); err != nil {
		return nil, err
	}

	// if we made it through the gaunlet, return the decoded value
	return payload, nil
}
