package dto

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	Email       string `json:"email" binding:"required,email" example:"user@example.com"`
	Password    string `json:"password" binding:"required,min=8,max=32" example:"password"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100" example:"John Doe"`
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
