package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/orinicee/ps-app-bc/internal/domain"
)

// ContentUseCases agrupa todos los casos de uso relacionados con contenido
type ContentUseCases struct {
	repo domain.ContentRepository
}

// NewContentUseCases crea una nueva instancia de los casos de uso
func NewContentUseCases(repo domain.ContentRepository) *ContentUseCases {
	if repo == nil {
		panic("repository is required")
	}
	return &ContentUseCases{repo: repo}
}

// UploadContent maneja la subida de nuevo contenido
func (uc *ContentUseCases) UploadContent(ctx context.Context, content *domain.Content) error {
	if content == nil {
		return errors.New("content is required")
	}

	// Validaciones de negocio
	if err := content.Validate(); err != nil {
		return err
	}

	// Reglas de negocio específicas
	if err := uc.applyUploadRules(content); err != nil {
		return err
	}

	// Asignar ID y timestamp si no existen
	if content.ID == uuid.Nil {
		content.ID = uuid.New()
	}
	if content.CreatedAt.IsZero() {
		content.CreatedAt = time.Now()
	}

	// Persistencia
	return uc.repo.Create(content)
}

// GetContent obtiene un contenido por su ID
func (uc *ContentUseCases) GetContent(ctx context.Context, id uuid.UUID) (*domain.Content, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid content ID")
	}

	content, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Aquí podrías agregar reglas de negocio adicionales
	// Por ejemplo, verificar permisos de acceso
	return content, nil
}

// ListContents obtiene una lista paginada de contenidos
func (uc *ContentUseCases) ListContents(ctx context.Context, page, pageSize int) ([]*domain.Content, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	contents, err := uc.repo.List(pageSize, offset)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// UpdateContent actualiza un contenido existente
func (uc *ContentUseCases) UpdateContent(ctx context.Context, content *domain.Content) error {
	if content == nil {
		return errors.New("content is required")
	}

	if content.ID == uuid.Nil {
		return errors.New("content ID is required")
	}

	// Validar que el contenido existe
	existing, err := uc.repo.GetByID(content.ID)
	if err != nil {
		return err
	}

	// Aplicar reglas de negocio para actualización
	if err := uc.applyUpdateRules(existing, content); err != nil {
		return err
	}

	// Persistir cambios
	return uc.repo.Update(content)
}

// DeleteContent elimina un contenido
func (uc *ContentUseCases) DeleteContent(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid content ID")
	}

	// Verificar que el contenido existe
	if _, err := uc.repo.GetByID(id); err != nil {
		return err
	}

	// Aplicar reglas de negocio para eliminación
	if err := uc.applyDeleteRules(id); err != nil {
		return err
	}

	return uc.repo.Delete(id)
}

// Reglas de negocio específicas
func (uc *ContentUseCases) applyUploadRules(content *domain.Content) error {
	// Implementar reglas específicas para subida
	// Por ejemplo:
	// - Verificar límites de tamaño
	// - Validar formato de URL
	// - Verificar permisos del usuario
	return nil
}

func (uc *ContentUseCases) applyUpdateRules(existing, new *domain.Content) error {
	// Implementar reglas específicas para actualización
	// Por ejemplo:
	// - Verificar que el usuario tiene permisos
	// - Validar cambios permitidos
	// - Mantener campos inmutables
	return nil
}

func (uc *ContentUseCases) applyDeleteRules(id uuid.UUID) error {
	// Implementar reglas específicas para eliminación
	// Por ejemplo:
	// - Verificar que el usuario tiene permisos
	// - Validar si el contenido puede ser eliminado
	return nil
}
