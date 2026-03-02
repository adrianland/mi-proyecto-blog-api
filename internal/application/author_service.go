package application

import (
	"fmt"

	"github.com/adrianland/mi-proyecto-blog-api/internal/domain"
)

type AuthorService struct {
	authorRepo domain.AuthorRepository
}

func NewAuthorService(authorRepo domain.AuthorRepository) *AuthorService {
	return &AuthorService{
		authorRepo: authorRepo,
	}
}

// CreateAuthor crea un nuevo autor
func (s *AuthorService) CreateAuthor(name, email string) (*domain.Author, error) {
	if name == "" || email == "" {
		return nil, fmt.Errorf("name and email are required")
	}

	author := &domain.Author{
		Name:  name,
		Email: email,
	}

	if err := s.authorRepo.Create(author); err != nil {
		return nil, fmt.Errorf("failed to create author: %w", err)
	}

	return author, nil
}

// GetAuthor obtiene un autor por ID
func (s *AuthorService) GetAuthor(id int) (*domain.Author, error) {
	author, err := s.authorRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}
	return author, nil
}

// GetAuthorSummary obtiene el resumen de un autor
func (s *AuthorService) GetAuthorSummary(authorID int) (*domain.AuthorSummary, error) {
	summary, err := s.authorRepo.GetSummary(authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author summary: %w", err)
	}
	return summary, nil
}

// ListAllAuthors lista todos los autores
func (s *AuthorService) ListAllAuthors() ([]domain.Author, error) {
	authors, err := s.authorRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list authors: %w", err)
	}
	return authors, nil
}
