package storage

import (
	"encoding/hex"

	"github.com/QuoineFinancial/liquid-chain/api/models"
	"github.com/QuoineFinancial/liquid-chain/event"
)

func parseEvent(liquidEvent *event.Event) *models.Event {
	attributes := []models.EventAttribute{}
	tmEvent := liquidEvent.ToTMEvent()
	for index, param := range liquidEvent.Parameters {
		attributes = append(attributes, models.EventAttribute{
			Key:   param.Name,
			Type:  param.Type.String(),
			Value: hex.EncodeToString(tmEvent.Attributes[index].Value),
		})
	}
	return &models.Event{
		Name:       liquidEvent.Name,
		Contract:   liquidEvent.ContractAddress.String(),
		Attributes: attributes,
	}
}
