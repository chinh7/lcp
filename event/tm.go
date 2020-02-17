package event

import (
	"encoding/hex"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

// ToTMEvent convert vertex event to tendermint event
func (event *Event) ToTMEvent() types.Event {
	attributes := []common.KVPair{}
	decodedParams, _ := abi.DecodeToBytes(event.Event.Parameters, event.Values)
	for index, param := range decodedParams {
		attributes = append(attributes, common.KVPair{
			Key:   []byte(hex.EncodeToString([]byte{byte(index)})),
			Value: []byte(hex.EncodeToString(param)),
		})
	}
	var eventName string
	if event.ContractAddress != nil {
		eventName = hex.EncodeToString(append(event.ContractAddress[:], event.GetIndexByte()...))
	} else {
		eventName = GetEventCode(event.Event).String()
	}
	return types.Event{
		Type:       eventName,
		Attributes: attributes,
	}
}
