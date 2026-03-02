package handler

import (
	"net/http"
	"strconv"

	"github.com/adrianland/mi-proyecto-blog-api/interfaces/dto"
	"github.com/adrianland/mi-proyecto-blog-api/internal/application"
	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	authorService  *application.AuthorService
	articleService *application.ArticleService
}

func NewAuthorHandler(authorService *application.AuthorService, articleService *application.ArticleService) *AuthorHandler {
	return &AuthorHandler{
		authorService:  authorService,
		articleService: articleService,
	}
}

// CreateAuthor POST /autores
func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var req dto.CreateAuthorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	author, err := h.authorService.CreateAuthor(req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "CREATION_ERROR",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response := dto.AuthorResponse{
		ID:        author.ID,
		Name:      author.Name,
		Email:     author.Email,
		CreatedAt: author.CreatedAt,
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Author created successfully",
		Data:    response,
	})
}

// GetAuthorSummary GET /autores/:id/resumen
func (h *AuthorHandler) GetAuthorSummary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid author ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	summary, err := h.authorService.GetAuthorSummary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "NOT_FOUND",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	response := dto.AuthorSummaryResponse{
		AuthorID:            summary.AuthorID,
		AuthorName:          summary.AuthorName,
		TotalArticles:       summary.TotalArticles,
		TotalPublished:      summary.TotalPublished,
		LastPublicationDate: summary.LastPublicationDate,
	}

	c.JSON(http.StatusOK, response)
}

// GetTopAuthors GET /autores/top
func (h *AuthorHandler) GetTopAuthors(c *gin.Context) {
	var query dto.TopAuthorsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		query.N = 3
	}

	if query.N < 1 {
		query.N = 3
	}

	topAuthors, err := h.articleService.GetTopAuthors(query.N)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "LIST_ERROR",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	responses := make([]dto.TopAuthorResponse, len(topAuthors))
	for i, author := range topAuthors {
		responses[i] = dto.TopAuthorResponse{
			AuthorID:         author.AuthorID,
			AuthorName:       author.AuthorName,
			ScoreAccumulated: author.ScoreAccumulated,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// ListArticlesByAuthor GET /autores/:id/articulos
func (h *AuthorHandler) ListArticlesByAuthor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid author ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var query dto.AuthorFilterQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		query.Page = 1
		query.PageSize = 10
	}

	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}
	if query.Page < 1 {
		query.Page = 1
	}

	articles, total, err := h.articleService.ListArticlesByAuthor(id, query.Status, query.Page, query.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "LIST_ERROR",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	responses := make([]dto.ArticleResponse, len(articles))
	for i, article := range articles {
		responses[i] = dto.ArticleResponse{
			ID:              article.ID,
			Title:           article.Title,
			Content:         article.Content,
			Status:          article.Status,
			AuthorID:        article.AuthorID,
			CreatedAt:       article.CreatedAt,
			PublicationDate: article.PublicationDate,
			Score:           article.Score,
			WordCount:       article.WordCount,
		}
	}

	totalPages := (total + int64(query.PageSize) - 1) / int64(query.PageSize)

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:       responses,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	})
}
