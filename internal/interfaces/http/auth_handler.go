package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/orinicee/ps-app-bc/internal/domain"
	"github.com/orinicee/ps-app-bc/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

type registerRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Role        string `json:"role" binding:"required,oneof=admin client"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	Token     string       `json:"token"`
	User      userResponse `json:"user"`
	ExpiresAt string       `json:"expires_at"`
}

type userResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.RegisterInput{
		Email:       req.Email,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Role:        domain.Role(req.Role),
	}

	result, err := h.authUseCase.Register(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := authResponse{
		Token: result.Token,
		User: userResponse{
			ID:          result.User.ID.String(),
			Email:       result.User.Email,
			FirstName:   result.User.FirstName,
			LastName:    result.User.LastName,
			PhoneNumber: result.User.PhoneNumber,
			Role:        string(result.User.Role),
		},
		ExpiresAt: result.ExpiresAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.authUseCase.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := authResponse{
		Token: result.Token,
		User: userResponse{
			ID:          result.User.ID.String(),
			Email:       result.User.Email,
			FirstName:   result.User.FirstName,
			LastName:    result.User.LastName,
			PhoneNumber: result.User.PhoneNumber,
			Role:        string(result.User.Role),
		},
		ExpiresAt: result.ExpiresAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
