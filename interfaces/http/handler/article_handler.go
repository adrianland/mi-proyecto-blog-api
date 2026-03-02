package handler

import (
	"net/http"
	"strconv"

	"github.com/adrianland/mi-proyecto-blog-api/interfaces/dto"
	"github.com/adrianland/mi-proyecto-blog-api/internal/application"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service *application.ArticleService
}

func NewArticleHandler(service *application.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		service: service,
	}
}

// CreateArticle POST /articulos
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req dto.CreateArticleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	article, err := h.service.CreateArticle(req.Title, req.Content, req.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "CREATION_ERROR",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response := dto.ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		Status:    article.Status,
		AuthorID:  article.AuthorID,
		CreatedAt: article.CreatedAt,
		WordCount: len(article.Content) / 5, // Aproximación simple
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Article created successfully",
		Data:    response,
	})
}

// GetArticle GET /articulos/:id
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid article ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	article, err := h.service.GetArticle(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "NOT_FOUND",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	response := dto.ArticleResponse{
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

	c.JSON(http.StatusOK, response)
}

// UpdateArticle PUT /articulos/:id
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid article ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	article, err := h.service.UpdateArticle(id, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "UPDATE_ERROR",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response := dto.ArticleResponse{
		ID:              article.ID,
		Title:           article.Title,
		Content:         article.Content,
		Status:          article.Status,
		AuthorID:        article.AuthorID,
		CreatedAt:       article.CreatedAt,
		PublicationDate: article.PublicationDate,
		WordCount:       article.WordCount,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Article updated successfully",
		Data:    response,
	})
}

// PublishArticle PUT /articulos/:id/publicar
func (h *ArticleHandler) PublishArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Invalid article ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	article, err := h.service.PublishArticle(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "PUBLISH_ERROR",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response := dto.ArticleResponse{
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

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Article published successfully",
		Data:    response,
	})
}

// ListPublishedArticles GET /articulos
func (h *ArticleHandler) ListPublishedArticles(c *gin.Context) {
	var query dto.PaginationQuery
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

	articles, total, err := h.service.ListPublishedArticles(query.Page, query.PageSize)
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
			AuthorName:      article.AuthorName,
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
