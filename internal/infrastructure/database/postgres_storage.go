package database

import (
	"database/sql"

	"github.com/orinicee/ps-app-bc/internal/domain"
)

// PostgresStorage implementa la interfaz domain.Storage para PostgreSQL
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage crea una nueva instancia de PostgresStorage
func NewPostgresStorage(db *sql.DB) domain.Storage {
	return &PostgresStorage{
		db: db,
	}
}

// HealthCheck implementa el método de la interfaz domain.Storage
func (s *PostgresStorage) HealthCheck() error {
	return s.db.Ping()
}

// Close implementa el método de la interfaz domain.Storage
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
