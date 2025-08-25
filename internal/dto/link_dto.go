package dto

type LinkCreateReq struct {
	Name        string `json:"name" validate:"required"`
	Url         string `json:"url" validate:"required,url"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

type LinkRes struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	SortOrder   int    `json:"sortOrder"`
}
