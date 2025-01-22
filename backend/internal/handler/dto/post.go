package dto

type CreatePostRequest struct {
	// 空のリクエストボディを表現するため、フィールドなし
}

type CreatePostResponse struct {
	Slug string `json:"slug"`
}
