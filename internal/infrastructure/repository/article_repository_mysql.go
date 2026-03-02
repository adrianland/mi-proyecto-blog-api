package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/adrianland/mi-proyecto-blog-api/internal/domain"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) domain.ArticleRepository {
	return &ArticleRepository{db: db}
}

// Create crea un nuevo artículo
func (r *ArticleRepository) Create(article *domain.Article) error {
	query := `INSERT INTO articles (title, content, status, author_id, created_at) 
	          VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, article.Title, article.Content, article.Status, article.AuthorID, article.CreatedAt)
	if err != nil {
		return fmt.Errorf("error inserting article: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	article.ID = int(id)
	return nil
}

// GetByID obtiene un artículo por ID
func (r *ArticleRepository) GetByID(id int) (*domain.Article, error) {
	query := `SELECT id, title, content, status, author_id, created_at, publication_date 
	          FROM articles WHERE id = ?`

	article := &domain.Article{}
	err := r.db.QueryRow(query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Status,
		&article.AuthorID,
		&article.CreatedAt,
		&article.PublicationDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		return nil, fmt.Errorf("error querying article: %w", err)
	}

	// Contar palabras
	article.WordCount = len(strings.Fields(article.Content))

	return article, nil
}

// Update actualiza un artículo
func (r *ArticleRepository) Update(article *domain.Article) error {
	query := `UPDATE articles SET title = ?, content = ?, status = ?, publication_date = ? 
	          WHERE id = ?`

	result, err := r.db.Exec(query, article.Title, article.Content, article.Status, article.PublicationDate, article.ID)
	if err != nil {
		return fmt.Errorf("error updating article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

// ListPublished lista artículos publicados con paginación
func (r *ArticleRepository) ListPublished(page, pageSize int) ([]domain.ArticleWithAuthor, int64, error) {
	// Contar total
	countQuery := `SELECT COUNT(*) FROM articles WHERE status = 'PUBLICADO'`
	var total int64
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting articles: %w", err)
	}

	// Obtener artículos
	offset := (page - 1) * pageSize
	query := `SELECT a.id, a.title, a.content, a.status, a.author_id, au.name, 
	                 a.created_at, a.publication_date
	          FROM articles a
	          JOIN authors au ON a.author_id = au.id
	          WHERE a.status = 'PUBLICADO'
	          ORDER BY a.publication_date DESC
	          LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying articles: %w", err)
	}
	defer rows.Close()

	var articles []domain.ArticleWithAuthor
	for rows.Next() {
		article := domain.ArticleWithAuthor{}
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Status,
			&article.AuthorID,
			&article.AuthorName,
			&article.CreatedAt,
			&article.PublicationDate,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning article: %w", err)
		}
		article.WordCount = len(strings.Fields(article.Content))
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating articles: %w", err)
	}

	return articles, total, nil
}

// ListByAuthor lista artículos por autor
func (r *ArticleRepository) ListByAuthor(authorID int, status string, page, pageSize int) ([]domain.Article, int64, error) {
	// Construir query dinámicamente
	whereClause := "WHERE author_id = ?"
	args := []interface{}{authorID}

	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	// Contar total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM articles %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting articles: %w", err)
	}

	// Obtener artículos
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(
		`SELECT id, title, content, status, author_id, created_at, publication_date
		 FROM articles %s
		 ORDER BY created_at DESC
		 LIMIT ? OFFSET ?`,
		whereClause,
	)

	args = append(args, pageSize, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying articles: %w", err)
	}
	defer rows.Close()

	var articles []domain.Article
	for rows.Next() {
		article := domain.Article{}
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Status,
			&article.AuthorID,
			&article.CreatedAt,
			&article.PublicationDate,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning article: %w", err)
		}
		article.WordCount = len(strings.Fields(article.Content))
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating articles: %w", err)
	}

	return articles, total, nil
}

// CountPublishedByAuthor cuenta los artículos publicados de un autor
func (r *ArticleRepository) CountPublishedByAuthor(authorID int) (int, error) {
	query := `SELECT COUNT(*) FROM articles WHERE author_id = ? AND status = 'PUBLICADO'`

	var count int
	err := r.db.QueryRow(query, authorID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting published articles: %w", err)
	}

	return count, nil
}

// GetAllPublished obtiene todos los artículos publicados (sin paginación)
func (r *ArticleRepository) GetAllPublished() ([]domain.Article, error) {
	query := `SELECT id, title, content, status, author_id, created_at, publication_date
	          FROM articles
	          WHERE status = 'PUBLICADO'
	          ORDER BY publication_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying articles: %w", err)
	}
	defer rows.Close()

	var articles []domain.Article
	for rows.Next() {
		article := domain.Article{}
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Status,
			&article.AuthorID,
			&article.CreatedAt,
			&article.PublicationDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning article: %w", err)
		}
		article.WordCount = len(strings.Fields(article.Content))
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating articles: %w", err)
	}

	return articles, nil
}
