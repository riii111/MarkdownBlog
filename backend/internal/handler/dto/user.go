package dto

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type RegisterUserRequest struct {
	Email       string `json:"email" binding:"required,email"` // emailタグでバリデーションが行える（Ginのvalidatorパッケージの仕様）
	Password    string `json:"password" binding:"required,min=8, max=32"`
	DisplayName string `json:"display_name" binding:"required, min=1, max=100"`
}

var (
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+-=[]{}|;:,.<>?]+$`)
)

func RegisterCustomValidations(v *validator.Validate) {
	// パスワードのバリデーション
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return passwordRegex.MatchString(fl.Field().String())
	})
}

type RegisterUserResponse struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"` // UI表示用
	CreatedAt   time.Time `json:"created_at"`
}
