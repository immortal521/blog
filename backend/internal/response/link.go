package response

type LinkResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	SortOrder   int    `json:"sortOrder"`
}
