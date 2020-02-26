package event

import (
	"encoding/binary"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
)

// Event is instance for vm to emit
type Event struct {
	*abi.Event
	Values          []byte
	ContractAddress *crypto.Address
}

// NewDeploymentEvent returns event when deploy contract
func NewDeploymentEvent(address crypto.Address) Event {
	values, _ := abi.Encode(deploymentEventABI.Parameters, []interface{}{address})
	return Event{&deploymentEventABI, values, nil}
}

// NewDetailsEvent returns extra transactions details
func NewDetailsEvent(height uint64, from crypto.Address, to crypto.Address, nonce uint64, result uint64) Event {
	values, _ := abi.Encode(detailEventABI.Parameters, []interface{}{height, from, to, nonce, result})
	return Event{&detailEventABI, values, nil}
}

// NewCustomEvent return event declared by user
func NewCustomEvent(event *abi.Event, values []byte, contract crypto.Address) Event {
	return Event{event, values, &contract}
}

// ParseCustomEventName return the crypto.Adress and index of an event name
func ParseCustomEventName(name []byte) (*crypto.Address, uint32, error) {
	address, err := crypto.AddressFromBytes(name[0:crypto.AddressLength])
	if err != nil {
		return nil, 0, err
	}
	index := binary.LittleEndian.Uint32(name[crypto.AddressLength:])
	return &address, index, nil
}
