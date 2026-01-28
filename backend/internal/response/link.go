package response

type LinkResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	SortOrder   int    `json:"sortOrder"`
}

type LinkOverview struct {
	Total    int `json:"total"`
	Normal   int `json:"normal"`
	Abnormal int `json:"abnormal"`
	Pending  int `json:"pending"`
}
