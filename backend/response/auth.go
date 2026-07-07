package response

// LoginRes is the response body for login/register endpoints.
// Role is plain string instead of entity.UserRole to decouple
// the response layer from the domain model.
type LoginRes struct {
	AccessToken  string  `json:"accessToken"`
	RefreshToken string  `json:"-"`
	UUID         string  `json:"uuid"`
	Avatar       *string `json:"avatar"`
	Username     string  `json:"username"`
	Role         string  `json:"role"`
}

type RefreshRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"-"`
}
