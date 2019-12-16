package storage

import (
	"encoding/hex"

	"github.com/QuoineFinancial/liquid-chain/api/models"
	"github.com/QuoineFinancial/liquid-chain/event"
)

func parseEvent(vertexEvent *event.Event) *models.Event {
	attributes := []models.EventAttribute{}
	tmEvent := vertexEvent.ToTMEvent()
	for index, param := range vertexEvent.Parameters {
		attributes = append(attributes, models.EventAttribute{
			Key:   param.Name,
			Type:  param.Type.String(),
			Value: hex.EncodeToString(tmEvent.Attributes[index].Value),
		})
	}
	return &models.Event{
		Name:       vertexEvent.Name,
		Contract:   vertexEvent.ContractAddress.String(),
		Attributes: attributes,
	}
}
