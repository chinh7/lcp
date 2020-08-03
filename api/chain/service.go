package chain

import (
	"github.com/QuoineFinancial/liquid-chain/api/resource"
	"github.com/QuoineFinancial/liquid-chain/consensus"
)

// Service is first service
type Service struct {
	tAPI resource.TendermintAPI
	app  *consensus.App
}

// NewService returns new instance of Service
func NewService(tAPI resource.TendermintAPI, app *consensus.App) *Service {
	return &Service{tAPI, app}
}
