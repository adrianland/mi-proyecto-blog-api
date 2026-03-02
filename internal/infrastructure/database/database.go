package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adrianland/mi-proyecto-blog-api/internal/infrastructure/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewConnection(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(cfg.DBMaxConnections)
	db.SetMaxIdleConns(cfg.DBMaxIdleConnections)
	db.SetConnMaxLifetime(cfg.DBConnectionMaxLifetime)

	// Intentar conectar con reintentos
	maxRetries := 5
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			return db, nil
		}
		lastErr = err
		time.Sleep(time.Second * 2)
	}

	return nil, fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, lastErr)
}

// RunMigrations ejecuta las migraciones de la base de datos
func RunMigrations(db *sql.DB) error {
	queries := []string{
		// Tabla de autores
		`CREATE TABLE IF NOT EXISTS authors (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_email (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// Tabla de artículos
		`CREATE TABLE IF NOT EXISTS articles (
			id INT PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR(255) NOT NULL,
			content LONGTEXT NOT NULL,
			status ENUM('BORRADOR', 'PUBLICADO') NOT NULL DEFAULT 'BORRADOR',
			author_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			publication_date TIMESTAMP NULL,
			FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE,
			INDEX idx_status (status),
			INDEX idx_author_id (author_id),
			INDEX idx_publication_date (publication_date),
			INDEX idx_author_status (author_id, status),
			FULLTEXT INDEX idx_content (content)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("migration error: %w", err)
		}
	}

	return nil
}
