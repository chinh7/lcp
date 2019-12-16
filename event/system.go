package event

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/tendermint/tendermint/abci/types"
)

type SystemEventCode byte

type DetailEvent struct {
	From   crypto.Address
	To     crypto.Address
	Nonce  uint64
	Result uint64
}

type DeploymentEvent struct {
	Address crypto.Address
}

func LoadDetailEvent(tmEvent types.Event) *DetailEvent {
	return &DetailEvent{
		From:   crypto.AddressFromBytes(tmEvent.Attributes[0].Value),
		To:     crypto.AddressFromBytes(tmEvent.Attributes[1].Value),
		Nonce:  binary.LittleEndian.Uint64(tmEvent.Attributes[2].Value),
		Result: binary.LittleEndian.Uint64(tmEvent.Attributes[3].Value),
	}
}

func LoadDeploymentEvent(tmEvent types.Event) *DeploymentEvent {
	return &DeploymentEvent{
		Address: crypto.AddressFromBytes(tmEvent.Attributes[0].Value),
	}
}

const (
	Detail     SystemEventCode = 0x1
	Deployment SystemEventCode = 0x2
)

var systemEventName map[SystemEventCode]string = map[SystemEventCode]string{
	Detail:     "detail",
	Deployment: "deployment",
}

var detailEventABI abi.Event = abi.Event{
	Name: systemEventName[Detail],
	Parameters: []*abi.Parameter{
		&abi.Parameter{
			Name: "from",
			Type: abi.Address,
		},
		&abi.Parameter{
			Name: "to",
			Type: abi.Address,
		},
		&abi.Parameter{
			Name: "nonce",
			Type: abi.Uint64,
		},
		&abi.Parameter{
			Name: "result",
			Type: abi.Uint64,
		},
	},
}

var deploymentEventABI abi.Event = abi.Event{
	Name: systemEventName[Deployment],
	Parameters: []*abi.Parameter{
		&abi.Parameter{
			Name: "address",
			Type: abi.Address,
		},
	},
}

func (code SystemEventCode) GetEvent() *abi.Event {
	return map[SystemEventCode]*abi.Event{
		Detail:     &detailEventABI,
		Deployment: &deploymentEventABI,
	}[code]
}

func getEventCode(event *abi.Event) SystemEventCode {
	return map[string]SystemEventCode{
		systemEventName[Detail]:     Detail,
		systemEventName[Deployment]: Deployment,
	}[event.Name]
}

func (code SystemEventCode) string() string {
	return hex.EncodeToString([]byte{byte(code)})
}
