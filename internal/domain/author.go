package domain

import (
	"time"
)

type Author struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthorSummary struct {
	AuthorID            int        `json:"author_id"`
	AuthorName          string     `json:"author_name"`
	TotalArticles       int        `json:"total_articles"`
	TotalPublished      int        `json:"total_published"`
	LastPublicationDate *time.Time `json:"last_publication_date"`
}

type TopAuthor struct {
	AuthorID         int     `json:"author_id"`
	AuthorName       string  `json:"author_name"`
	ScoreAccumulated float64 `json:"score_accumulated"`
}
