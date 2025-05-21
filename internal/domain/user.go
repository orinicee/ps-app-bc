package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleClient Role = "client"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"-"` // No se serializa en JSON
	Role        Role      `json:"role"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewUser(email, password, firstName, lastName, phoneNumber string, role Role) (*User, error) {

	now := time.Now()
	user := &User{
		ID:          uuid.New(),
		Email:       email,
		Password:    password,
		Role:        role,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Validar el usuario
	if err := user.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	return user, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if u.FirstName == "" {
		return errors.New("first name is required")
	}
	if u.LastName == "" {
		return errors.New("last name is required")
	}
	if u.PhoneNumber == "" {
		return errors.New("phone number is required")
	}
	if u.Role == "" {
		return errors.New("role is required")
	}
	return nil
}
