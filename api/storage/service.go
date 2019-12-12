package storage

import (
	"github.com/QuoineFinancial/vertex/api/resource"
	"github.com/QuoineFinancial/vertex/db"
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
