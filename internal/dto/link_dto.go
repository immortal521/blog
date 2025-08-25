package dto

type LinkCreateDTO struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

type LinkResponseDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	SortOrder   int    `json:"sortOrder"`
}
