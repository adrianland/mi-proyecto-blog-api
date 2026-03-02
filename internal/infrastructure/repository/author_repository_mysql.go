package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adrianland/mi-proyecto-blog-api/internal/domain"
)

type AuthorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) domain.AuthorRepository {
	return &AuthorRepository{db: db}
}

// Create crea un nuevo autor
func (r *AuthorRepository) Create(author *domain.Author) error {
	query := `INSERT INTO authors (name, email) VALUES (?, ?)`

	result, err := r.db.Exec(query, author.Name, author.Email)
	if err != nil {
		return fmt.Errorf("error inserting author: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	author.ID = int(id)
	author.CreatedAt = time.Now()
	return nil
}

// GetByID obtiene un autor por ID
func (r *AuthorRepository) GetByID(id int) (*domain.Author, error) {
	query := `SELECT id, name, email, created_at FROM authors WHERE id = ?`

	author := &domain.Author{}
	err := r.db.QueryRow(query, id).Scan(&author.ID, &author.Name, &author.Email, &author.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("author not found")
		}
		return nil, fmt.Errorf("error querying author: %w", err)
	}

	return author, nil
}

// GetAll obtiene todos los autores
func (r *AuthorRepository) GetAll() ([]domain.Author, error) {
	query := `SELECT id, name, email, created_at FROM authors ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying authors: %w", err)
	}
	defer rows.Close()

	var authors []domain.Author
	for rows.Next() {
		author := domain.Author{}
		err := rows.Scan(&author.ID, &author.Name, &author.Email, &author.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning author: %w", err)
		}
		authors = append(authors, author)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating authors: %w", err)
	}

	return authors, nil
}

// GetSummary obtiene el resumen de estadísticas de un autor
func (r *AuthorRepository) GetSummary(authorID int) (*domain.AuthorSummary, error) {
	// Verificar que el autor existe
	author, err := r.GetByID(authorID)
	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	query := `
	SELECT 
		COUNT(*) as total_articles,
		SUM(CASE WHEN status = 'PUBLICADO' THEN 1 ELSE 0 END) as total_published,
		MAX(CASE WHEN status = 'PUBLICADO' THEN publication_date END) as last_publication
	FROM articles
	WHERE author_id = ?
	`

	summary := &domain.AuthorSummary{
		AuthorID:   author.ID,
		AuthorName: author.Name,
	}

	var totalArticles sql.NullInt64
	var totalPublished sql.NullInt64
	var lastPublication sql.NullTime

	err = r.db.QueryRow(query, authorID).Scan(&totalArticles, &totalPublished, &lastPublication)
	if err != nil {
		return nil, fmt.Errorf("error querying summary: %w", err)
	}

	if totalArticles.Valid {
		summary.TotalArticles = int(totalArticles.Int64)
	}
	if totalPublished.Valid {
		summary.TotalPublished = int(totalPublished.Int64)
	}
	if lastPublication.Valid {
		summary.LastPublicationDate = &lastPublication.Time
	}

	return summary, nil
}
