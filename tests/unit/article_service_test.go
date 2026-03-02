package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/adrianland/mi-proyecto-blog-api/internal/application"
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

// TESTS

// TestCalculateScore - Prueba unitaria del cálculo de score
func TestCalculateScore(t *testing.T) {
	articleRepo := NewMockArticleRepository()
	authorRepo := NewMockAuthorRepository()

	// Crear autor
	author := &domain.Author{Name: "John Doe", Email: "john@example.com"}
	authorRepo.Create(author)

	service := application.NewArticleService(articleRepo, authorRepo)

	// Crear artículo publicado reciente
	now := time.Now()
	article := &domain.Article{
		ID:              1,
		Title:           "Test Article",
		Content:         strings.Repeat("word ", 200), // 200 palabras
		Status:          domain.StatusPublished,
		AuthorID:        author.ID,
		PublicationDate: &now,
	}

	articleRepo.Create(article)

	// Calcular score
	// score = (palabras * 0.1) + (articulos_publicados * 5) + bonus
	// = (200 * 0.1) + (1 * 5) + 50
	// = 20 + 5 + 50 = 75
	score := service.CalculateScore(article)

	if score < 70 || score > 80 {
		t.Errorf("Expected score around 75, got %f", score)
	}
}

// TestValidateArticleMinimumWords - Prueba validación de palabras mínimas
func TestValidateArticleMinimumWords(t *testing.T) {
	articleRepo := NewMockArticleRepository()
	authorRepo := NewMockAuthorRepository()

	author := &domain.Author{Name: "John Doe", Email: "john@example.com"}
	authorRepo.Create(author)

	service := application.NewArticleService(articleRepo, authorRepo)

	// Artículo con muy pocas palabras
	article := &domain.Article{
		Title:    "Short Article",
		Content:  "Too short",
		Status:   domain.StatusDraft,
		AuthorID: author.ID,
	}

	articleRepo.Create(article)

	// Intentar publicar
	_, err := service.PublishArticle(article.ID)
	if err == nil {
		t.Error("Expected error for article with less than 120 words")
	}
}

// TestValidateArticleRepeatedWords - Prueba validación de palabras repetidas
func TestValidateArticleRepeatedWords(t *testing.T) {
	articleRepo := NewMockArticleRepository()
	authorRepo := NewMockAuthorRepository()

	author := &domain.Author{Name: "John Doe", Email: "john@example.com"}
	authorRepo.Create(author)

	service := application.NewArticleService(articleRepo, authorRepo)

	// Artículo con muchas palabras repetidas (> 35%)
	content := strings.Repeat("word word word word word ", 20) // Muchísimas palabras repetidas
	article := &domain.Article{
		Title:    "Repetitive Article",
		Content:  content,
		Status:   domain.StatusDraft,
		AuthorID: author.ID,
	}

	articleRepo.Create(article)

	// Intentar publicar
	_, err := service.PublishArticle(article.ID)
	if err == nil {
		t.Error("Expected error for article with too many repeated words")
	}
}

// TestGetTopAuthors - Prueba el endpoint de top autores
func TestGetTopAuthors(t *testing.T) {
	articleRepo := NewMockArticleRepository()
	authorRepo := NewMockAuthorRepository()

	// Crear 3 autores
	author1 := &domain.Author{Name: "Author 1", Email: "a1@example.com"}
	author2 := &domain.Author{Name: "Author 2", Email: "a2@example.com"}
	author3 := &domain.Author{Name: "Author 3", Email: "a3@example.com"}

	authorRepo.Create(author1)
	authorRepo.Create(author2)
	authorRepo.Create(author3)

	// Crear artículos publicados
	now := time.Now()
	for i := 0; i < 3; i++ {
		article := &domain.Article{
			Title:           "Article " + string(rune(i)),
			Content:         strings.Repeat("word ", 200),
			Status:          domain.StatusPublished,
			AuthorID:        author1.ID,
			PublicationDate: &now,
		}
		articleRepo.Create(article)
	}

	for i := 0; i < 2; i++ {
		article := &domain.Article{
			Title:           "Article " + string(rune(i)),
			Content:         strings.Repeat("word ", 150),
			Status:          domain.StatusPublished,
			AuthorID:        author2.ID,
			PublicationDate: &now,
		}
		articleRepo.Create(article)
	}

	service := application.NewArticleService(articleRepo, authorRepo)

	topAuthors, err := service.GetTopAuthors(3)
	if err != nil {
		t.Fatalf("Error getting top authors: %v", err)
	}

	if len(topAuthors) == 0 {
		t.Error("Expected at least one top author")
	}

	// Author 1 debe tener más score que Author 2 (más artículos)
	if len(topAuthors) > 1 {
		if topAuthors[0].ScoreAccumulated < topAuthors[1].ScoreAccumulated {
			t.Error("Top authors not sorted correctly by score")
		}
	}
}
