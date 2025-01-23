package middleware

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/config"
)

type Store interface {
	NewStore() sessions.Store
}

type cookieStore struct{}

func NewCookieStore() Store {
	return &cookieStore{}
}

func (s *cookieStore) NewStore() sessions.Store {
	sessionKey := []byte(os.Getenv("SESSION_SECRET"))
	if len(sessionKey) < 32 {
		newKey := make([]byte, 32)
		copy(newKey, sessionKey)
		sessionKey = newKey
	}

	store := cookie.NewStore(sessionKey)
	store.Options(sessions.Options{
		Path:     os.Getenv("COOKIE_PATH"),
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		MaxAge:   config.GetSessionMaxAge(),
		Secure:   os.Getenv("COOKIE_SECURE") == "true",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return store
}
