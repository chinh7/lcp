package event

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/tendermint/tendermint/abci/types"
)

// SystemEventCode represents the code of system event
type SystemEventCode byte

// DetailEvent emitting when a tx is mined
type DetailEvent struct {
	Height   uint64
	From     crypto.Address
	To       crypto.Address
	Nonce    uint64
	Result   uint64
	GasPrice uint64
}

// DeploymentEvent emitting when contract is deployed
type DeploymentEvent struct {
	Address crypto.Address
}

func LoadDetailEvent(tmEvent types.Event) *DetailEvent {
	decodedValues := make([][]byte, len(tmEvent.GetAttributes()))
	for index, attribute := range tmEvent.GetAttributes() {
		decodedValues[index], _ = hex.DecodeString(string(attribute.Value))
	}

	return &DetailEvent{
		Height:   binary.LittleEndian.Uint64(decodedValues[0]),
		From:     crypto.AddressFromBytes(decodedValues[1]),
		To:       crypto.AddressFromBytes(decodedValues[2]),
		Nonce:    binary.LittleEndian.Uint64(decodedValues[3]),
		Result:   binary.LittleEndian.Uint64(decodedValues[4]),
		GasPrice: binary.LittleEndian.Uint64(decodedValues[5]),
	}
}

func LoadDeploymentEvent(tmEvent types.Event) *DeploymentEvent {
	addressByte, _ := hex.DecodeString(string(tmEvent.Attributes[0].Value))
	return &DeploymentEvent{
		Address: crypto.AddressFromBytes(addressByte),
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
			Name: "height",
			Type: abi.Uint64,
		},
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
		&abi.Parameter{
			Name: "gas_price",
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

func GetEventCode(event *abi.Event) SystemEventCode {
	return map[string]SystemEventCode{
		systemEventName[Detail]:     Detail,
		systemEventName[Deployment]: Deployment,
	}[event.Name]
}

func (code SystemEventCode) String() string {
	return hex.EncodeToString([]byte{byte(code)})
}
