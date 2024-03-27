package models

type (
	RegisterUserReq struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		HashPass string `json:"password" validate:"required"`
	}
)
