package models

// Pagination is model for result in list
type Pagination struct {
	LastPage    int `json:"last"`
	CurrentPage int `json:"current"`
	Total       int `json:"total"`
}
