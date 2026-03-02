package tests

import (
	"github.com/adrianland/mi-proyecto-blog-api/internal/domain"
)

// MockArticleRepository para tests
type MockArticleRepository struct {
	articles map[int]*domain.Article
	nextID   int
}

func NewMockArticleRepository() *MockArticleRepository {
	return &MockArticleRepository{
		articles: make(map[int]*domain.Article),
		nextID:   1,
	}
}

func (m *MockArticleRepository) Create(article *domain.Article) error {
	article.ID = m.nextID
	m.articles[m.nextID] = article
	m.nextID++
	return nil
}

func (m *MockArticleRepository) GetByID(id int) (*domain.Article, error) {
	if article, exists := m.articles[id]; exists {
		return article, nil
	}
	return nil, nil
}

func (m *MockArticleRepository) Update(article *domain.Article) error {
	if _, exists := m.articles[article.ID]; !exists {
		return nil
	}
	m.articles[article.ID] = article
	return nil
}

func (m *MockArticleRepository) ListPublished(page, pageSize int) ([]domain.ArticleWithAuthor, int64, error) {
	return []domain.ArticleWithAuthor{}, 0, nil
}

func (m *MockArticleRepository) ListByAuthor(authorID int, status string, page, pageSize int) ([]domain.Article, int64, error) {
	return []domain.Article{}, 0, nil
}

func (m *MockArticleRepository) CountPublishedByAuthor(authorID int) (int, error) {
	count := 0
	for _, article := range m.articles {
		if article.AuthorID == authorID && article.Status == domain.StatusPublished {
			count++
		}
	}
	return count, nil
}

func (m *MockArticleRepository) GetAllPublished() ([]domain.Article, error) {
	var articles []domain.Article
	for _, article := range m.articles {
		if article.Status == domain.StatusPublished {
			articles = append(articles, *article)
		}
	}
	return articles, nil
}

// MockAuthorRepository para tests
type MockAuthorRepository struct {
	authors map[int]*domain.Author
	nextID  int
}

func NewMockAuthorRepository() *MockAuthorRepository {
	return &MockAuthorRepository{
		authors: make(map[int]*domain.Author),
		nextID:  1,
	}
}

func (m *MockAuthorRepository) Create(author *domain.Author) error {
	author.ID = m.nextID
	m.authors[m.nextID] = author
	m.nextID++
	return nil
}

func (m *MockAuthorRepository) GetByID(id int) (*domain.Author, error) {
	if author, exists := m.authors[id]; exists {
		return author, nil
	}
	return nil, nil
}

func (m *MockAuthorRepository) GetAll() ([]domain.Author, error) {
	return []domain.Author{}, nil
}

func (m *MockAuthorRepository) GetSummary(authorID int) (*domain.AuthorSummary, error) {
	return &domain.AuthorSummary{}, nil
}
