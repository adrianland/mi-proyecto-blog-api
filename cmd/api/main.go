package main

import (
	"fmt"
	"log"

	"github.com/adrianland/mi-proyecto-blog-api/interfaces/http/routes"
	"github.com/adrianland/mi-proyecto-blog-api/internal/infrastructure/config"
	"github.com/adrianland/mi-proyecto-blog-api/internal/infrastructure/database"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()
	log.Printf("Starting API on port %s with environment: %s", cfg.ServerPort, cfg.ServerEnv)

	// Conectar a la base de datos
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verificar conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection successful")

	// Ejecutar migraciones
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Crear router y configurar rutas
	router := routes.SetupRoutes(db, cfg)

	// Iniciar servidor
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
