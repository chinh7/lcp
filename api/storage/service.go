package storage

import "github.com/QuoineFinancial/vertex/api/resource"

// Service is first service
type Service struct {
	tAPI resource.TendermintAPI
}

// NewService returns new instance of Service
func NewService(tAPI resource.TendermintAPI) *Service {
	return &Service{tAPI}
}
