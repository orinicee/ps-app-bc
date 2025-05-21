package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/orinicee/ps-app-bc/internal/domain"
)

type AuthUseCase struct {
	userRepo domain.UserRepository
	jwtKey   []byte
}

func NewAuthUseCase(userRepo domain.UserRepository, jwtKey []byte) *AuthUseCase {
	if userRepo == nil {
		panic("user repository is required")
	}
	if len(jwtKey) == 0 {
		panic("jwt key is required")
	}
	return &AuthUseCase{
		userRepo: userRepo,
		jwtKey:   jwtKey,
	}
}

type RegisterInput struct {
	Email       string
	Password    string
	FirstName   string
	LastName    string
	PhoneNumber string
	Role        domain.Role
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	Token     string
	User      *domain.User
	ExpiresAt time.Time
}

func (uc *AuthUseCase) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	// Verificar si el usuario ya existe
	existingUser, err := uc.userRepo.GetByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Crear nuevo usuario
	user, err := domain.NewUser(input.Email, input.Password, input.FirstName, input.LastName, input.PhoneNumber, input.Role)
	if err != nil {
		return nil, err
	}

	// Guardar usuario
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generar token
	token, expiresAt, err := uc.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	// Buscar usuario por email
	user, err := uc.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verificar contraseña
	if !user.ValidatePassword(input.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Verificar si el usuario está activo
	if !user.Active {
		return nil, errors.New("user is not active")
	}

	// Generar token
	token, expiresAt, err := uc.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *AuthUseCase) generateToken(user *domain.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(uc.jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return uc.jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return nil, err
		}

		user, err := uc.userRepo.GetByID(userID)
		if err != nil {
			return nil, err
		}

		if !user.Active {
			return nil, errors.New("user is not active")
		}

		return user, nil
	}

	return nil, errors.New("invalid token")
}
