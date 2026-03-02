package application

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/adrianland/mi-proyecto-blog-api/internal/domain"
)

type ArticleService struct {
	articleRepo domain.ArticleRepository
	authorRepo  domain.AuthorRepository
}

func NewArticleService(articleRepo domain.ArticleRepository, authorRepo domain.AuthorRepository) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		authorRepo:  authorRepo,
	}
}

// CreateArticle crea un nuevo artículo en estado BORRADOR
func (s *ArticleService) CreateArticle(title, content string, authorID int) (*domain.Article, error) {
	if title == "" || content == "" {
		return nil, fmt.Errorf("title and content are required")
	}

	author, err := s.authorRepo.GetByID(authorID)
	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	article := &domain.Article{
		Title:     title,
		Content:   content,
		Status:    domain.StatusDraft,
		AuthorID:  author.ID,
		CreatedAt: time.Now(),
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return article, nil
}

// GetArticle obtiene un artículo por ID
func (s *ArticleService) GetArticle(id int) (*domain.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("article not found: %w", err)
	}

	// Calcular score si está publicado
	if article.Status == domain.StatusPublished {
		article.Score = s.CalculateScore(article)
	}

	return article, nil
}

// UpdateArticle actualiza un artículo (solo si está en BORRADOR)
func (s *ArticleService) UpdateArticle(id int, title, content string) (*domain.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("article not found: %w", err)
	}

	if article.Status != domain.StatusDraft {
		return nil, fmt.Errorf("cannot edit published articles")
	}

	if title != "" {
		article.Title = title
	}
	if content != "" {
		article.Content = content
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return article, nil
}

// PublishArticle publica un artículo con validaciones
func (s *ArticleService) PublishArticle(id int) (*domain.Article, error) {
	article, err := s.articleRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("article not found: %w", err)
	}

	if article.Status == domain.StatusPublished {
		return nil, fmt.Errorf("article is already published")
	}

	// Validaciones
	if err := s.validateArticleForPublishing(article); err != nil {
		return nil, err
	}

	// Cambiar estado y asignar fecha de publicación
	now := time.Now()
	article.Status = domain.StatusPublished
	article.PublicationDate = &now

	if err := s.articleRepo.Update(article); err != nil {
		return nil, fmt.Errorf("failed to publish article: %w", err)
	}

	// Calcular y asignar score
	article.Score = s.CalculateScore(article)

	return article, nil
}

// ListPublishedArticles lista artículos publicados con paginación
func (s *ArticleService) ListPublishedArticles(page, pageSize int) ([]domain.ArticleWithAuthor, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	articles, total, err := s.articleRepo.ListPublished(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list articles: %w", err)
	}

	// Calcular score para cada artículo
	for i := range articles {
		article := &domain.Article{
			ID:              articles[i].ID,
			Title:           articles[i].Title,
			Content:         articles[i].Content,
			Status:          articles[i].Status,
			AuthorID:        articles[i].AuthorID,
			CreatedAt:       articles[i].CreatedAt,
			PublicationDate: articles[i].PublicationDate,
			WordCount:       articles[i].WordCount,
		}
		articles[i].Score = s.CalculateScore(article)
	}

	return articles, total, nil
}

// ListArticlesByAuthor lista artículos por autor
func (s *ArticleService) ListArticlesByAuthor(authorID int, status string, page, pageSize int) ([]domain.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	articles, total, err := s.articleRepo.ListByAuthor(authorID, status, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list articles: %w", err)
	}

	// Calcular score para artículos publicados
	for i := range articles {
		if articles[i].Status == domain.StatusPublished {
			articles[i].Score = s.CalculateScore(&articles[i])
		}
	}

	return articles, total, nil
}

