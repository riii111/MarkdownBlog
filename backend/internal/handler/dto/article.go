package dto

type CreateArticleRequest struct {
	// 空のリクエストボディを表現するため、フィールドなし
}

type CreateArticleResponse struct {
	Slug string `json:"slug"`
}
