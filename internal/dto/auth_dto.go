package dto

type CaptchaReq struct {
	Email string `json:"email" validate:"required,email"`
	Type  string `json:"type" validate:"omitempty,oneof=Register PasswordReset ChangeEmail"`
}
