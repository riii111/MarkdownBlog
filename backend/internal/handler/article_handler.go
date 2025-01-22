package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
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
