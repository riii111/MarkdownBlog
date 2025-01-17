package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/riii111/markdown-blog-api/internal/domain/model"
	"golang.org/x/crypto/bcrypt"
)

// カスタムエラーの定義
var (
	ErrEmailAlreadyExists = fmt.Errorf("email already exists")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) Register(ctx context.Context, email, password, displayName string) (*model.User, error) {
	// メールアドレスの重複チェック
	existingUser, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// ユーザーモデルの生成
	now := time.Now() // デフォルトはローカル時刻
	user := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		DisplayName:  displayName,
		BaseModel: model.BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// ユーザーの保存
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
