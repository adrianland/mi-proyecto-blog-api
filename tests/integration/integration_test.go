package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrianland/mi-proyecto-blog-api/interfaces/dto"
	"github.com/adrianland/mi-proyecto-blog-api/interfaces/http/handler"
	"github.com/adrianland/mi-proyecto-blog-api/internal/application"
	"github.com/adrianland/mi-proyecto-blog-api/internal/domain"
	"github.com/adrianland/mi-proyecto-blog-api/tests"
	"github.com/gin-gonic/gin"
)

// TestPublishArticleIntegration - Test de integración: crear y publicar artículo

func TestPublishArticleIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	articleRepo := tests.NewMockArticleRepository()
	authorRepo := tests.NewMockAuthorRepository()

	author := &domain.Author{
		Name:  "Test Author",
		Email: "test@example.com",
	}
	authorRepo.Create(author)

	articleService := application.NewArticleService(articleRepo, authorRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	router := gin.New()
	router.POST("/articulos", articleHandler.CreateArticle)
	router.PUT("/articulos/:id/publicar", articleHandler.PublishArticle)

	testContent := "La comunicación digital ha redefinido profundamente la manera en que las organizaciones producen y difunden información. Las plataformas en línea permiten compartir noticias con inmediatez y alcance global, conectando audiencias diversas en cuestión de segundos. Los equipos editoriales emplean sistemas especializados para estructurar contenidos claros, coherentes y atractivos. Las herramientas tecnológicas actuales facilitan la incorporación de recursos visuales, material audiovisual y gráficos interactivos que enriquecen la experiencia del lector. La producción informativa exige análisis riguroso, verificación constante y criterios éticos sólidos para garantizar credibilidad. Además, las métricas digitales ofrecen datos precisos sobre comportamiento, interacción y preferencias del público. Este entorno dinámico impulsa innovación permanente, adaptación estratégica y desarrollo de narrativas más efectivas. El ecosistema mediático evoluciona continuamente, influenciado por avances tecnológicos, cambios culturales y nuevas expectativas sociales."

	// Test 1: Crear artículo
	createReq := dto.CreateArticleRequest{
		Title:    "Test Article",
		Content:  testContent,
		AuthorID: author.ID,
	}

	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/articulos", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var createResp dto.SuccessResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)

	// DEBUG: Imprimir la respuesta
	t.Logf("Create Response Data: %+v", createResp.Data)

	articleData := createResp.Data.(map[string]interface{})
	articleID := int(articleData["id"].(float64))

	t.Logf("Article ID: %d", articleID)

	// Test 2: Publicar artículo
	publishReq := httptest.NewRequest("PUT", fmt.Sprintf("/articulos/%d/publicar", articleID), nil)
	publishW := httptest.NewRecorder()

	router.ServeHTTP(publishW, publishReq)

	// DEBUG: Imprimir respuesta de publicar
	t.Logf("Publish Status: %d", publishW.Code)
	t.Logf("Publish Body: %s", publishW.Body.String())

	if publishW.Code != http.StatusOK {
		t.Fatalf("Expected status 200 on publish, got %d. Body: %s", publishW.Code, publishW.Body.String())
	}
}

// TestArticleValidationOnPublish - Test de validaciones al publicar
func TestArticleValidationOnPublish(t *testing.T) {
	gin.SetMode(gin.TestMode)
	articleRepo := tests.NewMockArticleRepository()
	authorRepo := tests.NewMockAuthorRepository()

	author := &domain.Author{
		Name:  "Test Author",
		Email: "test@example.com",
	}
	authorRepo.Create(author)

	articleService := application.NewArticleService(articleRepo, authorRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	router := gin.New()
	router.POST("/articulos", articleHandler.CreateArticle)
	router.PUT("/articulos/:id/publicar", articleHandler.PublishArticle)

	// Test: Intentar publicar con menos de 120 palabras
	createReq := dto.CreateArticleRequest{
		Title:    "Short Article",
		Content:  "This is too short",
		AuthorID: author.ID,
	}

	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/articulos", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var createResp dto.SuccessResponse
	json.Unmarshal(w.Body.Bytes(), &createResp)
	articleData := createResp.Data.(map[string]interface{})
	articleID := int(articleData["id"].(float64))

	// Intentar publicar
	publishReq := httptest.NewRequest("PUT", fmt.Sprintf("/articulos/%d/publicar", articleID), nil)
	publishW := httptest.NewRecorder()

	router.ServeHTTP(publishW, publishReq)

	// Debe fallar
	if publishW.Code == http.StatusOK {
		t.Error("Expected error when publishing article with less than 120 words")
	}
}
