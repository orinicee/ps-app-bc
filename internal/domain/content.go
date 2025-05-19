package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ContentType representa los tipos válidos de contenido
type ContentType struct {
	TypeName    string `json:"type_name"`
	Description string `json:"description"`
}

// Content representa la entidad principal de contenido
type Content struct {
	ID          uuid.UUID   `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	URL         string      `json:"url"`
	Type        ContentType `json:"type"`
	IsFree      bool        `json:"is_free"`
	CreatedAt   time.Time   `json:"created_at"`
	CreatedBy   uuid.UUID   `json:"created_by"`
}

func (c *Content) Validate() error {
	if c.Title == "" {
		return errors.New("title is required")
	}
	if c.URL == "" {
		return errors.New("URL is required")
	}
	if c.Type.TypeName == "" {
		return errors.New("content type is required")
	}
	if c.CreatedBy == uuid.Nil {
		return errors.New("creator ID is required")
	}
	return nil
}
