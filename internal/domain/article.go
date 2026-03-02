package domain

import (
	"time"
)

const (
	StatusDraft     = "BORRADOR"
	StatusPublished = "PUBLICADO"
)

type Article struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Status          string     `json:"status"`
	AuthorID        int        `json:"author_id"`
	CreatedAt       time.Time  `json:"created_at"`
	PublicationDate *time.Time `json:"publication_date"`
	Score           float64    `json:"score,omitempty"`
	WordCount       int        `json:"word_count,omitempty"`
}

type ArticleWithAuthor struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	Status          string     `json:"status"`
	AuthorID        int        `json:"author_id"`
	AuthorName      string     `json:"author_name"`
	CreatedAt       time.Time  `json:"created_at"`
	PublicationDate *time.Time `json:"publication_date"`
	Score           float64    `json:"score,omitempty"`
	WordCount       int        `json:"word_count,omitempty"`
}
