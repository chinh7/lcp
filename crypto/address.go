package crypto

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"log"

	"github.com/QuoineFinancial/liquid-chain-rlp/rlp"
	"github.com/QuoineFinancial/liquid-chain/crc16"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
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
		log.Fatalf("byte write should not fail %v", err)
	}

	// write payload
	if _, err := raw.Write(src); err != nil {
		log.Fatalf("byte write should not fail %v", err)
	}

	// calculate and write checksum
	checksum := crc16.Checksum(raw.Bytes())
	if _, err := raw.Write(checksum); err != nil {
		log.Fatalf("byte write should not fail %v", err)
	}
	var address Address
	address.setBytes(raw.Bytes())
	return address
}

// AddressFromString parse an address string to Address
func AddressFromString(address string) (Address, error) {
	pubKeyBytes, err := decodeAddress(versionByteAccountID, address)
	if err != nil {
		return Address{}, err
	}
	pubkey := ed25519.PublicKey(pubKeyBytes)
	return AddressFromPubKey(pubkey), nil
}

// AddressFromBytes return an address given its bytes
func AddressFromBytes(b []byte) (Address, error) {
	var a Address
	_, err := decodeAddressBytes(versionByteAccountID, b)
	if err != nil {
		return Address{}, err
	}
	a.setBytes(b)
	return a, nil
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
	return decodeAddressBytes(expected, raw)
}

func decodeAddressBytes(expected byte, raw []byte) ([]byte, error) {
	version := byte(raw[0])
	payload := raw[1 : len(raw)-2]
	checksum := raw[len(raw)-2:]
	original := raw[0 : len(raw)-2]

	if version != expected {
		return nil, errors.Errorf("Unexpected version %x", version)
	}

	// checksum check
	if err := crc16.Validate(original, checksum); err != nil {
		return nil, err
	}
	return payload, nil
}

// NewDeploymentAddress returns new contract deployment address
func NewDeploymentAddress(senderAddress Address, senderNonce uint64) Address {
	senderBytes, _ := rlp.EncodeToBytes([]interface{}{senderAddress, senderNonce})
	res := blake2b.Sum256(senderBytes)
	address := AddressFromPubKey(res[:])
	return address
}
