package chain

import (
	"github.com/QuoineFinancial/liquid-chain/api/resource"
	"github.com/QuoineFinancial/liquid-chain/db"
)

// Service is first service
type Service struct {
	tAPI     resource.TendermintAPI
	database db.Database
}

// NewService returns new instance of Service
func NewService(tAPI resource.TendermintAPI, database db.Database) *Service {
	return &Service{tAPI, database}
}
