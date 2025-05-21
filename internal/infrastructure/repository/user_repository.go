package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/orinicee/ps-app-bc/internal/domain"
)

// UserRepository maneja las operaciones de base de datos para usuarios
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia del repositorio
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserta un nuevo usuario en la base de datos
func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password, first_name, last_name, phone_number, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	return r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Role,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone_number, role, active, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found with id: %v", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone_number, role, active, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found with email: %s", email)
	}
	if err != nil {
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return user, nil
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET email = $1, password = $2, first_name = $3, last_name = $4, 
			phone_number = $5, role = $6, active = $7, updated_at = $8
		WHERE id = $9`

	result, err := r.db.Exec(
		query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Role,
		user.Active,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found with id: %v", user.ID)
	}

	return nil
}

// Delete elimina un usuario
func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found with id: %v", id)
	}

	return nil
}

// List obtiene todos los usuarios con paginación
func (r *UserRepository) List(limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, phone_number, role, active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Role,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}
