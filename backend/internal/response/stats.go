package response

type DashboardStatsRes struct {
	TotalPV       int64 `json:"totalPV"`
	TotalUV       int64 `json:"totalUV"`
	TotalComments int64 `json:"totalComments"`
	TotalPosts    int64 `json:"totalPosts"`
}
