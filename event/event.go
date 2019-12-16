package event

import (
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
func NewDetailsEvent(from crypto.Address, to crypto.Address, nonce uint64, result uint64) Event {
	values, _ := abi.Encode(detailEventABI.Parameters, []interface{}{from, to, nonce, result})
	return Event{&detailEventABI, values, nil}
}

// NewCustomEvent return event declared by user
func NewCustomEvent(event *abi.Event, values []byte, contract crypto.Address) Event {
	return Event{event, values, &contract}
}
