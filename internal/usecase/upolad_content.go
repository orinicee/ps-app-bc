package usecase

import (
	"context"
	"errors"

	"github.com/orinicee/ps-app-bc/internal/domain"
)

type UploadContentUseCase struct {
	repo domain.ContentRepository
}

func NewUploadContentUseCase(repo domain.ContentRepository) *UploadContentUseCase {
	if repo == nil {
		panic("repository is required")
	}
	return &UploadContentUseCase{repo: repo}
}

func (uc *UploadContentUseCase) Execute(ctx context.Context, content *domain.Content) error {
	if content == nil {
		return errors.New("content is required")
	}

	// Validaciones de negocio
	if err := content.Validate(); err != nil {
		return err
	}

	// Reglas de negocio específicas
	if err := uc.applyBusinessRules(content); err != nil {
		return err
	}

	// Persistencia
	return uc.repo.Save(content)
}

func (uc *UploadContentUseCase) applyBusinessRules(content *domain.Content) error {
	// Implementa reglas de negocio específicas aquí
	return nil
}
