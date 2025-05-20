package domain

import (
	"github.com/google/uuid"
)

// ContentRepository define la interfaz para las operaciones de persistencia de contenido
type ContentRepository interface {
	Create(content *Content) error
	GetByID(id uuid.UUID) (*Content, error)
	List(limit, offset int) ([]*Content, error)
	Update(content *Content) error
	Delete(id uuid.UUID) error
}
