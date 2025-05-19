package repository

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/orinicee/ps-app-bc/internal/domain"
	"github.com/orinicee/ps-app-bc/internal/infrastructure/database"
)

func setupTestDB(_ *testing.T) *database.Config {
	return &database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres123",
		DBName:   "ps_app",
		SSLMode:  "disable",
	}
}

func createTestUser(db *sql.DB) (uuid.UUID, error) {
	userID := uuid.New()
	query := `
		INSERT INTO users (id, email, name)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := db.QueryRow(query, userID, "test@example.com", "Test User").Scan(&userID)
	return userID, err
}

func cleanupTestData(db *sql.DB) error {
	// Primero eliminar contenidos (por la llave foránea)
	_, err := db.Exec("DELETE FROM contents")
	if err != nil {
		return err
	}
	// Luego eliminar usuarios
	_, err = db.Exec("DELETE FROM users")
	return err
}

func TestContentRepository_Create(t *testing.T) {
	config := setupTestDB(t)
	db, err := database.NewConnection(*config)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	defer cleanupTestData(db) // Limpiar datos después del test

	// Crear un usuario de prueba
	userID, err := createTestUser(db)
	if err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	repo := NewContentRepository(db)

	// Crear un contenido de prueba
	content := &domain.Content{
		ID:          uuid.New(),
		Title:       "Test Content",
		Description: "Test Description",
		URL:         "https://test.com/video",
		Type: domain.ContentType{
			TypeName:    "video",
			Description: "Contenido tipo video",
		},
		IsFree:    true,
		CreatedBy: userID,
	}

	// Intentar crear el contenido
	err = repo.Create(content)
	if err != nil {
		t.Fatalf("Error creating content: %v", err)
	}

	// Verificar que se puede recuperar el contenido
	retrieved, err := repo.GetByID(content.ID)
	if err != nil {
		t.Fatalf("Error retrieving content: %v", err)
	}

	// Verificar que los datos coinciden
	if retrieved.Title != content.Title {
		t.Errorf("Expected title %s, got %s", content.Title, retrieved.Title)
	}
	if retrieved.Type.TypeName != content.Type.TypeName {
		t.Errorf("Expected type %s, got %s", content.Type.TypeName, retrieved.Type.TypeName)
	}
}

func TestContentRepository_List(t *testing.T) {
	config := setupTestDB(t)
	db, err := database.NewConnection(*config)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	defer cleanupTestData(db) // Limpiar datos después del test

	// Crear un usuario de prueba
	userID, err := createTestUser(db)
	if err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	repo := NewContentRepository(db)

	// contenidos de prueba
	contents := []*domain.Content{
		{
			ID:          uuid.New(),
			Title:       "Test Content 1",
			Description: "Test Description 1",
			URL:         "https://test.com/video1",
			Type: domain.ContentType{
				TypeName:    "video",
				Description: "Contenido tipo video",
			},
			IsFree:    true,
			CreatedBy: userID,
		},
		{
			ID:          uuid.New(),
			Title:       "Test Content 2",
			Description: "Test Description 2",
			URL:         "https://test.com/video2",
			Type: domain.ContentType{
				TypeName:    "audio",
				Description: "Contenido tipo audio",
			},
			IsFree:    false,
			CreatedBy: userID,
		},
	}

	// Crear los contenidos
	for _, content := range contents {
		err = repo.Create(content)
		if err != nil {
			t.Fatalf("Error creating test content: %v", err)
		}
	}

	// Obtener lista de contenidos
	retrieved, err := repo.List(10, 0)
	if err != nil {
		t.Fatalf("Error listing contents: %v", err)
	}

	// Verificar que obtenemos los contenidos creados
	if len(retrieved) != len(contents) {
		t.Errorf("Expected %d contents, got %d", len(contents), len(retrieved))
	}

	// Verificar que los contenidos tienen los campos requeridos
	for _, content := range retrieved {
		if content.ID == uuid.Nil {
			t.Error("Content has nil ID")
		}
		if content.Title == "" {
			t.Error("Content has empty title")
		}
		if content.Type.TypeName == "" {
			t.Error("Content has empty type")
		}
	}
}
