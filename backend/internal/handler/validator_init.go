// カスタムバリデーションの登録処理Wrapper

package handler

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ValidationRegistry struct {
	validator *validator.Validate
}

func NewValidationRegistry() (*ValidationRegistry, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return &ValidationRegistry{validator: v}, nil
	}
	return nil, fmt.Errorf("failed to get validator engine")
}

// アプリケーション全体のカスタムバリデーションを登録する
func (vr *ValidationRegistry) RegisterCustomValidations() error {
	if err := vr.registerEmailValidations(); err != nil {
		return err
	}

	return nil
}
