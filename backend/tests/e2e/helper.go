package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// エラーレスポンス検証ヘルパー
// テスト内でのエラーレスポンスの検証を簡略化するためのヘルパー関数
func AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, errorMsgContains string) {
	assert.Equal(t, expectedStatus, w.Code, fmt.Sprintf("Expected status code %d", expectedStatus))

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err, "Failed to unmarshal error response")

	assert.Contains(t, errorResponse, "error", "Response should contain error field")
	if errorMsgContains != "" {
		assert.Contains(t, errorResponse["error"], errorMsgContains,
			fmt.Sprintf("Error message should contain '%s'", errorMsgContains))
	}
}

// ユーザー作成＋ログインのワークフローを簡略化するヘルパー関数
// テストユーザーを作成し、そのままログインまで行う
func CreateAndLoginTestUser(t *testing.T, router *gin.Engine) (dto.RegisterUserResponse, string, string) {
	// テストユーザーデータ
	testEmail := fmt.Sprintf("test-%s@example.com", uuid.New().String())
	testUser := dto.RegisterUserRequest{
		Email:       testEmail,
		Password:    "password123",
		DisplayName: "Test User",
	}

	// ユーザー登録リクエスト
	body, err := json.Marshal(testUser)
	require.NoError(t, err, "Failed to marshal test user")

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")

	// レスポンスをパース
	var response dto.RegisterUserResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response")

	// ログインしてセッショントークンを取得
	loginReq := dto.LoginRequest{
		Email:    testEmail,
		Password: testUser.Password,
	}

	loginBody, err := json.Marshal(loginReq)
	require.NoError(t, err, "Failed to marshal login request")

	loginReqObj := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(loginBody))
	loginReqObj.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReqObj)

	require.Equal(t, http.StatusOK, loginW.Code, "Expected status code 200 for login")

	// セッションCookieを取得
	cookies := loginW.Result().Cookies()
	var sessionToken string
	for _, cookie := range cookies {
		// テスト用のセッション名を使用
		if cookie.Name == "test-session" {
			sessionToken = cookie.Value
			break
		}
	}

	require.NotEmpty(t, sessionToken, "Session token should not be empty")

	// ユーザー情報とセッショントークンを返す
	return response, sessionToken, testEmail
}

// HTTPリクエストを実行するヘルパー関数
func PerformRequest(router *gin.Engine, method, path string, body interface{}, sessionToken string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// セッショントークンがある場合はCookieを設定
	if sessionToken != "" {
		req.AddCookie(&http.Cookie{
			Name:  "test-session",
			Value: sessionToken,
		})
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// テスト用のコンテキストを作成
func CreateTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// レスポンスボディからJSONを取得するヘルパー関数
func GetResponseJSON(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), target)
	require.NoError(t, err, "Failed to unmarshal response body")
}
