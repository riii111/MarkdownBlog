package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/internal/usecase"
)

type PostHandler struct {
	postUsecase *usecase.PostUsecase
}

func NewPostHandler(postUsecase *usecase.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new blog post
// @Tags posts
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "Post creation data"
// @Success 201 {object} dto.CreatePostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	// ユーザーIDの取得（認証ミドルウェアで設定されていることを前提）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	post, err := h.postUsecase.CreatePost(
		c.Request.Context(),
		userID.(uuid.UUID),
		req.Title,
		req.Content,
		req.SeriesID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.CreatePostResponse{
		Slug: post.Slug,
	})
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a blog post by slug
// @Tags posts
// @Produce json
// @Param slug path string true "Post slug"
// @Success 204 "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /posts/{slug} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid post ID",
		})
		return
	}

	if err := h.postUsecase.DeletePost(c.Request.Context(), userID.(uuid.UUID), postID); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
