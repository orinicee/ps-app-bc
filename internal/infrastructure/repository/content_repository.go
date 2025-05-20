package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/orinicee/ps-app-bc/internal/domain"
)

// ContentRepository maneja las operaciones de base de datos para contenidos
type ContentRepository struct {
	db *sql.DB
}

// NewContentRepository crea una nueva instancia del repositorio
func NewContentRepository(db *sql.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

// Create inserta un nuevo contenido en la base de datos
func (r *ContentRepository) Create(content *domain.Content) error {
	query := `
		INSERT INTO contents (id, title, description, url, type, is_free, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	return r.db.QueryRow(
		query,
		content.ID,
		content.Title,
		content.Description,
		content.URL,
		content.Type.TypeName,
		content.IsFree,
		content.CreatedBy,
	).Scan(&content.ID)
}

// GetByID obtiene un contenido por su ID
func (r *ContentRepository) GetByID(id uuid.UUID) (*domain.Content, error) {
	query := `
		SELECT c.id, c.title, c.description, c.url, 
			   c.type, ct.description, c.is_free, 
			   c.created_at, c.created_by
		FROM contents c
		JOIN content_types ct ON c.type = ct.type_name
		WHERE c.id = $1`

	content := &domain.Content{}
	var contentType domain.ContentType

	err := r.db.QueryRow(query, id).Scan(
		&content.ID,
		&content.Title,
		&content.Description,
		&content.URL,
		&contentType.TypeName,
		&contentType.Description,
		&content.IsFree,
		&content.CreatedAt,
		&content.CreatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("content not found with id: %v", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error querying content: %w", err)
	}

	content.Type = contentType
	return content, nil
}

// List obtiene todos los contenidos con paginación
func (r *ContentRepository) List(limit, offset int) ([]*domain.Content, error) {
	query := `
		SELECT c.id, c.title, c.description, c.url, 
			   c.type, ct.description, c.is_free, 
			   c.created_at, c.created_by
		FROM contents c
		JOIN content_types ct ON c.type = ct.type_name
		ORDER BY c.created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying contents: %w", err)
	}
	defer rows.Close()

	var contents []*domain.Content

	for rows.Next() {
		content := &domain.Content{}
		var contentType domain.ContentType

		err := rows.Scan(
			&content.ID,
			&content.Title,
			&content.Description,
			&content.URL,
			&contentType.TypeName,
			&contentType.Description,
			&content.IsFree,
			&content.CreatedAt,
			&content.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning content: %w", err)
		}

		content.Type = contentType
		contents = append(contents, content)
	}

	return contents, nil
}

// Update actualiza un contenido existente
func (r *ContentRepository) Update(content *domain.Content) error {
	query := `
		UPDATE contents 
		SET title = $1, description = $2, url = $3, 
			type = $4, is_free = $5
		WHERE id = $6`

	result, err := r.db.Exec(
		query,
		content.Title,
		content.Description,
		content.URL,
		content.Type.TypeName,
		content.IsFree,
		content.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating content: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("content not found with id: %v", content.ID)
	}

	return nil
}

// Delete elimina un contenido
func (r *ContentRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM contents WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting content: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("content not found with id: %v", id)
	}

	return nil
}
