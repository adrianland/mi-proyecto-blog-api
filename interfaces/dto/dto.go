package dto

import "time"

// Request DTOs

type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	AuthorID int    `json:"author_id" binding:"required"`
}

type UpdateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreateAuthorRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type PaginationQuery struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}

type AuthorFilterQuery struct {
	Status   string `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type TopAuthorsQuery struct {
	N int `form:"n,default=3"`
}

// Response DTOs

type ArticleResponse struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Status          string     `json:"status"`
	AuthorID        int        `json:"author_id"`
	AuthorName      string     `json:"author_name,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	PublicationDate *time.Time `json:"publication_date"`
	Score           float64    `json:"score,omitempty"`
	WordCount       int        `json:"word_count,omitempty"`
}

type AuthorResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthorSummaryResponse struct {
	AuthorID            int        `json:"author_id"`
	AuthorName          string     `json:"author_name"`
	TotalArticles       int        `json:"total_articles"`
	TotalPublished      int        `json:"total_published"`
	LastPublicationDate *time.Time `json:"last_publication_date"`
}

type TopAuthorResponse struct {
	AuthorID         int     `json:"author_id"`
	AuthorName       string  `json:"author_name"`
	ScoreAccumulated float64 `json:"score_accumulated"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
