package storage

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"

	"github.com/QuoineFinancial/liquid-chain/abi"
	"github.com/QuoineFinancial/liquid-chain/api/models"
	"github.com/QuoineFinancial/liquid-chain/crypto"
	"github.com/QuoineFinancial/liquid-chain/event"
)

func parseEvent(liquidEvent *event.Event) *models.Event {
	attributes := []models.EventAttribute{}
	tmEvent := liquidEvent.ToTMEvent()
	for index, param := range liquidEvent.Parameters {
		valueByte, _ := hex.DecodeString(string(tmEvent.Attributes[index].Value))
		var value string
		if param.Type == abi.Address {
			address, _ := crypto.AddressFromBytes(valueByte)
			value = address.String()
		} else {
			value = strconv.FormatUint(binary.LittleEndian.Uint64(valueByte), 10)
		}
		attributes = append(attributes, models.EventAttribute{
			Key:   param.Name,
			Type:  param.Type.String(),
			Value: value,
		})
	}
	return &models.Event{
		Name:       liquidEvent.Name,
		Contract:   liquidEvent.ContractAddress.String(),
		Attributes: attributes,
	}
}
