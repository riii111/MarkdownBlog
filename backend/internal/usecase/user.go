package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/riii111/markdown-blog-api/internal/domain/model"
	"golang.org/x/crypto/bcrypt"
)

// カスタムエラーの定義
var (
	ErrEmailAlreadyExists = fmt.Errorf("email already exists")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrRegistrationFailed = fmt.Errorf("registration failed")
)

type UserUsecase struct {
	userRepo model.UserRepository
}

func NewUserUsecase(userRepo model.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) Register(ctx context.Context, email, password, displayName string) (*model.User, error) {
	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return nil, ErrRegistrationFailed
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

	// ユーザーの保存（ユニーク制約違反を含むすべてのエラーを一般的なエラーとして扱う）
	if err := u.userRepo.Create(ctx, user); err != nil {
		// エラーの詳細はログに記録するが、呼び出し元には一般的なエラーのみを返す
		log.Printf("failed to create user: %v", err)
		return nil, ErrRegistrationFailed
	}

	return user, nil
}
