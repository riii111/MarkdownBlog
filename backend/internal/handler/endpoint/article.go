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
// @Description Get published articles
// @Tags articles
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Number of articles per page"
// @Success 200 {object} dto.ArticleListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/articles [get]
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	page := config.GetDefaultPage()
	perPage := config.GetDefaultPerPage()

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil {
			perPage = pp
		}
	}

	// バリデーション
	page, perPage = config.ValidateAndAdjustPagination(page, perPage)

	articles, total, err := h.articleUsecase.GetPublishedArticles(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to fetch articles",
		})
		return
	}

	c.JSON(http.StatusOK, dto.NewArticleListResponse(articles, page, perPage, total))
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
