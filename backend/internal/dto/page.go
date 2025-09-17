package dto

type Page[T any] struct {
	Total uint `json:"total"`
	Page  uint `json:"page"`
	Limit uint `json:"limit"`
	Data  []T  `json:"data"`
}
