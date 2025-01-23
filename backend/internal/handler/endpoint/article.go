package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/config"
	"github.com/riii111/markdown-blog-api/internal/usecase"
)

type ArticleHandler struct {
	articleUsecase *usecase.ArticleUsecase
}

func NewArticleHandler(articleUsecase *usecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		articleUsecase: articleUsecase,
	}
}

// CreateArticle godoc
// @Summary Create a new article
// @Description Create a new blog article
// @Tags articles
// @Accept json
// @Produce json
// @Param request body dto.CreateArticleRequest true "Article creation data"
// @Success 201 {object} dto.CreateArticleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	article, err := h.articleUsecase.CreateArticle(
		c.Request.Context(),
		userID,
		"",
		"",
		nil,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response := dto.CreateArticleResponse{
		Slug: article.Slug,
	}
	c.JSON(http.StatusCreated, response)
}

// DeleteArticle godoc
// @Summary Delete an article
// @Description Delete a blog article by slug
// @Tags articles
// @Produce json
// @Param slug path string true "Article slug"
// @Success 204 "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/articles/{slug} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	slug := c.Param("slug")

	if err := h.articleUsecase.DeleteArticleBySlug(c.Request.Context(), userID, slug); err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Article not found",
			})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Unauthorized",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete article",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetArticles godoc
// @Summary Get published articles
// @Description Get published articles with cursor-based pagination
// @Tags articles
// @Produce json
// @Param limit query int false "Number of articles per page"
// @Param cursor query string false "Cursor for pagination"
// @Success 200 {object} dto.ArticleListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/articles [get]
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	limit := config.GetDefaultPerPage()
	var cursor *string

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if cursorStr := c.Query("cursor"); cursorStr != "" {
		cursor = &cursorStr
	}

	// limitの検証
	if limit <= 0 || limit > config.GetMaxPerPage() {
		limit = config.GetDefaultPerPage()
	}

	articles, nextCursor, err := h.articleUsecase.GetPublishedArticles(c.Request.Context(), limit, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to fetch articles",
		})
		return
	}

	response := dto.NewArticleListResponse(articles, limit, nextCursor)
	c.JSON(http.StatusOK, response)
}

// GetArticle godoc
// @Summary Get an article detail
// @Description Get a blog article by slug
// @Tags articles
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} dto.ArticleDetailResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/articles/{slug} [get]
func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")

	article, err := h.articleUsecase.GetArticleBySlug(c.Request.Context(), slug)
	if err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Article not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to fetch article",
		})
		return
	}

	c.JSON(http.StatusOK, dto.NewArticleDetailResponse(article))
}