// GetTopAuthors retorna los N autores con mayor score acumulado
func (s *ArticleService) GetTopAuthors(n int) ([]domain.TopAuthor, error) {
	articles, err := s.articleRepo.GetAllPublished()
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}

	// Mapa para acumular scores por autor
	authorScores := make(map[int]float64)
	authorNames := make(map[int]string)

	for _, article := range articles {
		score := s.CalculateScore(&article)
		authorScores[article.AuthorID] += score

		// Guardar nombre del autor si no lo tenemos
		if _, exists := authorNames[article.AuthorID]; !exists {
			author, err := s.authorRepo.GetByID(article.AuthorID)
			if err == nil {
				authorNames[article.AuthorID] = author.Name
			}
		}
	}

	// Convertir a slice y ordenar
	topAuthors := make([]domain.TopAuthor, 0, len(authorScores))
	for authorID, score := range authorScores {
		topAuthors = append(topAuthors, domain.TopAuthor{
			AuthorID:         authorID,
			AuthorName:       authorNames[authorID],
			ScoreAccumulated: math.Round(score*100) / 100,
		})
	}

	// Ordenar por score descendente
	for i := 0; i < len(topAuthors)-1; i++ {
		for j := i + 1; j < len(topAuthors); j++ {
			if topAuthors[j].ScoreAccumulated > topAuthors[i].ScoreAccumulated {
				topAuthors[i], topAuthors[j] = topAuthors[j], topAuthors[i]
			}
		}
	}

	// Limitar a N
	if n > 0 && n < len(topAuthors) {
		topAuthors = topAuthors[:n]
	}

	return topAuthors, nil
}

// ========== VALIDACIONES Y CÁLCULOS INTERNOS ==========

// validateArticleForPublishing valida que el artículo cumpla los requisitos
func (s *ArticleService) validateArticleForPublishing(article *domain.Article) error {
	wordCount := countWords(article.Content)

	// Validar mínimo de palabras
	if wordCount < 120 {
		return fmt.Errorf("article must have at least 120 words, current: %d", wordCount)
	}

	// Validar palabras repetidas
	repeatedPercent := calculateRepeatedWordsPercentage(article.Content)
	if repeatedPercent > 35 {
		return fmt.Errorf("article has %.2f%% repeated words, maximum allowed is 35%%", repeatedPercent)
	}

	return nil
}

// CalculateScore calcula el score dinámico de un artículo publicado
func (s *ArticleService) CalculateScore(article *domain.Article) float64 {
	if article.Status != domain.StatusPublished || article.PublicationDate == nil {
		return 0
	}

	score := 0.0

	// Componente 1: Palabras del artículo
	wordCount := countWords(article.Content)
	score += float64(wordCount) * 0.1

	// Componente 2: Artículos publicados del autor
	publishedCount, err := s.articleRepo.CountPublishedByAuthor(article.AuthorID)
	if err == nil {
		score += float64(publishedCount) * 5
	}

	// Componente 3: Bonus por recencia
	score += s.calculateBonusReciente(*article.PublicationDate)

	return math.Round(score*100) / 100
}

// calculateBonusReciente calcula el bonus por publicación reciente
func (s *ArticleService) calculateBonusReciente(publicationDate time.Time) float64 {
	now := time.Now()
	duration := now.Sub(publicationDate)

	if duration < 24*time.Hour {
		return 50
	}
	if duration < 72*time.Hour {
		return 20
	}
	return 0
}

// countWords cuenta las palabras en el contenido
func countWords(content string) int {
	words := strings.Fields(content)
	return len(words)
}

// calculateRepeatedWordsPercentage calcula el porcentaje de palabras repetidas
func calculateRepeatedWordsPercentage(content string) float64 {
	words := strings.Fields(strings.ToLower(content))
	if len(words) == 0 {
		return 0
	}

	// Contar ocurrencias
	wordCount := make(map[string]int)
	for _, word := range words {
		// Limpiar puntuación
		word = strings.Trim(word, ".,!?;:\"'()-")
		if word != "" {
			wordCount[word]++
		}
	}

	// Calcular palabras repetidas
	totalRepeated := 0
	for _, count := range wordCount {
		if count > 1 {
			totalRepeated += count - 1
		}
	}

	return (float64(totalRepeated) / float64(len(words))) * 100
}
