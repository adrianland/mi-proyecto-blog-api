package routes

import (
	"database/sql"

	"github.com/adrianland/mi-proyecto-blog-api/interfaces/http/handler"
	"github.com/adrianland/mi-proyecto-blog-api/interfaces/http/middleware"
	"github.com/adrianland/mi-proyecto-blog-api/internal/application"
	"github.com/adrianland/mi-proyecto-blog-api/internal/infrastructure/config"
	"github.com/adrianland/mi-proyecto-blog-api/internal/infrastructure/repository"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB, cfg *config.Config) *gin.Engine {
	// Crear repositorios
	articleRepo := repository.NewArticleRepository(db)
	authorRepo := repository.NewAuthorRepository(db)

	// Crear servicios
	articleService := application.NewArticleService(articleRepo, authorRepo)
	authorService := application.NewAuthorService(authorRepo)

	// Crear handlers
	articleHandler := handler.NewArticleHandler(articleService)
	authorHandler := handler.NewAuthorHandler(authorService, articleService)

	// Configurar router
	router := gin.Default()

	// Middlewares globales
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRequests, cfg.RateLimitWindow)
	router.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Rutas de Autores
	authors := router.Group("/autores")
	{
		authors.POST("", authorHandler.CreateAuthor)
		authors.GET("/top", authorHandler.GetTopAuthors)
		authors.GET("/:id/resumen", authorHandler.GetAuthorSummary)
		authors.GET("/:id/articulos", authorHandler.ListArticlesByAuthor)
	}

	// Rutas de Artículos
	articles := router.Group("/articulos")
	{
		articles.POST("", articleHandler.CreateArticle)
		articles.GET("", articleHandler.ListPublishedArticles)
		articles.GET("/:id", articleHandler.GetArticle)
		articles.PUT("/:id", articleHandler.UpdateArticle)
		articles.PUT("/:id/publicar", articleHandler.PublishArticle)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	return router
}
