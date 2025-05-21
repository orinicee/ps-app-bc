package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/orinicee/ps-app-bc/internal/domain"
	"github.com/orinicee/ps-app-bc/internal/infrastructure/database"
)

func TestUserRepository_Create(t *testing.T) {
	config := setupTestDB(t)
	db, err := database.NewConnection(*config)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	defer cleanupTestData(db)

	repo := NewUserRepository(db)

	// Crear un usuario de prueba
	user := &domain.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		Password:    "hashed_password",
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "1234567890",
		Role:        domain.RoleClient,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Intentar crear el usuario
	err = repo.Create(user)
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	// Verificar que se puede recuperar el usuario
	retrieved, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Error retrieving user: %v", err)
	}

	// Verificar que los datos coinciden
	if retrieved.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved.Email)
	}
	if retrieved.FirstName != user.FirstName {
		t.Errorf("Expected first name %s, got %s", user.FirstName, retrieved.FirstName)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	config := setupTestDB(t)
	db, err := database.NewConnection(*config)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	defer cleanupTestData(db)

	repo := NewUserRepository(db)

	// Crear un usuario de prueba
	user := &domain.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		Password:    "hashed_password",
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "1234567890",
		Role:        domain.RoleClient,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Crear el usuario
	err = repo.Create(user)
	if err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	// Buscar por email
	retrieved, err := repo.GetByEmail(user.Email)
	if err != nil {
		t.Fatalf("Error retrieving user by email: %v", err)
	}

	// Verificar que los datos coinciden
	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %v, got %v", user.ID, retrieved.ID)
	}
	if retrieved.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	config := setupTestDB(t)
	db, err := database.NewConnection(*config)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	defer cleanupTestData(db)

	repo := NewUserRepository(db)

	// Crear un usuario de prueba
	user := &domain.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		Password:    "hashed_password",
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "1234567890",
		Role:        domain.RoleClient,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Crear el usuario
	err = repo.Create(user)
	if err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	// Modificar el usuario
	user.FirstName = "Updated"
	user.LastName = "Name"
	user.UpdatedAt = time.Now()

	// Actualizar el usuario
	err = repo.Update(user)
	if err != nil {
		t.Fatalf("Error updating user: %v", err)
	}

	// Verificar los cambios
	updated, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Error retrieving updated user: %v", err)
	}

	if updated.FirstName != "Updated" {
		t.Errorf("Expected first name 'Updated', got %s", updated.FirstName)
	}
	if updated.LastName != "Name" {
		t.Errorf("Expected last name 'Name', got %s", updated.LastName)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	config := setupTestDB(t)
	db, err := database.NewConnection(*config)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	defer cleanupTestData(db)

	repo := NewUserRepository(db)

	// Crear un usuario de prueba
	user := &domain.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		Password:    "hashed_password",
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "1234567890",
		Role:        domain.RoleClient,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Crear el usuario
	err = repo.Create(user)
	if err != nil {
		t.Fatalf("Error creating test user: %v", err)
	}

	// Eliminar el usuario
	err = repo.Delete(user.ID)
	if err != nil {
		t.Fatalf("Error deleting user: %v", err)
	}

	// Verificar que el usuario ya no existe
	_, err = repo.GetByID(user.ID)
	if err == nil {
		t.Error("Expected error when retrieving deleted user, got nil")
	}
}

func TestUserRepository_List(t *testing.T) {
	config := setupTestDB(t)
	db, err := database.NewConnection(*config)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	defer cleanupTestData(db)

	repo := NewUserRepository(db)

	// Crear varios usuarios de prueba
	users := []*domain.User{
		{
			ID:          uuid.New(),
			Email:       "test1@example.com",
			Password:    "hashed_password",
			FirstName:   "Test",
			LastName:    "User 1",
			PhoneNumber: "1234567890",
			Role:        domain.RoleClient,
			Active:      true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Email:       "test2@example.com",
			Password:    "hashed_password",
			FirstName:   "Test",
			LastName:    "User 2",
			PhoneNumber: "1234567891",
			Role:        domain.RoleAdmin,
			Active:      true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Crear los usuarios
	for _, user := range users {
		err = repo.Create(user)
		if err != nil {
			t.Fatalf("Error creating test user: %v", err)
		}
	}

	// Probar paginación
	tests := []struct {
		name     string
		page     int
		pageSize int
		want     int
	}{
		{"first page", 1, 1, 1},
		{"second page", 2, 1, 1},
		{"full page", 1, 2, 2},
		{"empty page", 3, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrieved, err := repo.List(tt.pageSize, (tt.page-1)*tt.pageSize)
			if err != nil {
				t.Fatalf("Error listing users: %v", err)
			}

			if len(retrieved) != tt.want {
				t.Errorf("Expected %d users, got %d", tt.want, len(retrieved))
			}
		})
	}
}
