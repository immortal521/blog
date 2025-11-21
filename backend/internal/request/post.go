package request

type SummarizeReq struct {
	Content string `json:"content" validate:"required"`
}
