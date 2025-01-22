package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/riii111/markdown-blog-api/internal/domain/model"
	"golang.org/x/crypto/bcrypt"
)

// カスタムエラーの定義
var (
	ErrEmailAlreadyExists = fmt.Errorf("email already exists")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrRegistrationFailed = fmt.Errorf("registration failed")
	ErrSessionNotFound    = fmt.Errorf("session not found")
)

type UserUsecase struct {
	userRepo    model.UserRepository
	sessionRepo model.SessionRepository
}

func NewUserUsecase(userRepo model.UserRepository, sessionRepo model.SessionRepository) *UserUsecase {
	return &UserUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// ユーザー登録
func (u *UserUsecase) Register(ctx context.Context, email, password, displayName string) (*model.User, error) {
	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return nil, ErrRegistrationFailed
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		DisplayName:  displayName,
	}

	// ユーザーの保存（ユニーク制約違反を含むすべてのエラーを一般的なエラーとして扱う）
	if err := u.userRepo.Create(ctx, user); err != nil {
		// エラーの詳細はログに記録するが、呼び出し元には一般的なエラーのみを返す
		log.Printf("failed to create user: %v", err)
		return nil, ErrRegistrationFailed
	}

	return user, nil
}

// ログイン
func (u *UserUsecase) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Printf("FindByEmail error: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// ログアウト
func (u *UserUsecase) Logout(ctx context.Context, sessionToken string) error {
	// セッションを検索
	session, err := u.sessionRepo.FindByToken(sessionToken)
	if err != nil {
		return err
	}
	if session == nil {
		return ErrSessionNotFound
	}

	// セッションを削除
	return u.sessionRepo.Delete(session.ID)
}

// セッショントークンの生成
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
